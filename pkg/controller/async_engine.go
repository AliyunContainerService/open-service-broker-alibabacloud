package controller

import (
	"fmt"
	"sync"
	"time"

	"github.com/AliyunContainerService/open-service-broker-alibabacloud/pkg/brokerapi"
	"github.com/golang/glog"
)

const AsyncLoopTime = 10

type Async interface {
	Run()
	AddAsyncInstance(instanceInfo *InstanceRunInfo, persistent bool) error
	DeleteAsyncInstance(instanceID string) error
	UpdateAsyncInstance(instanceInfo *InstanceRunInfo) error
	GetAsyncInstance(instanceId string) *InstanceRunInfo
	GetAsyncBinding(instanceId, bindingID string) *BindRunInfo
	GetProvisionInstanceLastOperation(instanceInfo *InstanceRunInfo)
	GetDeprovisionInstanceLastOperation(instanceInfo *InstanceRunInfo)
	GetBindingLastOperation(instanceInfo *InstanceRunInfo, bindInfo *BindRunInfo)
	GetUnBindingLastOperation(instanceInfo *InstanceRunInfo, bindInfo *BindRunInfo)
}

type AsyncEngine struct {
	storageProvider      *StorageProvider
	instancesRunMap      map[string]*InstanceRunInfo
	instancesRunMapMutex *sync.Mutex
}

func NewAsyncEngine() *AsyncEngine {

	asyncEngine := &AsyncEngine{
		instancesRunMap:      make(map[string]*InstanceRunInfo),
		instancesRunMapMutex: new(sync.Mutex),
	}

	sp := NewStorageProvider()
	if sp != nil {
		asyncEngine.storageProvider = sp
		asyncEngine.initStorageProvider()
	} else {
		asyncEngine.storageProvider = nil
	}

	asyncEngine.Run()

	return asyncEngine
}

func (a *AsyncEngine) initStorageProvider() {

	instanceInfo, err := a.storageProvider.GetAllObject()
	if err == nil {
		glog.Infof("Found existed instance info in storage, loading to async engine.")
		for index := range instanceInfo {
			err = a.AddAsyncInstance(instanceInfo[index], false)
			if err != nil {
				glog.Warningln(err)
				return
			}
		}
		return
	} else {
		glog.Infof("Not found existed instance info in storage, start to initiate.")
		err := a.storageProvider.Init()
		if err != nil {
			glog.Infof("Failed to initiate service broker info in storage")
			return
		}
	}

}

func (a *AsyncEngine) Run() {
	go a.asyncRun()
}

func (a *AsyncEngine) AddAsyncInstance(instanceInfo *InstanceRunInfo, persistent bool) error {
	if _, ok := a.instancesRunMap[instanceInfo.InstanceId]; ok {
		glog.Infof("Can't add instance info in async engine, since there is one already exist!")
		return fmt.Errorf("Can't add instance info in async engine, since there is one already exist!")
	}
	a.instancesRunMapMutex.Lock()
	a.instancesRunMap[instanceInfo.InstanceId] = instanceInfo
	a.instancesRunMapMutex.Unlock()
	if persistent && a.storageProvider != nil {
		a.storageProvider.WriteObject(*instanceInfo)
	}
	return nil
}

func (a *AsyncEngine) DeleteAsyncInstance(instanceID string) error {
	if _, ok := a.instancesRunMap[instanceID]; ok {
		a.instancesRunMapMutex.Lock()
		delete(a.instancesRunMap, instanceID)
		a.instancesRunMapMutex.Unlock()
		if a.storageProvider != nil {
			a.storageProvider.DeleteObject(instanceID)
		}
		return nil
	}

	return fmt.Errorf("There is no such instance when delete instance.")
}

func (a *AsyncEngine) UpdateAsyncInstance(instanceInfo *InstanceRunInfo) error {
	if _, ok := a.instancesRunMap[instanceInfo.InstanceId]; ok {

		err := a.DeleteAsyncInstance(instanceInfo.InstanceId)
		if err != nil {
			return fmt.Errorf("UpdateAsyncInstance delete exist instance info failed. err:%v", err)
		}
		err = a.AddAsyncInstance(instanceInfo, true)
		if err != nil {
			return fmt.Errorf("UpdateAsyncInstance add new instance info failed.err:%v", err)
		}
		return nil
	}

	return fmt.Errorf("There is no such instance when update instance.")
}

func (a *AsyncEngine) GetAsyncInstance(instanceId string) *InstanceRunInfo {
	if value, ok := a.instancesRunMap[instanceId]; ok {
		return value
	}
	return nil
}

func (a *AsyncEngine) GetAsyncBinding(instanceId, bindingID string) *BindRunInfo {
	if value, ok := a.instancesRunMap[instanceId]; ok {
		if binding, ok := value.Bindings[bindingID]; ok {
			return binding
		}
	}
	return nil
}

func (a *AsyncEngine) asyncBindingRun(instanceInfo *InstanceRunInfo) {
	for _, bindingInfo := range instanceInfo.Bindings {
		switch bindingInfo.Status {
		case StateBindingInProgress:
			go a.GetBindingLastOperation(instanceInfo, bindingInfo)
		case StateUnBindingInProgress:
			go a.GetUnBindingLastOperation(instanceInfo, bindingInfo)
		default:
			continue
		}
	}
	return
}

