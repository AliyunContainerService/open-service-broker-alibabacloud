package polardb

import (
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/golang/glog"
)

type ConnectionInfo struct {
	ConnectionString string
	Port             string
}

func (c *PolarDBBroker) getDBClusterConnectionInfo(dbInstanceID string) (*ConnectionInfo, error) {
	req := polardb.CreateDescribeDBClusterEndpointsRequest()
	req.DBClusterId = dbInstanceID
	resp, err := c.client.DescribeDBClusterEndpoints(req)
	err = handleErrorResponse(resp.BaseResponse, err)
	if err != nil {
		glog.Errorf("Failed to DescribeDBClusterEndpoints: %v", err)
		return nil, err
	}
	info := &ConnectionInfo{}
	for _, dbEndpoint := range resp.Items {
		if dbEndpoint.EndpointType == "Cluster" {
			info.ConnectionString = dbEndpoint.AddressItems[0].ConnectionString
			info.Port = dbEndpoint.AddressItems[0].Port
		}
	}
	if info.ConnectionString == "" {
		return nil, errors.New("get connectionInfo fail")
	}

	return info, nil
}

func (c *PolarDBBroker) getDBClusterByID(instanceID string) (string, error) {

	req := polardb.CreateDescribeDBClustersRequest()
	req.DBClusterDescription = instanceID
	req.Tag = &[]polardb.DescribeDBClustersTag{{Value: "true", Key: SERVICE_CATALOG_TAG_KEY}}
	req.PageSize = requests.NewInteger(DEFAULT_PAGE_SIZE)
	req.PageNumber = "1"

	for {
		resp, err := c.client.DescribeDBClusters(req)

		err = handleErrorResponse(resp.BaseResponse, err)
		if err != nil {
			glog.Errorf("Failed to DescribeDBClusters: %v", err)
			return "", err
		} else {
			for _, dbInstance := range resp.Items.DBCluster {
				if dbInstance.DBClusterDescription == instanceID {
					return dbInstance.DBClusterId, nil
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

func (c *PolarDBBroker) CheckDBClusterStatus(instanceID, status string, parameterIn map[string]interface{}) (bool, error) {

	req := polardb.CreateDescribeDBClustersRequest()
	req.DBClusterDescription = instanceID
	req.Tag = &[]polardb.DescribeDBClustersTag{{Value: "true", Key: SERVICE_CATALOG_TAG_KEY}}
	req.PageSize = requests.NewInteger(DEFAULT_PAGE_SIZE)
	req.PageNumber = "1"

	glog.Infof("checkDBClusterStatus in: instanceId:\n%v\n", instanceID)

	for {
		resp, err := c.client.DescribeDBClusters(req)

		err = handleErrorResponse(resp.BaseResponse, err)
		if err != nil {
			glog.Errorf("Failed to DescribeDBClusters: %v", err)
			return false, err
		} else {
			for _, dbInstance := range resp.Items.DBCluster {
				glog.Infof("checkDBClusterStatus of dbInstance:\n%v\n", dbInstance)
				if dbInstance.DBClusterDescription == instanceID {
					if dbInstance.DBClusterStatus == status {
						glog.Infof("checkDBClusterStatus found instance status is running.")
						return true, nil
					}
					glog.Infof("checkDBClusterStatus found instance status is %s ", dbInstance.DBClusterStatus)
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

func (c *PolarDBBroker) CheckAccountStatus(instanceID string, accountName string, status string) (bool, error) {

	dBInstanceID, err := c.getDBClusterByID(instanceID)
	if err != nil {
		glog.Infof("checkAccountStatus: Get dbInstanceID failed:\n%v\n", err)
		return false, err
	}

	args := polardb.CreateDescribeAccountsRequest()
	args.DBClusterId = dBInstanceID
	args.AccountName = accountName

	resp, err := c.client.DescribeAccounts(args)

	err = handleErrorResponse(resp.BaseResponse, err)
	if err != nil {
		return false, err
	}

	if len(resp.Accounts) < 1 {
		return false, err
	}

	for _, account := range resp.Accounts {
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
