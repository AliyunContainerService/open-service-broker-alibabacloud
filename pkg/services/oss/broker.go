package oss

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

type OSSBroker struct {
	region          string
	zoneID          string
	vpcID           string
	vSwitchID       string
	accessKeyId     string
	accessKeySecret string
}

func CreateBroker() brokerapi.ServiceBroker {
	ossBroker := createBrokerImpl()
	if ossBroker == nil {
		return nil
	}
	glog.Infof("OSSBroker is created: %v", ossBroker)
	return ossBroker
}

func (c *OSSBroker) Catalog() (*brokerapi.Catalog, error) {
	return getCatalog(), nil
}

func (c *OSSBroker) Provision(instanceID, serviceID, planID string, parameterIn map[string]interface{}) (map[string]interface{}, error) {
	glog.Infof("Creating OSS Service Instance:%s", instanceID)

	accessKeyId, accessKeySecret, region, acl, storageClass, err := GetAccessInfoFromInstanceParameters(parameterIn)
	if err != nil {
		glog.Warningln(err)
		return nil, err
	}

	if accessKeyId == "" || accessKeySecret == "" {
		glog.Infof("Created Instance %s failed for no specify accessKeyId and accessKeySecret.:\n%v\n", instanceID)
		return nil, fmt.Errorf("Created Instance " + instanceID + " failed for no specify accessKeyId and accessKeySecret.")
	}
	if region == "" {
		region = c.region
	}

	//servicePlanMetadata, err := findServicePlanMetadata(serviceID, planID, parameterIn)
	//if err != nil {
	//	glog.Warningln(err)
	//	return nil, err
	//}
	//glog.Infof("Get servicePlan Info success serviceplan:%v", servicePlanMetadata)
	glog.Infof("Prepare to provision OSS bucket in region: %v", region)

	bucketName, err := CreateOssBucket(instanceID, region, accessKeyId, accessKeySecret, acl, storageClass)
	if err != nil {
		glog.Infof("Create OSS bucket for instance failed:\n%v\n", instanceID, err)
		return nil, err
	}

	err = CreateRamForInstance(instanceID, region, accessKeyId, accessKeySecret)
	if err != nil {
		glog.Infof("Create Ram Info for instance %s failed:\n%v\n", instanceID, err)
		err = DeleteOssBucket(instanceID, region, accessKeyId, accessKeySecret)
		if err != nil {
			glog.Infof("Delete OSS bucket for instance %s failed, please delete manually.error: %v\n", instanceID, err)
			return nil, err
		}
		glog.Infof("Delete OSS bucket for instance %s success.", instanceID)
		return nil, err
	}

	parameterIn["region"] = region
	parameterIn["accessKeyId"] = accessKeyId
	parameterIn["accessKeySecret"] = accessKeySecret

	glog.Infof("Create OSS bucket %s success.", bucketName)

	return parameterIn, nil
}

func (c *OSSBroker) GetInstanceStatus(instanceID, serviceID, planID string,
	parameterIn map[string]interface{}) (bool, error) {
	region, accessKeyId, accessKeySecret, err := c.getServiceInstanceInfo(parameterIn)
	if err != nil {
		err := fmt.Errorf("Get info (region,accessKeyId and accessKeySecret) for instance %s failed when get last operation.", instanceID)
		glog.Warningln(err)
		return false, err
	}

	ossOk, _ := CheckOssInstanceStatus(instanceID, region, accessKeyId, accessKeySecret)
	if ossOk == false {
		glog.Infof("Oss of instance %s is not ready", instanceID)
		return false, nil
	}

	ramOk := CheckRamForInstance(instanceID, region, accessKeyId, accessKeySecret)
	if ramOk == false {
		glog.Infof("Ram of instance %s is not ready", instanceID)
		return false, nil
	}

	glog.Infof("Service Instance %s lastoperation is :%s\n", instanceID)
	return true, nil
}

func (c *OSSBroker) GetServiceInstance(id string) (string, error) {
	return "", errors.New("Unimplemented")
}

func (c *OSSBroker) Deprovision(instanceID, serviceID, planID string, parameterIn map[string]interface{}) error {

	region, accessKeyId, accessKeySecret, err := c.getServiceInstanceInfo(parameterIn)
	if err != nil {
		err := fmt.Errorf("Get info (region,accessKeyId and accessKeySecret) for instance %s failed when remove instance.", instanceID)
		glog.Warningln(err)
		return err
	}

	err = DeleteOssBucket(instanceID, region, accessKeyId, accessKeySecret)
	if err != nil {
		glog.Infof("Delete OSS bucket for instance %s failed error: %v\n", instanceID, err)
		return err
	}

	err = DeleteRamForInstance(instanceID, region, accessKeyId, accessKeySecret)
	if err != nil {
		glog.Infof("Delete RAM for instance %s failed error: %v\n", instanceID, err)
		return err
	}

	return nil
}

