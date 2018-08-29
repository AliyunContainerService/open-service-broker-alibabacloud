package teslamaxcompute

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

// QueryResourceInventory invokes the teslamaxcompute.QueryResourceInventory API synchronously
// api document: https://help.aliyun.com/api/teslamaxcompute/queryresourceinventory.html
func (client *Client) QueryResourceInventory(request *QueryResourceInventoryRequest) (response *QueryResourceInventoryResponse, err error) {
	response = CreateQueryResourceInventoryResponse()
	err = client.DoAction(request, response)
	return
}

// QueryResourceInventoryWithChan invokes the teslamaxcompute.QueryResourceInventory API asynchronously
// api document: https://help.aliyun.com/api/teslamaxcompute/queryresourceinventory.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) QueryResourceInventoryWithChan(request *QueryResourceInventoryRequest) (<-chan *QueryResourceInventoryResponse, <-chan error) {
	responseChan := make(chan *QueryResourceInventoryResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.QueryResourceInventory(request)
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

// QueryResourceInventoryWithCallback invokes the teslamaxcompute.QueryResourceInventory API asynchronously
// api document: https://help.aliyun.com/api/teslamaxcompute/queryresourceinventory.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) QueryResourceInventoryWithCallback(request *QueryResourceInventoryRequest, callback func(response *QueryResourceInventoryResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *QueryResourceInventoryResponse
		var err error
		defer close(result)
		response, err = client.QueryResourceInventory(request)
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

// QueryResourceInventoryRequest is the request struct for api QueryResourceInventory
type QueryResourceInventoryRequest struct {
	*requests.RpcRequest
}

// QueryResourceInventoryResponse is the response struct for api QueryResourceInventory
type QueryResourceInventoryResponse struct {
	*responses.BaseResponse
	Code      int    `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
	RequestId string `json:"RequestId" xml:"RequestId"`
	Data      Data   `json:"Data" xml:"Data"`
}

// CreateQueryResourceInventoryRequest creates a request to invoke QueryResourceInventory API
func CreateQueryResourceInventoryRequest() (request *QueryResourceInventoryRequest) {
	request = &QueryResourceInventoryRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("TeslaMaxCompute", "2018-01-04", "QueryResourceInventory", "", "")
	return
}

// CreateQueryResourceInventoryResponse creates a response to parse from QueryResourceInventory response
func CreateQueryResourceInventoryResponse() (response *QueryResourceInventoryResponse) {
	response = &QueryResourceInventoryResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
