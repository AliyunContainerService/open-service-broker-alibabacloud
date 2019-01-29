package rds

import (
	"errors"
	"fmt"
	"github.com/AliyunContainerService/open-service-broker-alibabacloud/pkg/brokerapi"
	"github.com/AliyunContainerService/open-service-broker-alibabacloud/pkg/util"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/utils"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/dchest/uniuri"
	"github.com/golang/glog"
	"strconv"
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

type RDSBroker struct {
	client    *rds.Client
	region    string
	zoneID    string
	vpcID     string
	vSwitchID string
}

func CreateBroker() brokerapi.ServiceBroker {
	rdsBroker := createBrokerImpl()
	if rdsBroker == nil {
		return nil
	}
	glog.Infof("RDSBroker is created: %v", rdsBroker)
	return rdsBroker
}

func (c *RDSBroker) Catalog() (*brokerapi.Catalog, error) {
	return getCatalog(), nil
}

func (c *RDSBroker) Provision(instanceID, serviceID, planID string, parameterIn map[string]interface{}) (map[string]interface{}, error) {
	glog.Infof("Created Database Service Instance:\n%v\n", instanceID)

	servicePlanMetadata, err := findServicePlanMetadata(serviceID, planID, parameterIn)
	if err != nil {
		glog.Warningln(err)
		return nil, err
	}
	glog.Infof("Get servicePlan Info success serviceplan:\n%v\n", servicePlanMetadata)

	/*	if servicePlanMetadata.Engine != "MySQL" {
			err := fmt.Errorf("Now only support mysql engine for service id: %s, plan id: %s engine type:%s",
				serviceID, planID, servicePlanMetadata.Engine)
			glog.Warningln(err)
			return nil, err
		}
	*/
	if err := c.CreateNewClientFromStsToken(); err != nil {
		glog.Infof("Create Rds client failed:\n%v\n", err)
		return nil, err
	}

	vSwitchID, zone, err := c.checkVSwitch(c.vSwitchID, c.zoneID)
	if err != nil {
		glog.Infof("There is no available rds in region %s zone %s:%v\n", c.region, c.zoneID, err.Error())
		return nil, err
	}

	request := rds.CreateCreateDBInstanceRequest()

	request.RegionId = c.region
	request.ZoneId = zone
	request.DBInstanceDescription = instanceID //Set service instance ID into description
	request.Engine = servicePlanMetadata.Engine
	request.EngineVersion = servicePlanMetadata.EngineVersion
	if storage, err := strconv.Atoi(servicePlanMetadata.Storage); err == nil {
		request.DBInstanceStorage = requests.NewInteger(storage)
	} else {
		glog.Infof("Get DBInstanceStorage failed:\n%v\n", err)
		return nil, err
	}
	request.DBInstanceClass = servicePlanMetadata.Class
	request.ClientToken = utils.GetUUIDV4()
	if servicePlanMetadata.NetworkType == "VPC" {
		request.InstanceNetworkType = servicePlanMetadata.NetworkType
		if servicePlanMetadata.VpcID != "" {
			request.VPCId = servicePlanMetadata.VpcID
		} else {
			request.VPCId = c.vpcID
		}
		if servicePlanMetadata.VSwitchID != "" {
			request.VSwitchId = servicePlanMetadata.VSwitchID
		} else {
			request.VSwitchId = vSwitchID
		}
	}

	request.DBInstanceNetType = "Intranet"
	request.SecurityIPList = "0.0.0.0/0"
	request.PayType = "Postpaid"

	glog.Infof("Creating DB instance for service instance %s request:%v...", instanceID, request)

	response, err := c.client.CreateDBInstance(request)

	if err != nil {
		glog.Warningln(err)
		return nil, err
	}

	dbInstanceID := response.DBInstanceId

	tagReq := rds.CreateAddTagsToResourceRequest()
	tagReq.DBInstanceId = dbInstanceID
	tagReq.Tag1Key = SERVICE_CATALOG_TAG_KEY
	tagReq.Tag1Value = SERVICE_CATALOG_TAG_VALUE

	glog.Infof("Add tags to DB instance %s got response:%v", dbInstanceID, response)

	if err := c.CreateNewClientFromStsToken(); err != nil {
		glog.Infof("Create Rds client failed:\n%v\n", err)
		return nil, err
	}
	_, err = c.client.AddTagsToResource(tagReq)
	if err != nil {
		glog.Warningln(err)
		return nil, err
	}

	return parameterIn, nil
}

func (c *RDSBroker) GetInstanceStatus(instanceID, serviceID, planID string,
	parameterIn map[string]interface{}) (bool, error) {
	return c.CheckDBInstanceStatus(instanceID, "Running")
}

func (c *RDSBroker) GetServiceInstance(id string) (string, error) {
	return "", errors.New("Unimplemented")
}

func (c *RDSBroker) Deprovision(instanceID, serviceID, planID string, parameterIn map[string]interface{}) error {
	dbInstanceID, err := c.getDBInstanceByID(instanceID)
	if err != nil {
		glog.Warningln(err)
		return err
	}

	if err := c.CreateNewClientFromStsToken(); err != nil {
		glog.Infof("Create Rds client failed:\n%v\n", err)
		return err
	}

	req := rds.CreateDeleteDBInstanceRequest()
	req.DBInstanceId = dbInstanceID

	glog.Infof("Remove  instance %s DBinstance %s request:%v.", instanceID, dbInstanceID, req)

	_, err = c.client.DeleteDBInstance(req)

	if err != nil {
		glog.Warningln(err)
		return err
	}

	return nil
}

func (c *RDSBroker) Bind(instanceID, serviceID, planID, bindingID string,
	parameterInstance, parameterIn map[string]interface{}) (map[string]interface{}, brokerapi.Credential, error) {

	dbInstanceID, err := c.getDBInstanceByID(instanceID)
	if err != nil {
		glog.Warningln(err)
		return nil, nil, err
	}

	accountName, accountPassword := GetAccountInfoFromBindingParameters(parameterIn)

	err = c.CreateDatabaseAccount(instanceID, dbInstanceID, accountName, accountPassword)
	if err != nil {
		err = c.ResetDatabaseAccount(instanceID, dbInstanceID, accountName, accountPassword)
		if err != nil {
			glog.Warningln(err)
			return nil, nil, err
		}
	}

	info, err := c.getDBInstanceConnectionInfo(dbInstanceID)
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

func (c *RDSBroker) UnBind(instanceID, serviceID, planID, bindingID string,
	parameterInstance, parameterIn map[string]interface{}) error {
	dbInstanceID, err := c.getDBInstanceByID(instanceID)
	if err != nil {
		glog.Infof("UnBind faield for not found instance %s DBinstance. error:%v", instanceID, err)
		return err
	}

	accountName, _ := GetAccountInfoFromBindingParameters(parameterIn)
	glog.Infof("UnBind instance %s DBinstance %s's account name:%s.", instanceID, dbInstanceID, accountName)
	resetPassword := uniuri.New()
	err = c.ResetDatabaseAccount(instanceID, dbInstanceID, accountName, resetPassword)
	if err != nil {
		glog.Warningln(err)
		return err
	}
	return nil
}

func (c *RDSBroker) GetBindingStatus(instanceID, serviceID, planID, bindingID string,
	parameterInstance, parameterIn map[string]interface{}) (bool, error) {
	accountName, _ := GetAccountInfoFromBindingParameters(parameterIn)
	return c.CheckAccountStatus(instanceID, accountName, "Available")
}

func releaseName(id string) string {
	return "i-" + id
}

func createBrokerImpl() *RDSBroker {

	userMetaData, err := GetCloudServiceMetaData()
	if err != nil {
		glog.Infof("Failed to get cloud service metadate.")
		return nil
	}

	return &RDSBroker{
		client:    nil,
		region:    userMetaData.Region,
		zoneID:    userMetaData.ZoneID,
		vpcID:     userMetaData.VpcID,
		vSwitchID: userMetaData.VSwitchID,
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

func (c *RDSBroker) CreateDatabaseAccount(instanceID, dbInstanceID, accountName, accountPassword string) error {

	createAccountRequest := rds.CreateCreateAccountRequest()
	createAccountRequest.DBInstanceId = dbInstanceID
	createAccountRequest.AccountName = accountName
	createAccountRequest.AccountPassword = accountPassword
	createAccountRequest.AccountType = "Super"

	if err := c.CreateNewClientFromStsToken(); err != nil {
		glog.Infof("Create Rds client failed:\n%v\n", err)
		return err
	}

	glog.Infof("CreateDatabaseAccount DBinstance %s's account.", dbInstanceID)

	response, err := c.client.CreateAccount(createAccountRequest)
	if err != nil {
		glog.Infof("Failed to CreateDatabase:%v\n", err)
		return err
	}

	glog.Infof("CreateDatabaseAccount DBinstance %s's account response:%v.", dbInstanceID, response)
	return nil
}

func (c *RDSBroker) ResetDatabaseAccount(instanceID, dbInstanceID, accountName, password string) error {

	if err := c.CreateNewClientFromStsToken(); err != nil {
		glog.Infof("Create Rds client failed:\n%v\n", err)
		return err
	}

	resetPasswordRequest := rds.CreateResetAccountPasswordRequest()
	resetPasswordRequest.AccountName = accountName
	resetPasswordRequest.AccountPassword = password
	resetPasswordRequest.DBInstanceId = dbInstanceID

	response, err := c.client.ResetAccountPassword(resetPasswordRequest)

	if err != nil {
		glog.Infof("Failed to reset password %s :%v\n", password, err)
		return err
	}
	glog.Infof("Reset database acccount success \n%v\n", response)
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

func (c *RDSBroker) CreateNewClientFromStsToken() error {
	accessMetaData, err := util.GetAccessMetaData()
	if err != nil {
		glog.Infof("Failed to get access metadata err:%v\n", err)
		return err
	}

	client, err := rds.NewClientWithStsToken(c.region, accessMetaData.AccessKeyId,
		accessMetaData.AccessKeySecret, accessMetaData.SecurityToken)
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

func (c *RDSBroker) CreateVpcClientFromStsToken() (*vpc.Client, error) {
	accessMetaData, err := util.GetAccessMetaData()
	if err != nil {
		glog.Infof("Failed to get access metadate err:%v\n", err)
		return nil, err
	}

	client, err := vpc.NewClientWithStsToken(c.region, accessMetaData.AccessKeyId,
		accessMetaData.AccessKeySecret, accessMetaData.SecurityToken)
	if err != nil {
		glog.Infof("Failed to create rds client err:%v\n", err)
		return nil, err
	}
	return client, nil
}

func (c *RDSBroker) checkVSwitch(vSwitchIdLocal, zoneIDLocal string) (vSwitchId, zoneID string, err error) {
	err = c.CreateNewClientFromStsToken()
	if err != nil {
		return
	}
	req := rds.CreateDescribeRegionsRequest()
	req.RegionId = c.region
	glog.Infof("Describe regions request:%v:\n", req)
	response, err := c.client.DescribeRegions(req)
	if err != nil {
		glog.Infof("Describe Rds region failed:%v\n", err)
		return
	}
	glog.Infof("Describe regions response:%v:\n", response)
	regions := response.Regions.RDSRegion

	for _, region := range regions {
		if region.ZoneId == zoneIDLocal {
			vSwitchId = vSwitchIdLocal
			zoneID = zoneIDLocal
			return
		}
	}

	vpcClient, err := c.CreateVpcClientFromStsToken()
	if err != nil {
		glog.Infof("Create VPC client failed:%v\n", err)
		return
	}

	vpcReq := vpc.CreateDescribeVSwitchesRequest()
	vpcReq.RegionId = c.region
	vpcReq.VpcId = c.vpcID
	glog.Infof("Describe vswitch request:%v:\n", req)
	vpcResponse, err := vpcClient.DescribeVSwitches(vpcReq)
	if err != nil {
		glog.Infof("Describe vswitch failed:%v\n", err)
		return
	}
	glog.Infof("Describe vswitch response:%v:\n", vpcResponse)
	vSwitchs := vpcResponse.VSwitches.VSwitch

	for _, vSwitch := range vSwitchs {
		for _, region := range regions {
			if vSwitch.ZoneId == region.ZoneId {
				vSwitchId = vSwitch.VSwitchId
				zoneID = vSwitch.ZoneId
				err = nil
				glog.Infof("Found zone with availble Rds. zoneID:%s, vSwitchId:%s.\n", zoneID, vSwitchId)
				return
			}
		}
	}

	err = fmt.Errorf("Zone with availble Rds not found")
	glog.Infof("Zone with availble Rds not found.")
	return
}

func (c *RDSBroker) GetEcsVpcIps() (ips string, err error) {
	vpcClient, err := c.CreateVpcClientFromStsToken()
	if err != nil {
		glog.Infof("Create VPC client failed:%v\n", err)
		return
	}

	req := vpc.CreateDescribeVpcAttributeRequest()
	req.RegionId = c.region
	req.VpcId = c.vpcID
	glog.Infof("Describe vpc request:%v:\n", req)
	vpcResponse, err := vpcClient.DescribeVpcAttribute(req)
	if err != nil {
		glog.Infof("Describe vpc failed:%v\n", err)
		return
	}
	glog.Infof("Describe vpc response:%v:\n", vpcResponse)

	ips = vpcResponse.CidrBlock
	return
}

func (c *RDSBroker) SetRdsSecurityIps(ips, dbInstanceID string) error {
	err := c.CreateNewClientFromStsToken()
	if err != nil {
		return err
	}
	req := rds.CreateModifySecurityIpsRequest()
	req.RegionId = c.region
	req.SecurityIps = ips
	req.DBInstanceId = dbInstanceID
	glog.Infof("Set rds Security ips request:%v:\n", req)
	response, err := c.client.ModifySecurityIps(req)
	if err != nil {
		glog.Infof("Describe Rds region failed:%v\n", err)
		return err
	}
	glog.Infof("Set rds response:%v:\n", response)

	return nil
}

func (c *RDSBroker) SetRdsIpsForEcs(instanceID, dbInstanceID string) error {

	ips, err := c.GetEcsVpcIps()
	if err != nil {
		glog.Infof("Get vpc cidr failed:\n%v\n", err)
		err := fmt.Errorf("Get vpc cidr failed")
		return err
	}

	err = c.SetRdsSecurityIps(ips, dbInstanceID)
	if err != nil {
		glog.Infof("Set Rds security ips failed:\n%v\n", err)
		err := fmt.Errorf("Set Rds security ips failed.")
		return err
	}
	return nil
}
