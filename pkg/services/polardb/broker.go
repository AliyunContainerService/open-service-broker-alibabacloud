package polardb

import (
	"errors"
	"fmt"
	"github.com/AliyunContainerService/open-service-broker-alibabacloud/pkg/brokerapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/utils"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/golang/glog"
	"strings"
	"time"
)

var (
	poolSize         = 10
	maxTaskQueueSize = 10000

	defaultUsername = "myaccount"
	defaultPassword = "DB_paasw0rd"
)

type errNoSuchInstance struct {
	instanceID string
}

func (e errNoSuchInstance) Error() string {
	return fmt.Sprintf("no such instance with ID %s", e.instanceID)
}

type PolarDBBroker struct {
	client    *polardb.Client
}

func CreateBroker() brokerapi.ServiceBroker {
	polardbBroker := createBrokerImpl()
	if polardbBroker == nil {
		return nil
	}
	glog.Infof("PolarDBBroker is created: %v", polardbBroker)
	return polardbBroker
}

func (c *PolarDBBroker) Catalog() (*brokerapi.Catalog, error) {
	return getCatalog(), nil
}

func (c *PolarDBBroker) Provision(instanceID, serviceID, planID string, parameterIn map[string]interface{}) (map[string]interface{}, error) {
	glog.Infof("Created Database Service Instance:\n%v\n", instanceID)
	servicePlanMetadata, err := findServicePlanMetadata(serviceID, planID, parameterIn)
	if err != nil {
		glog.Warningln(err)
		return nil, err
	}
	glog.Infof("Get servicePlan Info success serviceplan:\n%v\n", servicePlanMetadata)

	if err := c.CreateNewClientFromAK(
		servicePlanMetadata.RegionId, servicePlanMetadata.AccessKeyId, servicePlanMetadata.AccessKeySecret); err != nil {
		glog.Infof("Create PolarDB client failed:\n%v\n", err)
		return nil, err
	}

	request := polardb.CreateCreateDBClusterRequest()
	request.RegionId = servicePlanMetadata.RegionId
	request.ZoneId = servicePlanMetadata.ZoneId
	request.DBClusterDescription = instanceID //Set service instance ID into description
	request.DBType = servicePlanMetadata.Engine
	request.DBVersion = servicePlanMetadata.EngineVersion

	request.DBNodeClass = servicePlanMetadata.Class
	request.ClientToken = utils.GetUUID()
	request.ClusterNetworkType = "VPC"
	request.VPCId = servicePlanMetadata.VpcID
	request.VSwitchId = servicePlanMetadata.VSwitchID
	request.PayType = "Postpaid"
	glog.Infof("Creating PolarDB instance with request: %v...", request)
	response, err := c.client.CreateDBCluster(request)
	if err != nil {
		glog.Warningln(err)
		return map[string]interface{}{"createSuccess": false}, err
	}
	glog.Infof("Backend PolarDB create instance got response: %v", response)
	// do add tags may fail, so do sleep .
	time.Sleep(3 * time.Second)

	// add service broker tags to RDS instance
	tagReq := polardb.CreateTagResourcesRequest()

	tagReq.ResourceId = &[]string{response.DBClusterId}
	tagReq.ResourceType = "cluster"
	tagReq.Tag = &[]polardb.TagResourcesTag{{Value: SERVICE_CATALOG_TAG_VALUE, Key: SERVICE_CATALOG_TAG_KEY}}

	_, err = c.client.TagResources(tagReq)
	if err != nil {
		glog.Infof("PolarDB instance %s is created, but gets error when adding tags.  %v", response.DBClusterId, err)
		return nil, err
	}
	return parameterIn, nil
}

func (c *PolarDBBroker) GetInstanceStatus(instanceID, serviceID, planID string,
	parameterIn map[string]interface{}) (bool, error) {
	servicePlanMetadata, err := findServicePlanMetadata(serviceID, planID, parameterIn)
	if err != nil {
		glog.Warningln(err)
		return false, err
	}
	if err := c.CreateNewClientFromAK(
		servicePlanMetadata.RegionId, servicePlanMetadata.AccessKeyId, servicePlanMetadata.AccessKeySecret); err != nil {
		glog.Infof("Create Rds client failed:\n%v\n", err)
		return false, err
	}
	return c.CheckDBClusterStatus(instanceID, "Running", parameterIn)
}

func (c *PolarDBBroker) GetServiceInstance(id string) (string, error) {
	return "", errors.New("Unimplemented")
}