func (c *OSSBroker) Bind(instanceID, serviceID, planID, bindingID string,
	parameterInstance, parameterIn map[string]interface{}) (map[string]interface{}, brokerapi.Credential, error) {

	region, accessKeyId, accessKeySecret, err := c.getServiceInstanceInfo(parameterInstance)
	if err != nil {
		err := fmt.Errorf("Get info (region,accessKeyId and accessKeySecret) for instance %s failed when create binding.", instanceID)
		glog.Warningln(err)
		return nil, nil, err
	}

	user, accessKey, err := CreateRamUserForBinding(bindingID, instanceID, region, accessKeyId, accessKeySecret)
	if err != nil {
		glog.Warningln(err)
		return nil, nil, err
	}

	parameterIn["userAccessKeyID"] = accessKey.AccessKeyId

	glog.Infof("createServiceBindingImpl begin to return response.")
	return parameterIn, brokerapi.Credential{
		"bucket":            GetBucketNameForInstance(instanceID),
		"username":          user.UserName,
		"access_key_id":     accessKey.AccessKeyId,
		"secret_access_key": accessKey.AccessKeySecret,
		"host":              GetOSSEndPoint(region),
		"uri": "oss://" + accessKey.AccessKeyId + ":" + accessKey.AccessKeySecret + "@" +
			GetOSSEndPoint(c.region) + ":" + GetBucketNameForInstance(instanceID),
	}, nil
}

func (c *OSSBroker) UnBind(instanceID, serviceID, planID, bindingID string,
	parameterInstance, parameterIn map[string]interface{}) error {
	region, accessKeyId, accessKeySecret, err := c.getServiceInstanceInfo(parameterInstance)
	if err != nil {
		err := fmt.Errorf("Get info (region,accessKeyId and accessKeySecret) for instance %s failed when remove instance.", instanceID)
		glog.Warningln(err)
		return err
	}

	userAccessKeyID := ""
	vUserAccessKeyID, ok := parameterIn["userAccessKeyID"]
	if ok {
		if value, ok := vUserAccessKeyID.(string); ok {
			userAccessKeyID = value
		} else {
			err = fmt.Errorf("UnBind can not get userAccessKeyID for instance %s binding :%s.", instanceID, bindingID)
			return err
		}
	}

	err = DeleteRamUserForBinding(bindingID, instanceID, region, userAccessKeyID, accessKeyId, accessKeySecret)
	if err != nil {
		glog.Infof("UnBind faield to delete bindingID %s. error:%v", bindingID, err)
		return err
	}

	return nil
}

func (c *OSSBroker) GetBindingStatus(instanceID, serviceID, planID, bindingID string,
	parameterInstance, parameterIn map[string]interface{}) (bool, error) {
	region, accessKeyId, accessKeySecret, err := c.getServiceInstanceInfo(parameterInstance)
	if err != nil {
		err := fmt.Errorf("Get info (region,accessKeyId and accessKeySecret) for instance %s failed when remove instance.", instanceID)
		glog.Warningln(err)
		return false, err
	}
	err = GetRamUserForBinding(bindingID, instanceID, region, accessKeyId, accessKeySecret)
	if err != nil {
		return false, err
	}
	return true, nil
}

func createBrokerImpl() *OSSBroker {

	userMetaData, err := GetCloudServiceMetaData()
	if err != nil {
		glog.Infof("Failed to get cloud service metadate.")
		return nil
	}

	return &OSSBroker{
		region:    userMetaData.Region,
		zoneID:    userMetaData.ZoneID,
		vpcID:     userMetaData.VpcID,
		vSwitchID: userMetaData.VSwitchID,
	}
}

func (c *OSSBroker) getServiceInstanceInfo(parameterIn map[string]interface{}) (region, accessKeyId, accessKeySecret string, err error) {
	glog.Infof("OSS broker getServiceInstanceInfo with input param: %v", parameterIn)
	vRegion, ok := parameterIn["region"]
	if ok {
		if value, ok := vRegion.(string); ok {
			region = value
		} else {
			err = fmt.Errorf("Service instance AccessKeyId  type is wrong suppose to be string.")
			return
		}
	}

	vAccessKeyId, ok := parameterIn["accessKeyId"]
	if ok {
		if value, ok := vAccessKeyId.(string); ok {
			accessKeyId = value
		} else {
			err = fmt.Errorf("Service instance AccessKeyId  type is wrong suppose to be string.")
			return
		}
	}

	vAccessKeySecret, ok := parameterIn["accessKeySecret"]
	if ok {
		if value, ok := vAccessKeySecret.(string); ok {
			accessKeySecret = value
		} else {
			err = fmt.Errorf("Service instance AccessKeyId  type is wrong suppose to be string.")
			return
		}
	}

	return
}

func getServiceInstanceNamespace(req *brokerapi.CreateServiceInstanceRequest) string {

	if req.ContextProfile.Namespace != "" {
		return req.ContextProfile.Namespace
	}
	return "default"
}
