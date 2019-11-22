package cloudauth

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

// DescribeRPSDK invokes the cloudauth.DescribeRPSDK API synchronously
// api document: https://help.aliyun.com/api/cloudauth/describerpsdk.html
func (client *Client) DescribeRPSDK(request *DescribeRPSDKRequest) (response *DescribeRPSDKResponse, err error) {
	response = CreateDescribeRPSDKResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeRPSDKWithChan invokes the cloudauth.DescribeRPSDK API asynchronously
// api document: https://help.aliyun.com/api/cloudauth/describerpsdk.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeRPSDKWithChan(request *DescribeRPSDKRequest) (<-chan *DescribeRPSDKResponse, <-chan error) {
	responseChan := make(chan *DescribeRPSDKResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeRPSDK(request)
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

// DescribeRPSDKWithCallback invokes the cloudauth.DescribeRPSDK API asynchronously
// api document: https://help.aliyun.com/api/cloudauth/describerpsdk.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeRPSDKWithCallback(request *DescribeRPSDKRequest, callback func(response *DescribeRPSDKResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeRPSDKResponse
		var err error
		defer close(result)
		response, err = client.DescribeRPSDK(request)
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

// DescribeRPSDKRequest is the request struct for api DescribeRPSDK
type DescribeRPSDKRequest struct {
	*requests.RpcRequest
	SourceIp string `position:"Query" name:"SourceIp"`
	Lang     string `position:"Query" name:"Lang"`
	TaskId   string `position:"Query" name:"TaskId"`
}

// DescribeRPSDKResponse is the response struct for api DescribeRPSDK
type DescribeRPSDKResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	SdkUrl    string `json:"SdkUrl" xml:"SdkUrl"`
}

// CreateDescribeRPSDKRequest creates a request to invoke DescribeRPSDK API
func CreateDescribeRPSDKRequest() (request *DescribeRPSDKRequest) {
	request = &DescribeRPSDKRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cloudauth", "2019-03-07", "DescribeRPSDK", "cloudauth", "openAPI")
	return
}

// CreateDescribeRPSDKResponse creates a response to parse from DescribeRPSDK response
func CreateDescribeRPSDKResponse() (response *DescribeRPSDKResponse) {
	response = &DescribeRPSDKResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
