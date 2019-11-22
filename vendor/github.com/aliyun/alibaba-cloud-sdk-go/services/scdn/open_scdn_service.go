package scdn

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

// OpenScdnService invokes the scdn.OpenScdnService API synchronously
// api document: https://help.aliyun.com/api/scdn/openscdnservice.html
func (client *Client) OpenScdnService(request *OpenScdnServiceRequest) (response *OpenScdnServiceResponse, err error) {
	response = CreateOpenScdnServiceResponse()
	err = client.DoAction(request, response)
	return
}

// OpenScdnServiceWithChan invokes the scdn.OpenScdnService API asynchronously
// api document: https://help.aliyun.com/api/scdn/openscdnservice.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) OpenScdnServiceWithChan(request *OpenScdnServiceRequest) (<-chan *OpenScdnServiceResponse, <-chan error) {
	responseChan := make(chan *OpenScdnServiceResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.OpenScdnService(request)
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

// OpenScdnServiceWithCallback invokes the scdn.OpenScdnService API asynchronously
// api document: https://help.aliyun.com/api/scdn/openscdnservice.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) OpenScdnServiceWithCallback(request *OpenScdnServiceRequest, callback func(response *OpenScdnServiceResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *OpenScdnServiceResponse
		var err error
		defer close(result)
		response, err = client.OpenScdnService(request)
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

// OpenScdnServiceRequest is the request struct for api OpenScdnService
type OpenScdnServiceRequest struct {
	*requests.RpcRequest
	StartDate         string           `position:"Query" name:"StartDate"`
	CcProtection      requests.Integer `position:"Query" name:"CcProtection"`
	SecurityToken     string           `position:"Query" name:"SecurityToken"`
	ProtectType       string           `position:"Query" name:"ProtectType"`
	DDoSBasic         requests.Integer `position:"Query" name:"DDoSBasic"`
	Bandwidth         requests.Integer `position:"Query" name:"Bandwidth"`
	DomainCount       requests.Integer `position:"Query" name:"DomainCount"`
	OwnerId           requests.Integer `position:"Query" name:"OwnerId"`
	EndDate           string           `position:"Query" name:"EndDate"`
	ElasticProtection requests.Integer `position:"Query" name:"ElasticProtection"`
}

// OpenScdnServiceResponse is the response struct for api OpenScdnService
type OpenScdnServiceResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateOpenScdnServiceRequest creates a request to invoke OpenScdnService API
func CreateOpenScdnServiceRequest() (request *OpenScdnServiceRequest) {
	request = &OpenScdnServiceRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("scdn", "2017-11-15", "OpenScdnService", "", "")
	return
}

// CreateOpenScdnServiceResponse creates a response to parse from OpenScdnService response
func CreateOpenScdnServiceResponse() (response *OpenScdnServiceResponse) {
	response = &OpenScdnServiceResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
