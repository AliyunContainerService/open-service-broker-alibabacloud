package ecs

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

// SignAgreement invokes the ecs.SignAgreement API synchronously
// api document: https://help.aliyun.com/api/ecs/signagreement.html
func (client *Client) SignAgreement(request *SignAgreementRequest) (response *SignAgreementResponse, err error) {
	response = CreateSignAgreementResponse()
	err = client.DoAction(request, response)
	return
}

// SignAgreementWithChan invokes the ecs.SignAgreement API asynchronously
// api document: https://help.aliyun.com/api/ecs/signagreement.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SignAgreementWithChan(request *SignAgreementRequest) (<-chan *SignAgreementResponse, <-chan error) {
	responseChan := make(chan *SignAgreementResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.SignAgreement(request)
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

// SignAgreementWithCallback invokes the ecs.SignAgreement API asynchronously
// api document: https://help.aliyun.com/api/ecs/signagreement.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SignAgreementWithCallback(request *SignAgreementRequest, callback func(response *SignAgreementResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *SignAgreementResponse
		var err error
		defer close(result)
		response, err = client.SignAgreement(request)
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

// SignAgreementRequest is the request struct for api SignAgreement
type SignAgreementRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	AgreementType        string           `position:"Query" name:"AgreementType"`
}

// SignAgreementResponse is the response struct for api SignAgreement
type SignAgreementResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateSignAgreementRequest creates a request to invoke SignAgreement API
func CreateSignAgreementRequest() (request *SignAgreementRequest) {
	request = &SignAgreementRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ecs", "2014-05-26", "SignAgreement", "ecs", "openAPI")
	return
}

// CreateSignAgreementResponse creates a response to parse from SignAgreement response
func CreateSignAgreementResponse() (response *SignAgreementResponse) {
	response = &SignAgreementResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
