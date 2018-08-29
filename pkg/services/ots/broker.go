package ots

import (
	"errors"
	"fmt"
	"github.com/AliyunContainerService/open-service-broker-alibabacloud/pkg/brokerapi"
	"github.com/golang/glog"
)

type errNoSuchInstance struct {
	instanceID string
}

func (e errNoSuchInstance) Error() string {
	return fmt.Sprintf("no such instance with ID %s", e.instanceID)
}

type OTSBroker struct {
	region          string
	zoneID          string
	vpcID           string
	vSwitchID       string
	accessKeyId     string
	accessKeySecret string
}

func CreateBroker() brokerapi.ServiceBroker {
	otsBroker := createBrokerImpl()
	return otsBroker
}

func (c *OTSBroker) Catalog() (*brokerapi.Catalog, error) {
	return getCatalog(), nil
}

func (c *OTSBroker) Provision(instanceID, serviceID, planID string, parameterIn map[string]interface{}) (map[string]interface{}, error) {

	glog.Infof("Created table store instance:\n%v\n", instanceID)

	servicePlanMetadata, err := findServicePlanMetadata(serviceID, planID, parameterIn)
	if err != nil {
		glog.Warningln(err)
		return nil, err
	}
	glog.Infof("Get servicePlan info success serviceplan:\n%v\n", servicePlanMetadata)

	otsUrl, err := CreateOtsInstance(instanceID, c.region)
	if err != nil {
		glog.Infof("Create OTS  instance failed:\n%v\n", instanceID, err)
		return nil, err
	}

	err = CreateRamForInstance(instanceID, c.region)
	if err != nil {
		glog.Infof("Create Ram Info for instance %s failed:\n%v\n", instanceID, err)
		err = DeleteOtsInstance(instanceID, c.region)
		if err != nil {
			glog.Infof("Delete OTS instance %s failed, please delete manually.error: %v\n", instanceID, err)
			return nil, err
		}
		glog.Infof("Delete OTS  instance %s success.", instanceID)
		return nil, err
	}

	parameterIn["otsUrl"] = otsUrl

	glog.Infof("Create OTS instance %s success.")

	return parameterIn, nil
}

func (c *OTSBroker) GetInstanceStatus(instanceID, serviceID, planID string,
	parameterIn map[string]interface{}) (bool, error) {

	otsOk, _ := CheckOtsInstanceStatus(instanceID, c.region)
	if otsOk == false {
		glog.Infof("Oss of instance %s is not ready", instanceID)
		return false, nil
	}

	ramOk := CheckRamForInstance(instanceID, c.region)
	if ramOk == false {
		glog.Infof("Ram of instance %s is not ready", instanceID)
		return false, nil
	}

	glog.Infof("Service Instance %s lastoperation is :%s\n", instanceID)
	return true, nil
}

func (c *OTSBroker) GetServiceInstance(id string) (string, error) {
	return "", errors.New("Unimplemented")
}

func (c *OTSBroker) Deprovision(instanceID, serviceID, planID string, parameterIn map[string]interface{}) error {

	err := DeleteOtsInstance(instanceID, c.region)
	if err != nil {
		glog.Infof("Delete OTS instance %s failed error: %v\n", instanceID, err)
		return err
	}

	err = DeleteRamForInstance(instanceID, c.region)
	if err != nil {
		glog.Infof("Delete RAM for instance %s failed error: %v\n", instanceID, err)
		return err
	}

	return nil
}

func (c *OTSBroker) Bind(instanceID, serviceID, planID, bindingID string,
	parameterInstance, parameterIn map[string]interface{}) (map[string]interface{}, brokerapi.Credential, error) {

	user, accessKey, err := CreateRamUserForBinding(bindingID, instanceID, c.region)
	if err != nil {
		glog.Warningln(err)
		return nil, nil, err
	}

	parameterIn["userAccessKeyID"] = accessKey.AccessKeyId

	glog.Infof("create ots binding begin to return response.")
	return parameterIn, brokerapi.Credential{
		"instancename":      GetOtsNameForInstance(instanceID),
		"username":          user.UserName,
		"access_key_id":     accessKey.AccessKeyId,
		"secret_access_key": accessKey.AccessKeySecret,
		"host":              GetOTSInstanceUrl(instanceID, c.region),
	}, nil
}

func (c *OTSBroker) UnBind(instanceID, serviceID, planID, bindingID string,
	parameterInstance, parameterIn map[string]interface{}) error {

	userAccessKeyID := ""
	vUserAccessKeyID, ok := parameterIn["userAccessKeyID"]
	if ok {
		if value, ok := vUserAccessKeyID.(string); ok {
			userAccessKeyID = value
		} else {
			err := fmt.Errorf("UnBind can not get userAccessKeyID for instance %s binding :%s.", instanceID, bindingID)
			return err
		}
	}

	err := DeleteRamUserForBinding(bindingID, instanceID, c.region, userAccessKeyID)
	if err != nil {
		glog.Infof("UnBind faield to delete bindingID %s. error:%v", bindingID, err)
		return err
	}

	return nil
}

func (c *OTSBroker) GetBindingStatus(instanceID, serviceID, planID, bindingID string,
	parameterInstance, parameterIn map[string]interface{}) (bool, error) {
	err := GetRamUserForBinding(bindingID, instanceID, c.region)
	if err != nil {
		return false, err
	}
	return true, nil
}

func createBrokerImpl() *OTSBroker {

	userMetaData, err := GetCloudServiceMetaData()
	if err != nil {
		glog.Infof("Failed to get user metadate err:%v\n", err)
	}

	return &OTSBroker{
		region:    userMetaData.Region,
		zoneID:    userMetaData.ZoneID,
		vpcID:     userMetaData.VpcID,
		vSwitchID: userMetaData.VSwitchID,
	}
}
