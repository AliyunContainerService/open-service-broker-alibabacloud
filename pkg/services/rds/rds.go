package rds

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/denverdino/aliyungo/metadata"
	"github.com/golang/glog"
)

type ConnectionInfo struct {
	ConnectionString string
	Port             string
}

func (c *RDSBroker) getDBInstanceConnectionInfo(dbInstanceID string) (*ConnectionInfo, error) {
	if err := c.CreateNewClientFromStsToken(); err != nil {
		glog.Infof("getDBInstanceConnectionInfo: Create Rds client failed:\n%v\n", err)
		return nil, err
	}
	req := rds.CreateDescribeDBInstanceAttributeRequest()
	req.DBInstanceId = dbInstanceID
	resp, err := c.client.DescribeDBInstanceAttribute(req)
	err = handleErrorResponse(resp.BaseResponse, err)
	if err != nil {
		glog.Errorf("Failed to DescribeDBInstanceAttribute: %v", err)
		return nil, err
	}
	info := ConnectionInfo{}
	attr := resp.Items.DBInstanceAttribute[0]
	info.ConnectionString = attr.ConnectionString
	info.Port = attr.Port
	return &info, nil
}

func (c *RDSBroker) getDBInstanceByID(serviceID string) (string, error) {

	req := rds.CreateDescribeDBInstancesRequest()
	req.SearchKey = serviceID
	req.Tag1Key = SERVICE_CATALOG_TAG_KEY
	req.Tag1Value = "true"
	req.PageSize = requests.NewInteger(DEFAULT_PAGE_SIZE)
	req.PageNumber = "1"

	if err := c.CreateNewClientFromStsToken(); err != nil {
		glog.Infof("getDBInstanceByID: Create Rds client failed:\n%v\n", err)
		return "", err
	}

	for {
		resp, err := c.client.DescribeDBInstances(req)

		err = handleErrorResponse(resp.BaseResponse, err)
		if err != nil {
			glog.Errorf("Failed to DescribeDBInstances: %v", err)
			return "", err
		} else {
			for _, dbInstance := range resp.Items.DBInstance {
				if dbInstance.DBInstanceDescription == serviceID {
					return dbInstance.DBInstanceId, nil
				}
			}
			nextPage := getNextPageNumber(resp.PageNumber, DEFAULT_PAGE_SIZE, resp.TotalRecordCount)
			if nextPage == 0 {
				break
			} else {
				req.PageNumber = requests.NewInteger(nextPage)
			}
		}
	}

	return "", nil
}

func (c *RDSBroker) CheckDBInstanceStatus(instanceID, status string) (bool, error) {

	req := rds.CreateDescribeDBInstancesRequest()
	req.SearchKey = instanceID
	req.Tag1Key = SERVICE_CATALOG_TAG_KEY
	req.Tag1Value = "true"
	req.PageSize = requests.NewInteger(DEFAULT_PAGE_SIZE)
	req.PageNumber = "1"

	glog.Infof("checkDBInstanceStatus in: instanceId:\n%v\n", instanceID)
	if err := c.CreateNewClientFromStsToken(); err != nil {
		glog.Infof("checkDBInstanceStatus: Create Rds client failed:\n%v\n", err)
		return false, err
	}

	for {
		resp, err := c.client.DescribeDBInstances(req)

		err = handleErrorResponse(resp.BaseResponse, err)
		if err != nil {
			glog.Errorf("Failed to DescribeDBInstances: %v", err)
			return false, err
		} else {
			for _, dbInstance := range resp.Items.DBInstance {
				glog.Infof("checkDBInstanceStatus of dbInstance:\n%v\n", dbInstance)
				if dbInstance.DBInstanceDescription == instanceID {
					if dbInstance.DBInstanceStatus == status {
						err = c.SetRdsIpsForEcs(instanceID, dbInstance.DBInstanceId)
						if err != nil {
							glog.Infof("Set rds security ips for ecs failed:%v\n")
							return false, err
						}
						glog.Infof("checkDBInstanceStatus found instance status is running.")
						return true, nil
					}
					glog.Infof("checkDBInstanceStatus found instance but status %s is not running\n", dbInstance.DBInstanceStatus)
					return false, nil
				}
			}
			nextPage := getNextPageNumber(resp.PageNumber, DEFAULT_PAGE_SIZE, resp.TotalRecordCount)
			if nextPage == 0 {
				break
			} else {
				req.PageNumber = requests.NewInteger(nextPage)
			}
		}
	}

	return false, nil
}

func (c *RDSBroker) CheckAccountStatus(instanceID string, accountName string, status string) (bool, error) {

	dBInstanceID, err := c.getDBInstanceByID(instanceID)
	if err != nil {
		glog.Infof("checkAccountStatus: Get dbInstanceID failed:\n%v\n", err)
		return false, err
	}

	args := rds.CreateDescribeAccountsRequest()
	args.DBInstanceId = dBInstanceID
	args.AccountName = accountName

	resp, err := c.client.DescribeAccounts(args)

	err = handleErrorResponse(resp.BaseResponse, err)
	if err != nil {
		return false, err
	}

	accounts := resp.Accounts.DBInstanceAccount

	if len(accounts) < 1 {
		return false, err
	}

	for _, account := range accounts {
		glog.Infof("checkAccountStatus: account:\n%v\n", account)
		if account.AccountName == accountName && account.AccountStatus == status {
			return true, nil
		}
	}

	glog.Infof("checkAccountStatus: account %s not found.", accountName)
	return false, err
}

type UserMetaData struct {
	AccessKeyId     string
	AccessKeySecret string
	SecurityToken   string
	Region          string
	ZoneID          string
	VpcID           string
	VSwitchID       string
}

func GetCloudServiceMetaData() (*UserMetaData, error) {
	m := metadata.NewMetaData(nil)

	region, err := m.Region()
	if err != nil {
		glog.Errorf("Get Region error, %v", err.Error())
		return nil, err
	}
	vpcID, err := m.VpcID()
	if err != nil {
		glog.Errorf("Get VpcID error, %v", err.Error())
		return nil, err
	}
	vSwitchID, err := m.VswitchID()
	if err != nil {
		glog.Errorf("Get VswitchID error, %v", err.Error())
		return nil, err
	}
	zoneID, err := m.Zone()
	if err != nil {
		glog.Errorf("Get VswitchID error, %v", err.Error())
		return nil, err
	}
	userMetaData := UserMetaData{
		Region:    region,
		ZoneID:    zoneID,
		VpcID:     vpcID,
		VSwitchID: vSwitchID,
	}
	return &userMetaData, nil
}
