package oss

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/golang/glog"
)

func CreateRamNewClient(region, accessKeyId, accessKeySecret string) (*ram.Client, error) {
	//	accessMetaData, err := GetAccessMetaData();
	//	if err != nil {
	//		glog.Infof("Failed to get access metadate err:%v\n", err)
	//		return nil, err
	//	}

	//	client, err := ram.NewClientWithStsToken(region, accessMetaData.AccessKeyId, accessMetaData.AccessKeySecret, accessMetaData.SecurityToken)
	client, err := ram.NewClientWithAccessKey(region, accessKeyId, accessKeySecret)
	if err != nil {
		glog.Infof("Failed to create ram client err:%v\n", err)
		return nil, err
	}

	return client, nil
}

func getUserNameForBinding(bindingId string) string {
	return "ram-user" + bindingId
}

func getGroupNameForInstance(instanceId string) string {
	return "ram-group" + instanceId
}

func getGroupPolicyNameForInstance(instanceId string) string {
	return "ram-policy" + instanceId
}

func createRamUser(bindingId, region, accessKeyId, accessKeySecret string) (*ram.User, error) {
	client, err := CreateRamNewClient(region, accessKeyId, accessKeySecret)
	if err != nil {
		glog.Infof("Create RAM client failed when create ram user.:\n%v\n", err)
		return nil, err
	}

	req := ram.CreateCreateUserRequest()
	req.UserName = getUserNameForBinding(bindingId)
	req.Scheme = "https"
	response, err := client.CreateUser(req)
	if err != nil {
		glog.Infof("Create RAM user for binding %s failed:%v\n", bindingId, err)
		return nil, err
	}
	glog.Infof("Create RAM user for binding %s success. response:%v\n", bindingId, response)
	return &response.User, nil
}

func createRamAccessKeyForUser(bindingId, region, accessKeyId, accessKeySecret string) (*ram.AccessKey, error) {
	client, err := CreateRamNewClient(region, accessKeyId, accessKeySecret)
	if err != nil {
		glog.Infof("Create RAM client failed when create ram access key for user.:\n%v\n", err)
		return nil, err
	}

	req := ram.CreateCreateAccessKeyRequest()
	req.UserName = getUserNameForBinding(bindingId)
	req.Scheme = "https"
	response, err := client.CreateAccessKey(req)
	if err != nil {
		glog.Infof("Create RAM user for binding %s failed:%v\n", bindingId, err)
		return nil, err
	}
	glog.Infof("Create RAM user access key for binding %s success. response:%v\n", bindingId, response)
	return &response.AccessKey, nil
}

func addUserToGroup(bindingId, instanceID, region, accessKeyId, accessKeySecret string) error {
	client, err := CreateRamNewClient(region, accessKeyId, accessKeySecret)
	if err != nil {
		glog.Infof("Create RAM client failed when add user to group.:\n%v\n", err)
		return err
	}

	req := ram.CreateAddUserToGroupRequest()
	req.Scheme = "https"
	req.UserName = getUserNameForBinding(bindingId)
	req.GroupName = getGroupNameForInstance(instanceID)
	response, err := client.AddUserToGroup(req)
	if err != nil {
		glog.Infof("Add RAM user to group for binding %s failed:%v\n", bindingId, err)
		return err
	}
	glog.Infof("add RAM user to group for binding %s success. response:%v\n", bindingId, response)
	return nil
}

func CreateRamUserForBinding(bindingId, instanceID, region, accessKeyId, accessKeySecret string) (*ram.User, *ram.AccessKey, error) {
	user, err := createRamUser(bindingId, region, accessKeyId, accessKeySecret)
	if err != nil {
		return nil, nil, err
	}

	accessKey, err := createRamAccessKeyForUser(bindingId, region, accessKeyId, accessKeySecret)
	if err != nil {
		deleteRamUser(bindingId, region, accessKeyId, accessKeySecret)
		return nil, nil, err
	}

	err = addUserToGroup(bindingId, instanceID, region, accessKeyId, accessKeySecret)
	if err != nil {
		deleteRamAccessKeyForUser(bindingId, region, accessKey.AccessKeyId, accessKeyId, accessKeySecret)
		deleteRamUser(bindingId, region, accessKeyId, accessKeySecret)
		return nil, nil, err
	}

	return user, accessKey, nil
}

func deleteRamUser(bindingId, region, accessKeyId, accessKeySecret string) error {
	client, err := CreateRamNewClient(region, accessKeyId, accessKeySecret)
	if err != nil {
		glog.Infof("Create RAM client failed when delete ram user.:\n%v\n", err)
		return err
	}

	req := ram.CreateDeleteUserRequest()
	req.Scheme = "https"
	req.UserName = getUserNameForBinding(bindingId)
	response, err := client.DeleteUser(req)
	if err != nil {
		glog.Infof("Delete RAM user for binding %s failed:%v\n", bindingId, err)
		return err
	}
	glog.Infof("Delete RAM user for binding %s success. response:%v\n", bindingId, response)
	return nil
}