func (c *PolarDBBroker) Deprovision(instanceID, serviceID, planID string, parameterIn map[string]interface{}) error {
	if parameterIn == nil {
		return nil
	}
	servicePlanMetadata, err := findServicePlanMetadata(serviceID, planID, parameterIn)
	if err != nil {
		glog.Warningln(err)
		return err
	}
	if err := c.CreateNewClientFromAK(
		servicePlanMetadata.RegionId, servicePlanMetadata.AccessKeyId, servicePlanMetadata.AccessKeySecret); err != nil {
		glog.Infof("Create Rds client failed:\n%v\n", err)
		return err
	}
	dbInstanceID, err := c.getDBClusterByID(instanceID)
	if err != nil {
		glog.Infof("Failed to found PolarDB instance with description %s. Gets error %v",instanceID, err)
		return err
	}
	if dbInstanceID == "" {
		glog.Infof("Not found PolarDB instance with description %s", instanceID)
		return nil
	}
	req := polardb.CreateDeleteDBClusterRequest()
	req.DBClusterId = dbInstanceID
	glog.Infof("Deleting instance %v: [polardb %v %v].", instanceID, dbInstanceID, req)
	_, err = c.client.DeleteDBCluster(req)
	if err != nil {
		glog.Infof("Delete instance %v: [polardb %v] get error: %v.", instanceID, dbInstanceID, err)
		return err
	}

	return nil
}

func (c *PolarDBBroker) Bind(instanceID, serviceID, planID, bindingID string,
	parameterInstance, parameterIn map[string]interface{}) (map[string]interface{}, brokerapi.Credential, error) {

	servicePlanMetadata, err := findServicePlanMetadata(serviceID, planID, parameterInstance)
	if err != nil {
		glog.Warningln(err)
		return nil, nil, err
	}
	if err := c.CreateNewClientFromAK(
		servicePlanMetadata.RegionId, servicePlanMetadata.AccessKeyId, servicePlanMetadata.AccessKeySecret); err != nil {
		glog.Infof("Create Rds client failed:\n%v\n", err)
		return nil, nil, err
	}

	dbInstanceID, err := c.getDBClusterByID(instanceID)
	if err != nil {
		glog.Warningln(err)
		return nil, nil, err
	}

	accountName, accountPassword := GetAccountInfoFromBindingParameters(parameterIn)

	err = c.CreateDatabaseAccount(instanceID, dbInstanceID, accountName, accountPassword)
	if err != nil {
		if !strings.Contains(err.Error(), "InvalidAccountName.Duplicate") {
			glog.Warningln(err)
			return nil, nil, err
		} else {
			glog.Infof("%v %v %v InvalidAccountName.Duplicate, ignore!", instanceID, dbInstanceID, accountName)
		}
		//err = c.ResetDatabaseAccount(instanceID, dbInstanceID, accountName, accountPassword)
		//if err != nil {
		//	glog.Warningln(err)
		//	return nil, nil, err
		//}
	}

	if err = c.SetDatabaseSecurityIps(instanceID, dbInstanceID, servicePlanMetadata.SecurityIps); err != nil {
		glog.Error(err)
		return nil, nil, err
	}

	info, err := c.getDBClusterConnectionInfo(dbInstanceID)
	if err != nil {
		glog.Warningln(err)
		return nil, nil, err
	}

	host := info.ConnectionString
	port := info.Port
	glog.Infof("createServiceBindingImpl begin to return response.")
	return parameterIn, brokerapi.Credential{
		"uri":      "mysql://" + accountName + ":" + accountPassword + "@" + host + ":" + port + "/",
		"username": accountName,
		"password": accountPassword,
		"host":     host,
		"port":     port,
	}, nil
}

func (c *PolarDBBroker) UnBind(instanceID, serviceID, planID, bindingID string,
	parameterInstance, parameterIn map[string]interface{}) error {
	servicePlanMetadata, err := findServicePlanMetadata(serviceID, planID, parameterInstance)
	if err != nil {
		glog.Warningln(err)
		return err
	}
	if err := c.CreateNewClientFromAK(
		servicePlanMetadata.RegionId, servicePlanMetadata.AccessKeyId, servicePlanMetadata.AccessKeySecret); err != nil {
		glog.Infof("Create Rds client failed:\n%v\n", err)
		return err
	}

	//dbInstanceID, err := c.getDBClusterByID(instanceID)
	//if err != nil {
	//	glog.Infof("UnBind faield for not found instance %s DBinstance. error:%v", instanceID, err)
	//	return err
	//}

	//accountName, _ := GetAccountInfoFromBindingParameters(parameterIn)
	//glog.Infof("UnBind instance %s DBinstance %s's account name:%s.", instanceID, dbInstanceID, accountName)
	//resetPassword := uniuri.New()
	//err = c.ResetDatabaseAccount(instanceID, dbInstanceID, accountName, resetPassword)
	//if err != nil {
	//	glog.Warningln(err)
	//	return err
	//}
	return nil
}