func (a *AsyncEngine) asyncRun() {
	for {
		time.Sleep(AsyncLoopTime * time.Second)

		for _, instanceInfo := range a.instancesRunMap {
			//glog.Infof("AsyncRun in %+v\n", instanceInfo)
			a.instancesRunMapMutex.Lock()
			switch instanceInfo.Status {
			case StateProvisionInstanceInProgress:
				go a.GetProvisionInstanceLastOperation(instanceInfo)
			case StateDeprovisionInstanceInProgress:
				go a.GetDeprovisionInstanceLastOperation(instanceInfo)
			case StateProvisionInstanceSucceeded:
				a.asyncBindingRun(instanceInfo)
			default:
			}
			a.instancesRunMapMutex.Unlock()
		}
	}
}

func GetBrokerByName(brokerName string) brokerapi.ServiceBroker {
	baseController := GetBaseController()
	return baseController.getBrokerByName(brokerName)
}

func (a *AsyncEngine) GetProvisionInstanceLastOperation(instanceInfo *InstanceRunInfo) {
	a.instancesRunMapMutex.Lock()
	defer a.instancesRunMapMutex.Unlock()
	glog.Infof("GetProvisionInstanceILastOperation with instanceInfo: %v", instanceInfo)
	broker := GetBrokerByName(instanceInfo.BrokerName)
	if broker == nil {
		glog.Infof("getProvisionInstanceILastOperation failed.")
		return
	}

	status := ""
	ok, err := broker.GetInstanceStatus(instanceInfo.InstanceId, instanceInfo.ServiceID,
		instanceInfo.PlanID, instanceInfo.Parameter)
	if err != nil {
		status = StateProvisionInstanceFailed
	} else if ok == true {
		status = StateProvisionInstanceSucceeded
	} else {
		status = StateProvisionInstanceInProgress
	}
	instanceInfo.Status = status
	if a.storageProvider != nil {
		a.storageProvider.WriteObject(*instanceInfo)
	}
	glog.Infof("GetProvisionInstanceLastOperation:%v.", instanceInfo)
	return
}

func (a *AsyncEngine) GetDeprovisionInstanceLastOperation(instanceInfo *InstanceRunInfo) {
	a.instancesRunMapMutex.Lock()
	defer a.instancesRunMapMutex.Unlock()
	glog.Infof("getDeprovisionInstanceLastOperation in.")
	broker := GetBrokerByName(instanceInfo.BrokerName)
	if broker == nil {
		glog.Infof("getDeprovisionInstanceILastOperation failed.")
		return
	}

	status := ""
	ok, err := broker.GetInstanceStatus(instanceInfo.InstanceId, instanceInfo.ServiceID,
		instanceInfo.PlanID, instanceInfo.Parameter)
	if err != nil {
		status = StateDeprovisionInstanceFailed
	} else if ok == true {
		status = StateDeprovisionInstanceInProgress
	} else {
		status = StateDeprovisionInstanceSucceeded
	}
	instanceInfo.Status = status
	if a.storageProvider != nil {
		a.storageProvider.WriteObject(*instanceInfo)
	}
	glog.Infof("getDeprovisionInstanceLastOperation success:%v.", instanceInfo)
	return
}

func (a *AsyncEngine) GetBindingLastOperation(instanceInfo *InstanceRunInfo, bindInfo *BindRunInfo) {
	a.instancesRunMapMutex.Lock()
	defer a.instancesRunMapMutex.Unlock()
	glog.Infof("getBindingILastOperation in.")
	broker := GetBrokerByName(instanceInfo.BrokerName)
	if broker == nil {
		return
	}

	status := ""
	ok, err := broker.GetBindingStatus(instanceInfo.InstanceId, instanceInfo.ServiceID,
		instanceInfo.PlanID, bindInfo.BindingId, instanceInfo.Parameter, bindInfo.Parameter)
	if err != nil {
		status = StateBindingFailed
	} else if ok == true {
		status = StateBindingSucceeded
	} else {
		status = StateBindingInProgress
	}
	bindInfo.Status = status
	if a.storageProvider != nil {
		a.storageProvider.WriteObject(*instanceInfo)
	}
	return
}

func (a *AsyncEngine) GetUnBindingLastOperation(instanceInfo *InstanceRunInfo, bindInfo *BindRunInfo) {
	a.instancesRunMapMutex.Lock()
	defer a.instancesRunMapMutex.Unlock()
	broker := GetBrokerByName(instanceInfo.BrokerName)
	if broker == nil {
		return
	}
	status := ""
	ok, err := broker.GetBindingStatus(instanceInfo.InstanceId, instanceInfo.ServiceID,
		instanceInfo.PlanID, bindInfo.BindingId, instanceInfo.Parameter, bindInfo.Parameter)
	if err != nil {
		status = StateUnBindingFailed
	} else if ok == true {
		status = StateUnBindingInProgress
	} else {
		status = StateUnBindingSucceeded
	}
	bindInfo.Status = status
	if a.storageProvider != nil {
		a.storageProvider.WriteObject(*instanceInfo)
	}
	return
}
