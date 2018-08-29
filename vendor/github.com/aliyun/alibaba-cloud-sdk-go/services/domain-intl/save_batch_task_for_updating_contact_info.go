package domain_intl

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

// SaveBatchTaskForUpdatingContactInfo invokes the domain_intl.SaveBatchTaskForUpdatingContactInfo API synchronously
// api document: https://help.aliyun.com/api/domain-intl/savebatchtaskforupdatingcontactinfo.html
func (client *Client) SaveBatchTaskForUpdatingContactInfo(request *SaveBatchTaskForUpdatingContactInfoRequest) (response *SaveBatchTaskForUpdatingContactInfoResponse, err error) {
	response = CreateSaveBatchTaskForUpdatingContactInfoResponse()
	err = client.DoAction(request, response)
	return
}

// SaveBatchTaskForUpdatingContactInfoWithChan invokes the domain_intl.SaveBatchTaskForUpdatingContactInfo API asynchronously
// api document: https://help.aliyun.com/api/domain-intl/savebatchtaskforupdatingcontactinfo.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SaveBatchTaskForUpdatingContactInfoWithChan(request *SaveBatchTaskForUpdatingContactInfoRequest) (<-chan *SaveBatchTaskForUpdatingContactInfoResponse, <-chan error) {
	responseChan := make(chan *SaveBatchTaskForUpdatingContactInfoResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.SaveBatchTaskForUpdatingContactInfo(request)
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

// SaveBatchTaskForUpdatingContactInfoWithCallback invokes the domain_intl.SaveBatchTaskForUpdatingContactInfo API asynchronously
// api document: https://help.aliyun.com/api/domain-intl/savebatchtaskforupdatingcontactinfo.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SaveBatchTaskForUpdatingContactInfoWithCallback(request *SaveBatchTaskForUpdatingContactInfoRequest, callback func(response *SaveBatchTaskForUpdatingContactInfoResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *SaveBatchTaskForUpdatingContactInfoResponse
		var err error
		defer close(result)
		response, err = client.SaveBatchTaskForUpdatingContactInfo(request)
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

// SaveBatchTaskForUpdatingContactInfoRequest is the request struct for api SaveBatchTaskForUpdatingContactInfo
type SaveBatchTaskForUpdatingContactInfoRequest struct {
	*requests.RpcRequest
	UserClientIp        string           `position:"Query" name:"UserClientIp"`
	Lang                string           `position:"Query" name:"Lang"`
	RegistrantProfileId requests.Integer `position:"Query" name:"RegistrantProfileId"`
	ContactType         string           `position:"Query" name:"ContactType"`
	AddTransferLock     requests.Boolean `position:"Query" name:"AddTransferLock"`
	DomainName          *[]string        `position:"Query" name:"DomainName"  type:"Repeated"`
}

// SaveBatchTaskForUpdatingContactInfoResponse is the response struct for api SaveBatchTaskForUpdatingContactInfo
type SaveBatchTaskForUpdatingContactInfoResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	TaskNo    string `json:"TaskNo" xml:"TaskNo"`
}

// CreateSaveBatchTaskForUpdatingContactInfoRequest creates a request to invoke SaveBatchTaskForUpdatingContactInfo API
func CreateSaveBatchTaskForUpdatingContactInfoRequest() (request *SaveBatchTaskForUpdatingContactInfoRequest) {
	request = &SaveBatchTaskForUpdatingContactInfoRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Domain-intl", "2017-12-18", "SaveBatchTaskForUpdatingContactInfo", "domain", "openAPI")
	return
}

// CreateSaveBatchTaskForUpdatingContactInfoResponse creates a response to parse from SaveBatchTaskForUpdatingContactInfo response
func CreateSaveBatchTaskForUpdatingContactInfoResponse() (response *SaveBatchTaskForUpdatingContactInfoResponse) {
	response = &SaveBatchTaskForUpdatingContactInfoResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