func deleteRamAccessKeyForUser(bindingId, region, userAccessKeyId, accessKeyId, accessKeySecret string) error {
	client, err := CreateRamNewClient(region, accessKeyId, accessKeySecret)
	if err != nil {
		glog.Infof("Create RAM client failed when delete ram access key for user.:\n%v\n", err)
		return err
	}

	req := ram.CreateDeleteAccessKeyRequest()
	req.Scheme = "https"
	req.UserName = getUserNameForBinding(bindingId)
	req.UserAccessKeyId = userAccessKeyId
	response, err := client.DeleteAccessKey(req)
	if err != nil {
		glog.Infof("Delete RAM user for binding %s failed:%v\n", bindingId, err)
		return err
	}
	glog.Infof("Delete RAM user for binding %s success. response:%v\n", bindingId, response)
	return nil
}

func removeUserFromGroup(bindingId, instanceID, region, accessKeyId, accessKeySecret string) error {
	client, err := CreateRamNewClient(region, accessKeyId, accessKeySecret)
	if err != nil {
		glog.Infof("Create RAM client failed when remove user from group.:\n%v\n", err)
		return err
	}

	req := ram.CreateRemoveUserFromGroupRequest()
	req.Scheme = "https"
	req.UserName = getUserNameForBinding(bindingId)
	req.GroupName = getGroupNameForInstance(instanceID)
	response, err := client.RemoveUserFromGroup(req)
	if err != nil {
		glog.Infof("Remove RAM user from group for binding %s failed:%v\n", bindingId, err)
		return err
	}
	glog.Infof("Remove RAM user from group for binding %s success. response:%v\n", bindingId, response)
	return nil
}

func DeleteRamUserForBinding(bindingId, instanceID, region, userAccessKeyId, accessKeyId, accessKeySecret string) error {

	err := removeUserFromGroup(bindingId, instanceID, region, accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}

	err = deleteRamAccessKeyForUser(bindingId, region, userAccessKeyId, accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}

	err = deleteRamUser(bindingId, region, accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}

	return nil
}

func getRamUser(bindingId, region, accessKeyId, accessKeySecret string) (*ram.User, error) {
	client, err := CreateRamNewClient(region, accessKeyId, accessKeySecret)
	if err != nil {
		glog.Infof("Create RAM client failed when get ram user.:\n%v\n", err)
		return nil, err
	}

	req := ram.CreateGetUserRequest()
	req.UserName = getUserNameForBinding(bindingId)
	req.Scheme = "https"
	response, err := client.GetUser(req)
	if err != nil {
		glog.Infof("Get RAM user info for binding %s failed:%v\n", bindingId, err)
		return nil, err
	}
	glog.Infof("Get RAM user for binding %s success. response:%v\n", bindingId, response)
	return &response.User, nil
}

func GetRamUserForBinding(bindingId, instanceID, region, accessKeyId, accessKeySecret string) error {

	_, err := getRamUser(bindingId, region, accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}
	return nil
}

func createRamGroup(instanceId, region, accessKeyId, accessKeySecret string) error {
	client, err := CreateRamNewClient(region, accessKeyId, accessKeySecret)
	if err != nil {
		glog.Infof("Create RAM client failed:\n%v\n", err)
		return err
	}

	req := ram.CreateCreateGroupRequest()
	req.Scheme = "https"
	req.GroupName = getGroupNameForInstance(instanceId)
	response, err := client.CreateGroup(req)
	if err != nil {
		glog.Infof("Create RAM Group for instance %s failed:%v\n", instanceId, err)
		return err
	}
	glog.Infof("Create RAM Group for instance %s success. response:%v\n", instanceId, response)
	return nil
}

func createRamGroupPolicy(instanceId, region, accessKeyId, accessKeySecret string) error {
	client, err := CreateRamNewClient(region, accessKeyId, accessKeySecret)
	if err != nil {
		glog.Infof("Create RAM client failed:\n%v\n", err)
		return err
	}

	req := ram.CreateCreatePolicyRequest()
	req.Scheme = "https"
	req.PolicyName = getGroupPolicyNameForInstance(instanceId)
	req.PolicyDocument = GetOssRamPolicy(instanceId)
	response, err := client.CreatePolicy(req)
	if err != nil {
		glog.Infof("Create RAM Group Policy for instance %s failed:%v\n", instanceId, err)
		return err
	}
	glog.Infof("Create RAM Group Policy for instance %s success. response:%v\n", instanceId, response)
	return nil
}

