package controller

import (
	"fmt"
	"sync"

	"github.com/AliyunContainerService/open-service-broker-alibabacloud/pkg/brokerapi"
	"github.com/golang/glog"
)

const (
	StateProvisionInstanceInProgress = "Provision instance in progress"
	StateProvisionInstanceSucceeded  = "Provision instance succeeded"
	StateProvisionInstanceFailed     = "Provision instance failed"

	StateDeprovisionInstanceInProgress = "Deprovision instance in progress"
	StateDeprovisionInstanceSucceeded  = "Deprovision instance succeeded"
	StateDeprovisionInstanceFailed     = "Deprovision instance failed"

	StateBindingInProgress = "Binding in progress"
	StateBindingSucceeded  = "Binding succeeded"
	StateBindingFailed     = "Binding failed"

	StateUnBindingInProgress = "Unbinding in progress"
	StateUnBindingSucceeded  = "Unbinding succeeded"
	StateUnBindingFailed     = "Unbinding failed"
)

// BaseController implements the generic logic shared by all concrete service brokers.
// Like async process of create/delete service instance, bind/unbind and get last operation
type BaseController struct {
	asyncEngine     *AsyncEngine
	serviceMap      map[string]*brokerapi.Catalog
	brokerMap       map[string]brokerapi.ServiceBroker
	controllerMutex *sync.Mutex
}

var baseController *BaseController

func NewBaseController() *BaseController {
	asyncEngine := NewAsyncEngine()
	if asyncEngine == nil {
		glog.Infof("BaseControllerInit create asyncengine failed.")
		return nil
	}
	baseController = new(BaseController)
	baseController.asyncEngine = asyncEngine
	baseController.brokerMap = make(map[string]brokerapi.ServiceBroker)
	baseController.serviceMap = make(map[string]*brokerapi.Catalog)
	baseController.controllerMutex = new(sync.Mutex)
	return baseController
}

func GetBaseController() *BaseController {
	return baseController
}

func RegisterBrokers(brokers map[string]brokerapi.ServiceBroker) bool {
	registered := false
	for brokerName, broker := range brokers {
		glog.Infof("Registering service broker of %s", brokerName)
		if broker == nil {
			glog.Errorf("Attempt to register a nil broker, which is invalid.")
			registered = registered || false
			continue
		}
		controller := GetBaseController()
		controller.controllerMutex.Lock()
		//defer controller.controllerMutex.Unlock()
		if _, existed := controller.brokerMap[brokerName]; !existed {
			controller.brokerMap[brokerName] = broker
		}
		if _, existed := controller.serviceMap[brokerName]; existed {
			glog.Infof("Service broker: %s has been registered already. \n", brokerName)
			registered = true
			continue
		}
		// store service broker's Catalog() info
		catalog, err := broker.Catalog()
		if err != nil {
			glog.Errorf("Failed to get Catalog information of service broker: %s", brokerName)
			return false
		}
		controller.serviceMap[brokerName] = catalog
		controller.controllerMutex.Unlock()
		glog.Infof("Success to register service broker %s - %v", brokerName, controller.serviceMap[brokerName])
		registered = true
	}
	return registered
}

func RegisterBroker(brokerName string, broker brokerapi.ServiceBroker) bool {
	glog.Infof("Registering service broker of %s", brokerName)
	if broker == nil {
		glog.Errorf("Attempt to register a nil broker, which is invalid.")
		return false
	}
	controller := GetBaseController()
	controller.controllerMutex.Lock()
	defer controller.controllerMutex.Unlock()
	if _, existed := controller.brokerMap[brokerName]; !existed {
		controller.brokerMap[brokerName] = broker
	}
	if _, existed := controller.serviceMap[brokerName]; existed {
		glog.Infof("Service broker: %s has been registered already. \n", brokerName)
		return true
	}
	// store service broker's Catalog() info
	catalog, err := broker.Catalog()
	if err != nil {
		glog.Errorf("Failed to get Catalog information of service broker: %s", brokerName)
		return false
	}
	controller.serviceMap[brokerName] = catalog
	glog.Infof("Success to register service broker %s", brokerName)
	return true
}

