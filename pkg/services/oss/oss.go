package oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/denverdino/aliyungo/metadata"
	"github.com/golang/glog"
	"strings"
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

func createOssNewClient(region, accessKeyId, accessKeySecret string) (*oss.Client, error) {
	client, err := oss.New(GetOSSEndPoint(region), accessKeyId, accessKeySecret)
	if err != nil {
		glog.Infof("Failed to create oss client err:%v\n", err)
		return nil, err
	}

	return client, nil
}

func GetBucketNameForInstance(instanceId string) string {
	return "bucket-" + instanceId
}

func CreateOssBucket(instanceId, region, accessKeyId, accessKeySecret, acl, storageClass string) (string, error) {
	client, err := createOssNewClient(region, accessKeyId, accessKeySecret)
	if err != nil {
		glog.Infof("Create OSS client failed:\n%v\n", err)
		return "", err
	}

	bucketName := GetBucketNameForInstance(instanceId)

	if acl == "" {
		acl = "private"
	}

	if storageClass == "" {
		storageClass = "Standard"
	}

	err = client.CreateBucket(bucketName, oss.StorageClass(oss.StorageClassType(storageClass)), oss.ACL(oss.ACLType(acl)))
	if err != nil {
		glog.Infof("Create OSS Bucket failed:\n%v\n", err)
		return "", err
	}

	return bucketName, nil
}

func DeleteOssBucket(instanceId, region, accessKeyId, accessKeySecret string) error {
	client, err := createOssNewClient(region, accessKeyId, accessKeySecret)
	if err != nil {
		glog.Infof("Create OSS client failed when delete bucket:\n%v\n", err)
		return err
	}

	bucketName := GetBucketNameForInstance(instanceId)
	err = client.DeleteBucket(bucketName)
	if err != nil {
		glog.Infof("delete OSS bucket for instance %s failed:\n%v\n", instanceId, err)
		return err
	}
	return nil
}

func CheckOssInstanceStatus(instanceID, region, accessKeyId, accessKeySecret string) (bool, error) {

	client, err := createOssNewClient(region, accessKeyId, accessKeySecret)
	if err != nil {
		glog.Infof("Create OSS client failed:\n%v\n", err)
		return false, err
	}

	lsRes, err := client.ListBuckets()
	if err != nil {
		glog.Infof("List OSS bucket failed error: %v\n", err)
		return false, err
	}

	for _, bucket := range lsRes.Buckets {
		if GetBucketNameForInstance(instanceID) == bucket.Name {
			glog.Infof("Found OSS bucket instance %s success.", instanceID)
			return true, nil
		}
	}

	glog.Infof("Not found OSS bucket instance %s failed error: %v\n", instanceID, err)

	return false, nil
}

func GetOSSEndPoint(region string) string {
	return "oss-" + region + ".aliyuncs.com"
}

var policy = `{
  "Version": "1",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "oss:*"
      ],
      "Resource": [
        "acs:oss:*:*:${bucketName}",
        "acs:oss:*:*:${bucketName}/*"
      ],
      "Condition": {}
    }
  ]
}`

func GetOssRamPolicy(instanceID string) string {
	return strings.Replace(policy, "${bucketName}", GetBucketNameForInstance(instanceID), -1)
}
