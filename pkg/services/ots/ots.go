package ots

import (
	"fmt"
	"github.com/AliyunContainerService/open-service-broker-alibabacloud/pkg/util"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/denverdino/aliyungo/metadata"
	"github.com/golang/glog"
)

type ConnectionInfo struct {
	ConnectionString string
	Port             string
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
		glog.Errorf("Get Zone error, %v", err.Error())
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

func createOtsNewClient(instanceID, region string) (*ots.Client, error, string) {
	accessMetaData, err := util.GetAccessMetaData()
	if err != nil {
		glog.Infof("Failed to get access metadate err:%v\n", err)
		return nil, err, ""
	}
	client, err := ots.NewClientWithStsToken(region, accessMetaData.AccessKeyId,
		accessMetaData.AccessKeySecret, accessMetaData.SecurityToken)
	if err != nil {
		err := fmt.Errorf("Failed to create ots client err.")
		glog.Infof("Failed to create ots client err.")
		return nil, err, ""
	}

	return client, nil, accessMetaData.AccessKeyId
}

func GetOtsNameForInstance(instanceId string) string {
	strLen := len(instanceId)
	return "b-" + instanceId[strLen-12:strLen]
}

func CreateOtsInstance(instanceId, region string) (string, error) {
	client, err, accessKeyId := createOtsNewClient(instanceId, region)
	if err != nil {
		glog.Infof("Create OTS client failed:\n%v\n", err)
		return "", err
	}

	req := ots.CreateInsertInstanceRequest()
	req.InstanceName = GetOtsNameForInstance(instanceId)
	req.Domain = GetOTSEndPoint(region)
	req.AccessKeyId = accessKeyId

	rsp, err := client.InsertInstance(req)
	if err != nil {
		glog.Infof("Create OTS instance failed:\n%v\n", err)
		return "", err
	}
	glog.Infof("Create OTS instance success reponse:%v.", rsp)
	instanceUrl := GetOTSInstanceUrl(instanceId, region)
	return instanceUrl, nil
}

func DeleteOtsInstance(instanceId, region string) error {
	client, err, accessKeyId := createOtsNewClient(instanceId, region)
	if err != nil {
		glog.Infof("Create OTS client failed when delete bucket:\n%v\n", err)
		return err
	}

	req := ots.CreateDeleteInstanceRequest()
	req.InstanceName = GetOtsNameForInstance(instanceId)
	req.Domain = GetOTSEndPoint(region)
	req.AccessKeyId = accessKeyId

	rsp, err := client.DeleteInstance(req)
	if err != nil {
		glog.Infof("Delete OTS instance for instance %s failed:\n%v\n", instanceId, err)
		return err
	}
	glog.Infof("Delete OTS instance success reponse:%v.", rsp)
	return nil
}

func CheckOtsInstanceStatus(instanceID, region string) (bool, error) {

	client, err, accessKeyId := createOtsNewClient(instanceID, region)
	if err != nil {
		glog.Infof("Create OTS client failed:\n%v\n", err)
		return false, err
	}

	req := ots.CreateGetInstanceRequest()
	req.InstanceName = GetOtsNameForInstance(instanceID)
	req.Domain = GetOTSEndPoint(region)
	req.AccessKeyId = accessKeyId
	req.Method = "GET"

	rsp, err := client.GetInstance(req)
	if err != nil {
		glog.Infof("Get OTS instance failed error: %v\n", err)
		return false, err
	}
	glog.Infof("Get OTS instance success.reponse: %v\n", rsp)

	return true, nil
}

func GetOTSInstanceUrl(instanceID, region string) string {
	return GetOtsNameForInstance(instanceID) + "." + region + ".ots.aliyuncs.com"
}

func GetOTSEndPoint(region string) string {
	return "ots." + region + ".aliyuncs.com"
}

var policy = `{
  "Version": "1",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ots:*"
      ],
      "Resource": [
        "*"
      ],
      "Condition": {}
    }
  ]
}`

func GetOtsRamPolicy(instanceID string) string {
	return policy
}