func (c *BaseController) GetCatalog() (*brokerapi.Catalog, error) {
	c.controllerMutex.Lock()
	defer c.controllerMutex.Unlock()
	var ServicesSlice []*brokerapi.Service
	for _, catalog := range c.serviceMap {
		for _, service := range catalog.Services {
			ServicesSlice = append(ServicesSlice, service)
		}
	}
	return &brokerapi.Catalog{Services: ServicesSlice}, nil
}

func (c *BaseController) GetServiceBroker(serviceID string) (brokerName string, broker brokerapi.ServiceBroker) {
	c.controllerMutex.Lock()
	defer c.controllerMutex.Unlock()
	brokerName = ""
	for name, catalog := range c.serviceMap {
		for _, service := range catalog.Services {
			if service.ID == serviceID {
				brokerName = name
				break
			}
		}
	}

	if brokerName == "" {
		return "", nil
	}

	broker, ok := c.brokerMap[brokerName]
	if !ok {
		return "", nil
	}
	return brokerName, broker
}

func (c *BaseController) getBrokerByName(brokerName string) brokerapi.ServiceBroker {
	c.controllerMutex.Lock()
	defer c.controllerMutex.Unlock()
	if broker, ok := c.brokerMap[brokerName]; ok {
		return broker
	}
	return nil
}

type BindRunInfo struct {
	BindingId string
	Status    string
	Parameter map[string]interface{}
}

type InstanceRunInfo struct {
	InstanceId string
	ServiceID  string
	PlanID     string
	Bindings   map[string]*BindRunInfo
	Status     string
	BrokerName string
	Parameter  map[string]interface{}
}

func mergeParameter(parameterIn map[string]interface{},
	parameterOut map[string]interface{}) {
	for OutKey, OutValue := range parameterOut {
		for InKey := range parameterIn {
			if InKey == OutKey {
				parameterIn[InKey] = OutValue
			}
		}
		parameterIn[OutKey] = OutValue
	}
	return
}

func (c *BaseController) CreateServiceInstance(instanceId, serviceID, planID string,
	parameterIn map[string]interface{}) error {
	instanceGet := c.asyncEngine.GetAsyncInstance(instanceId)
	if instanceGet != nil {
		switch instanceGet.Status {
		case StateProvisionInstanceInProgress:
			return fmt.Errorf("Duplicated provisioning request of instance %s. It is in progress.", instanceId)
		case StateProvisionInstanceFailed:
			return fmt.Errorf("Duplicated provisioning request of instance %s. It was failed", instanceId)
		case StateProvisionInstanceSucceeded:
			return fmt.Errorf("Duplicated provisioning request of instance %s. It was succeeded", instanceId)
		}
		return fmt.Errorf("Duplicated provisioning request of instance. %v", instanceGet)
	}
	brokerName, broker := c.GetServiceBroker(serviceID)
	instanceInfo := new(InstanceRunInfo)
	instanceInfo.InstanceId = instanceId
	instanceInfo.ServiceID = serviceID
	instanceInfo.PlanID = planID
	instanceInfo.Bindings = make(map[string]*BindRunInfo)
	instanceInfo.Status = StateProvisionInstanceInProgress
	instanceInfo.BrokerName = brokerName
	instanceInfo.Parameter = parameterIn

	err := c.asyncEngine.AddAsyncInstance(instanceInfo, true)
	if err != nil {
		glog.Infof("CreateServiceInstance AddAsyncInstance has error:%v", err)
		return err
	}
	parameterOut, err := broker.Provision(instanceInfo.InstanceId,
		instanceInfo.ServiceID, instanceInfo.PlanID, instanceInfo.Parameter)
	if err != nil {
		glog.Infof("CreateServiceInstance broker provision instance gets error:%v", err)
		// update InstanceRunInfo.Status to StateProvisionInstanceFailed,
		// stop asyncEngine continuously check instance provision status
		instanceInfo.Status = StateProvisionInstanceFailed
		err = c.asyncEngine.UpdateAsyncInstance(instanceInfo)
		if err != nil {
			glog.Infof("Provision instance failed, while UpdateAsyncInstance has error: %v", err)
		}
		return err
	}

	if parameterOut != nil {
		mergeParameter(parameterIn, parameterOut)
		err = c.asyncEngine.UpdateAsyncInstance(instanceInfo)
		if err != nil {
			glog.Infof("UpdateAsyncInstance err: %v", err)
			return err
		}
	}

	glog.Infof("CreateServiceInstance finished.")
	return nil
}