func (c *PolarDBBroker) GetBindingStatus(instanceID, serviceID, planID, bindingID string,
	parameterInstance, parameterIn map[string]interface{}) (bool, error) {
	accountName, _ := GetAccountInfoFromBindingParameters(parameterIn)
	return c.CheckAccountStatus(instanceID, accountName, "Available")
}

func createBrokerImpl() *PolarDBBroker {
	return &PolarDBBroker{
		client:    nil,
	}
}

const SERVICE_CATALOG_TAG_KEY = "service_catalog"
const SERVICE_CATALOG_TAG_VALUE = "true"
const DEFAULT_PAGE_SIZE = 50

func getNextPageNumber(pageNumber, pageSize, totalCount int) int {

	if pageNumber*pageSize > totalCount {
		return 0
	} else {
		return pageNumber + 1
	}

}

func (c *PolarDBBroker) SetDatabaseSecurityIps(instanceID, dbInstanceID, securityIps string) error {
	request := polardb.CreateModifyDBClusterAccessWhitelistRequest()
	request.DBClusterId = dbInstanceID
	request.SecurityIps = securityIps

	glog.Infof("PolarDB CreateDatabaseAccount DBinstance %s's account.", dbInstanceID)

	response, err := c.client.ModifyDBClusterAccessWhitelist(request)
	if err != nil {
		glog.Infof("Failed to ModifyDBClusterAccessWhitelist:%v\n", err)
		return err
	}

	glog.Infof("ModifyDBClusterAccessWhitelist DBinstance %s's response:%v.", dbInstanceID, response)
	return nil
}

func (c *PolarDBBroker) CreateDatabaseAccount(instanceID, dbInstanceID, accountName, accountPassword string) error {
	createPolarDBAccountRequest := polardb.CreateCreateAccountRequest()
	createPolarDBAccountRequest.DBClusterId = dbInstanceID
	createPolarDBAccountRequest.AccountName = accountName
	createPolarDBAccountRequest.AccountPassword = accountPassword
	createPolarDBAccountRequest.AccountType = "Super"

	glog.Infof("PolarDB CreateDatabaseAccount DBinstance %s's account.", dbInstanceID)

	response, err := c.client.CreateAccount(createPolarDBAccountRequest)
	if err != nil {
		glog.Infof("Failed to CreateAccount:%v\n", err)
		return err
	}

	glog.Infof("CreateDatabaseAccount DBinstance %s's account response:%v.", dbInstanceID, response)
	return nil
}

func (c *PolarDBBroker) ResetDatabaseAccount(instanceID, dbInstanceID, accountName, password string) error {
	resetAccountPasswordRequest := polardb.CreateResetAccountRequest()
	resetAccountPasswordRequest.AccountName = accountName
	resetAccountPasswordRequest.AccountPassword = password
	resetAccountPasswordRequest.DBClusterId = dbInstanceID

	response, err := c.client.ResetAccount(resetAccountPasswordRequest)

	if err != nil {
		glog.Infof("Failed to reset password %s :%v\n", password, err)
		return err
	}
	glog.Infof("Reset database account success \n%v\n", response)
	return nil
}

func handleErrorResponse(resp *responses.BaseResponse, err error) error {
	if err == nil {
		if !resp.IsSuccess() {
			err = fmt.Errorf("response status:%d error: %s", resp.GetHttpStatus(), resp.GetHttpContentString())
		}
	}
	if err != nil {
		glog.Warningln(err)
	}
	return err
}

func (c *PolarDBBroker) CreateNewClientFromAK(region, accessId, accessSecret string) error {
	client, err := polardb.NewClientWithAccessKey(region, accessId, accessSecret)
	if err != nil {
		glog.Infof("Failed to create rds client err:%v\n", err)
		return err
	}
	client.EnableAsync(poolSize, maxTaskQueueSize)
	c.client = client
	return nil
}

func getServiceInstanceNamespace(req *brokerapi.CreateServiceInstanceRequest) string {

	if req.ContextProfile.Namespace != "" {
		return req.ContextProfile.Namespace
	}
	return "default"
}