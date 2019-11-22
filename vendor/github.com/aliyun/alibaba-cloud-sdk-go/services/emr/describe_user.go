package emr

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// DescribeUser invokes the emr.DescribeUser API synchronously
// api document: https://help.aliyun.com/api/emr/describeuser.html
func (client *Client) DescribeUser(request *DescribeUserRequest) (response *DescribeUserResponse, err error) {
	response = CreateDescribeUserResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeUserWithChan invokes the emr.DescribeUser API asynchronously
// api document: https://help.aliyun.com/api/emr/describeuser.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeUserWithChan(request *DescribeUserRequest) (<-chan *DescribeUserResponse, <-chan error) {
	responseChan := make(chan *DescribeUserResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeUser(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// DescribeUserWithCallback invokes the emr.DescribeUser API asynchronously
// api document: https://help.aliyun.com/api/emr/describeuser.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeUserWithCallback(request *DescribeUserRequest, callback func(response *DescribeUserResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeUserResponse
		var err error
		defer close(result)
		response, err = client.DescribeUser(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// DescribeUserRequest is the request struct for api DescribeUser
type DescribeUserRequest struct {
	*requests.RpcRequest
	ResourceOwnerId requests.Integer `position:"Query" name:"ResourceOwnerId"`
	AliyunUserId    string           `position:"Query" name:"AliyunUserId"`
}

// DescribeUserResponse is the response struct for api DescribeUser
type DescribeUserResponse struct {
	*responses.BaseResponse
	Paging         bool                         `json:"Paging" xml:"Paging"`
	RequestId      string                       `json:"RequestId" xml:"RequestId"`
	AliyunUserId   string                       `json:"AliyunUserId" xml:"AliyunUserId"`
	UserName       string                       `json:"UserName" xml:"UserName"`
	UserType       string                       `json:"UserType" xml:"UserType"`
	Status         string                       `json:"Status" xml:"Status"`
	IsSuperAdmin   string                       `json:"IsSuperAdmin" xml:"IsSuperAdmin"`
	Description    string                       `json:"Description" xml:"Description"`
	GmtCreate      string                       `json:"GmtCreate" xml:"GmtCreate"`
	GmtModified    string                       `json:"GmtModified" xml:"GmtModified"`
	RoleDTOList    RoleDTOListInDescribeUser    `json:"RoleDTOList" xml:"RoleDTOList"`
	GroupDTOList   GroupDTOListInDescribeUser   `json:"GroupDTOList" xml:"GroupDTOList"`
	AccountDTOList AccountDTOListInDescribeUser `json:"AccountDTOList" xml:"AccountDTOList"`
}

// CreateDescribeUserRequest creates a request to invoke DescribeUser API
func CreateDescribeUserRequest() (request *DescribeUserRequest) {
	request = &DescribeUserRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Emr", "2016-04-08", "DescribeUser", "emr", "openAPI")
	return
}

// CreateDescribeUserResponse creates a response to parse from DescribeUser response
func CreateDescribeUserResponse() (response *DescribeUserResponse) {
	response = &DescribeUserResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