func (c *BaseController) DeleteServiceInstance(instanceId, serviceID,
	planID string, parameterIn map[string]string) error {
	instanceInfo := c.asyncEngine.GetAsyncInstance(instanceId)
	if instanceInfo == nil {
		glog.Infof("Not found service instance %s record in async engine when delete.", instanceId)
		// InstanceId not found in async engine, but exists in catalog api server.
		// it might because of data inconsistent between async engine and catalog.
		// Try to deprovision backend instance anyway, while ignore checking related binding.
		// Since it's not able to find out binding info in async engine anymore.
		_, broker := c.GetServiceBroker(serviceID)
		err := broker.Deprovision(instanceId, "", "", nil)
		if err != nil {
			glog.Infof("Broker %v deprovision instance %s failed.", broker, instanceId)
			return err
		}
	} else {
		glog.Infof("Found service instance in async engine: %v with status %v", instanceInfo.InstanceId, instanceInfo.Status)
		if len(instanceInfo.Bindings) != 0 {
			glog.Infof("There are active bindings of service instance %s, please delete first.", instanceInfo.InstanceId)
			return fmt.Errorf("There are active bindings of Instance %s, please delete first.", instanceInfo.InstanceId)
		}
		if instanceInfo.Status != StateProvisionInstanceSucceeded {
			glog.Infof("Instance %v status is %v. It should not be deprovisioned.", instanceInfo.InstanceId, instanceInfo.Status)
			return fmt.Errorf("Instance %v status is %v. It should not be deprovisioned.", instanceInfo.InstanceId, instanceInfo.Status)
		}
		_, broker := c.GetServiceBroker(serviceID)
		err := broker.Deprovision(instanceInfo.InstanceId,
			instanceInfo.ServiceID, instanceInfo.PlanID, instanceInfo.Parameter)
		if err != nil {
			glog.Infof("Broker %v deprovision instance %s failed.", broker, instanceInfo.InstanceId)
			return err
		}
		instanceInfo.Status = StateDeprovisionInstanceInProgress
		err = c.asyncEngine.DeleteAsyncInstance(instanceInfo.InstanceId)
		if err != nil {
			glog.Infof("Async engine failed to DeleteAsyncInstance %s, with error %v.", instanceInfo.InstanceId, err)
			return err
		}
	}
	glog.Infof("Delete service instance %s success!", instanceInfo.InstanceId)
	return nil
}

func (c *BaseController) GetServiceInstanceStatus(instanceId, serviceID,
	planID string, operation string) (string, error) {
	glog.Infof("GetServiceInstanceStatus in.")
	instanceInfo := c.asyncEngine.GetAsyncInstance(instanceId)
	if instanceInfo == nil {
		glog.Infof("GetServiceInstanceStatus failed.")
		return "", fmt.Errorf("Not found instance %s when get instance status.", instanceId)
	}
	if instanceInfo.Status == StateProvisionInstanceInProgress {
		c.asyncEngine.GetProvisionInstanceLastOperation(instanceInfo)
	} else if instanceInfo.Status == StateDeprovisionInstanceInProgress {
		c.asyncEngine.GetDeprovisionInstanceLastOperation(instanceInfo)
	}
	glog.Infof("GetServiceInstanceStatus out.")
	return instanceInfo.Status, nil
}

