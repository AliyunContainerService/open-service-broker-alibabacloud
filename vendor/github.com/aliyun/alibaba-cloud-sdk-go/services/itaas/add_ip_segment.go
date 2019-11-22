package itaas

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

// AddIPSegment invokes the itaas.AddIPSegment API synchronously
// api document: https://help.aliyun.com/api/itaas/addipsegment.html
func (client *Client) AddIPSegment(request *AddIPSegmentRequest) (response *AddIPSegmentResponse, err error) {
	response = CreateAddIPSegmentResponse()
	err = client.DoAction(request, response)
	return
}

// AddIPSegmentWithChan invokes the itaas.AddIPSegment API asynchronously
// api document: https://help.aliyun.com/api/itaas/addipsegment.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) AddIPSegmentWithChan(request *AddIPSegmentRequest) (<-chan *AddIPSegmentResponse, <-chan error) {
	responseChan := make(chan *AddIPSegmentResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.AddIPSegment(request)
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

// AddIPSegmentWithCallback invokes the itaas.AddIPSegment API asynchronously
// api document: https://help.aliyun.com/api/itaas/addipsegment.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) AddIPSegmentWithCallback(request *AddIPSegmentRequest, callback func(response *AddIPSegmentResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *AddIPSegmentResponse
		var err error
		defer close(result)
		response, err = client.AddIPSegment(request)
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

// AddIPSegmentRequest is the request struct for api AddIPSegment
type AddIPSegmentRequest struct {
	*requests.RpcRequest
	Clientappid string `position:"Query" name:"Clientappid"`
	Ipsegment   string `position:"Query" name:"Ipsegment"`
	Memo        string `position:"Query" name:"Memo"`
	Sysfrom     string `position:"Query" name:"Sysfrom"`
	Operator    string `position:"Query" name:"Operator"`
}

// AddIPSegmentResponse is the response struct for api AddIPSegment
type AddIPSegmentResponse struct {
	*responses.BaseResponse
	RequestId string                  `json:"RequestId" xml:"RequestId"`
	ErrorCode int                     `json:"ErrorCode" xml:"ErrorCode"`
	ErrorMsg  string                  `json:"ErrorMsg" xml:"ErrorMsg"`
	Success   bool                    `json:"Success" xml:"Success"`
	ErrorList ErrorListInAddIPSegment `json:"ErrorList" xml:"ErrorList"`
}

// CreateAddIPSegmentRequest creates a request to invoke AddIPSegment API
func CreateAddIPSegmentRequest() (request *AddIPSegmentRequest) {
	request = &AddIPSegmentRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("ITaaS", "2017-05-05", "AddIPSegment", "itaas", "openAPI")
	return
}

// CreateAddIPSegmentResponse creates a response to parse from AddIPSegment response
func CreateAddIPSegmentResponse() (response *AddIPSegmentResponse) {
	response = &AddIPSegmentResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