func attachRamGroupPolicy(instanceId, region, accessKeyId, accessKeySecret string) error {
	client, err := CreateRamNewClient(region, accessKeyId, accessKeySecret)
	if err != nil {
		glog.Infof("Create RAM client failed:\n%v\n", err)
		return err
	}

	req := ram.CreateAttachPolicyToGroupRequest()
	req.Scheme = "https"
	req.PolicyName = getGroupPolicyNameForInstance(instanceId)
	req.GroupName = getGroupNameForInstance(instanceId)
	req.PolicyType = "Custom"
	response, err := client.AttachPolicyToGroup(req)
	if err != nil {
		glog.Infof("Attach RAM Group Policy for instance %s failed:%v\n", instanceId, err)
		return err
	}
	glog.Infof("Attach RAM Group Policy for instance %s success. response:%v\n", instanceId, response)
	return nil
}

func CreateRamForInstance(instanceId, region, accessKeyId, accessKeySecret string) error {
	err := createRamGroup(instanceId, region, accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}

	err = createRamGroupPolicy(instanceId, region, accessKeyId, accessKeySecret)
	if err != nil {
		deleteRamGroup(instanceId, region, accessKeyId, accessKeySecret)
		return err
	}

	err = attachRamGroupPolicy(instanceId, region, accessKeyId, accessKeySecret)
	if err != nil {
		deleteRamGroupPolicy(instanceId, region, accessKeyId, accessKeySecret)
		deleteRamGroup(instanceId, region, accessKeyId, accessKeySecret)
		return err
	}
	return nil
}

func deleteRamGroup(instanceId, region, accessKeyId, accessKeySecret string) error {
	client, err := CreateRamNewClient(region, accessKeyId, accessKeySecret)
	if err != nil {
		glog.Infof("Create RAM client failed when delete ram group:%v\n", err)
		return err
	}

	req := ram.CreateDeleteGroupRequest()
	req.Scheme = "https"
	req.GroupName = getGroupNameForInstance(instanceId)
	response, err := client.DeleteGroup(req)
	if err != nil {
		glog.Infof("Delete RAM Group for instance %s failed:%v\n", instanceId, err)
		return err
	}
	glog.Infof("Delete RAM Group for instance %s success. response:%v\n", instanceId, response)
	return nil
}

func deleteRamGroupPolicy(instanceId, region, accessKeyId, accessKeySecret string) error {
	client, err := CreateRamNewClient(region, accessKeyId, accessKeySecret)
	if err != nil {
		glog.Infof("Create RAM client failed when delete ram group policy:\n%v\n", err)
		return err
	}

	req := ram.CreateDeletePolicyRequest()
	req.Scheme = "https"
	req.PolicyName = getGroupPolicyNameForInstance(instanceId)
	response, err := client.DeletePolicy(req)
	if err != nil {
		glog.Infof("Delete RAM Group Policy for instance %s failed:%v\n", instanceId, err)
		return err
	}
	glog.Infof("Delete RAM Group Policy for instance %s success. response:%v\n", instanceId, response)
	return nil
}

func detachRamGroupPolicy(instanceId, region, accessKeyId, accessKeySecret string) error {
	client, err := CreateRamNewClient(region, accessKeyId, accessKeySecret)
	if err != nil {
		glog.Infof("Create RAM client failed when detach ram group policy:\n%v\n", err)
		return err
	}

	req := ram.CreateDetachPolicyFromGroupRequest()
	req.Scheme = "https"
	req.PolicyName = getGroupPolicyNameForInstance(instanceId)
	req.GroupName = getGroupNameForInstance(instanceId)
	req.PolicyType = "Custom"
	response, err := client.DetachPolicyFromGroup(req)
	if err != nil {
		glog.Infof("Detach RAM Group Policy for instance %s failed:%v\n", instanceId, err)
		return err
	}
	glog.Infof("Detach RAM Group Policy for instance %s success. response:%v\n", instanceId, response)
	return nil
}

func DeleteRamForInstance(instanceId, region, accessKeyId, accessKeySecret string) error {
	err := detachRamGroupPolicy(instanceId, region, accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}

	err = deleteRamGroupPolicy(instanceId, region, accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}

	err = deleteRamGroup(instanceId, region, accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}
	return nil
}

func checkRamGroupPolicy(instanceId, region, accessKeyId, accessKeySecret string) bool {
	client, err := CreateRamNewClient(region, accessKeyId, accessKeySecret)
	if err != nil {
		glog.Infof("Create RAM client failed when check ram group policy:%v\n", err)
		return false
	}

	req := ram.CreateGetPolicyRequest()
	req.Scheme = "https"
	req.PolicyName = getGroupPolicyNameForInstance(instanceId)
	req.PolicyType = "Custom"
	response, err := client.GetPolicy(req)
	if err != nil {
		glog.Infof("List RAM Group policy for instance %s failed:%v\n", instanceId, err)
		return false
	}

	if response.Policy.AttachmentCount == 1 {
		glog.Infof("Found RAM Group for instance %s.", instanceId)
		return true
	}

	glog.Infof("Not found RAM Group for instance %s success. response:%v\n", instanceId, response)
	return false
}

func CheckRamForInstance(instanceId, region, accessKeyId, accessKeySecret string) bool {
	return checkRamGroupPolicy(instanceId, region, accessKeyId, accessKeySecret)
}