func (c *BaseController) CreateServiceBinding(instanceId, serviceID,
	planID, bindingId string, parameterIn map[string]interface{}) (brokerapi.Credential, error) {
	glog.Infof("CreateServiceBinding in.")
	instanceInfo := c.asyncEngine.GetAsyncInstance(instanceId)
	if instanceInfo == nil {
		glog.Infof("Not found instance %s when create binding %s.", instanceId, bindingId)
		return nil, fmt.Errorf("Not found instance %s when create binding %s.", instanceId, bindingId)
	}

	bindingInfo := c.asyncEngine.GetAsyncBinding(instanceId, bindingId)
	if bindingInfo != nil {
		glog.Infof("Binding already exist instance %s when create binding %s.", instanceId, bindingId)
		return nil, fmt.Errorf("Binding already exist instance %s when create binding %s.", instanceId, bindingId)
	}

	if instanceInfo.Status != StateProvisionInstanceSucceeded {
		glog.Infof("Instance status is wrong %s when create binding %s.", instanceId, bindingId)
		return nil, fmt.Errorf("Instance status is wrong %s when create binding %s.", instanceId, bindingId)
	}

	if _, ok := instanceInfo.Bindings[bindingId]; ok {
		glog.Infof("The Instance %s has already been binded to bind %s.", instanceId, bindingId)
		return nil, fmt.Errorf("The Instance %s has already been binded to bind %s.", instanceId, bindingId)
	}

	bindInfo := BindRunInfo{BindingId: bindingId,
		Status:    StateBindingInProgress,
		Parameter: parameterIn}

	instanceInfo.Bindings[bindingId] = &bindInfo

	err := c.asyncEngine.UpdateAsyncInstance(instanceInfo)
	if err != nil {
		glog.Infof("Update Instance %s binding %s failed.", instanceId, bindingId)
		return nil, err
	}

	broker := c.getBrokerByName(instanceInfo.BrokerName)
	if broker == nil {
		glog.Infof("Broker %s is not found!", instanceInfo.BrokerName)
		return nil, fmt.Errorf("Broker %s is not found!", instanceInfo.BrokerName)
	}
	parameterOut, credential, err := broker.Bind(instanceInfo.InstanceId, instanceInfo.ServiceID,
		instanceInfo.PlanID, bindInfo.BindingId, instanceInfo.Parameter, bindInfo.Parameter)
	if err != nil {
		glog.Infof("Broker %s failed!", instanceInfo.BrokerName)
		return nil, err
	}

	if parameterOut != nil {
		mergeParameter(parameterIn, parameterOut)
		err = c.asyncEngine.UpdateAsyncInstance(instanceInfo)
		if err != nil {
			glog.Infof("Update Instance %s binding %s failed second.", instanceId, bindingId)
			return nil, err
		}
	}
	glog.Infof("CreateServiceBinding out.")
	return credential, nil
}

func (c *BaseController) DeleteServiceBinding(instanceId, serviceID,
	planID, bindingId string, parameterIn map[string]string) error {
	instanceInfo := c.asyncEngine.GetAsyncInstance(instanceId)
	if instanceInfo == nil {
		return fmt.Errorf("Not found instance %s when delete binding %s.", instanceId, bindingId)
	}

	if instanceInfo.Status != StateProvisionInstanceSucceeded {
		return fmt.Errorf("Instance status is wrong %s when delete binding %s.", instanceId, bindingId)
	}

	bindInfo, ok := instanceInfo.Bindings[bindingId]
	if !ok {
		return fmt.Errorf("The binding %s of instance %s has already been deleted.", bindingId, instanceId)
	}

	delete(instanceInfo.Bindings, bindingId)

	bindInfo.Status = StateUnBindingInProgress

	broker := c.getBrokerByName(instanceInfo.BrokerName)
	if broker == nil {
		return fmt.Errorf("Broker %s is not found!", instanceInfo.BrokerName)
	}
	err := broker.UnBind(instanceInfo.InstanceId, instanceInfo.ServiceID,
		instanceInfo.PlanID, bindingId, instanceInfo.Parameter, bindInfo.Parameter)
	if err != nil {
		return err
	}

	err = c.asyncEngine.UpdateAsyncInstance(instanceInfo)
	if err != nil {
		return err
	}
	glog.Infof("DeleteServiceBinding success!")
	return nil
}

func (c *BaseController) GetServiceBindingStatus(instanceId, serviceID, planID, bindingId string, operation string) (string, error) {
	instanceInfo := c.asyncEngine.GetAsyncInstance(instanceId)
	if instanceInfo == nil {
		return "", fmt.Errorf("Not found instance %s when get binding status.", instanceId)
	}

	if instanceInfo.Status != StateProvisionInstanceSucceeded {
		return instanceInfo.Status, fmt.Errorf("Instance status is wrong %s when get binding status.", instanceInfo.InstanceId)
	}

	bindInfo, ok := instanceInfo.Bindings[bindingId]
	if !ok {
		return "", fmt.Errorf("The binding %s of instance %s is not found.", bindingId, instanceId)
	}

	if bindInfo.Status == StateBindingInProgress {
		c.asyncEngine.GetBindingLastOperation(instanceInfo, bindInfo)
	} else if bindInfo.Status == StateUnBindingInProgress {
		c.asyncEngine.GetUnBindingLastOperation(instanceInfo, bindInfo)
	}

	return bindInfo.Status, nil
}
