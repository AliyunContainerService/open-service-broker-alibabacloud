package drds

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

// DescribeDrdsInstanceDbMonitor invokes the drds.DescribeDrdsInstanceDbMonitor API synchronously
// api document: https://help.aliyun.com/api/drds/describedrdsinstancedbmonitor.html
func (client *Client) DescribeDrdsInstanceDbMonitor(request *DescribeDrdsInstanceDbMonitorRequest) (response *DescribeDrdsInstanceDbMonitorResponse, err error) {
	response = CreateDescribeDrdsInstanceDbMonitorResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDrdsInstanceDbMonitorWithChan invokes the drds.DescribeDrdsInstanceDbMonitor API asynchronously
// api document: https://help.aliyun.com/api/drds/describedrdsinstancedbmonitor.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDrdsInstanceDbMonitorWithChan(request *DescribeDrdsInstanceDbMonitorRequest) (<-chan *DescribeDrdsInstanceDbMonitorResponse, <-chan error) {
	responseChan := make(chan *DescribeDrdsInstanceDbMonitorResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDrdsInstanceDbMonitor(request)
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

// DescribeDrdsInstanceDbMonitorWithCallback invokes the drds.DescribeDrdsInstanceDbMonitor API asynchronously
// api document: https://help.aliyun.com/api/drds/describedrdsinstancedbmonitor.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDrdsInstanceDbMonitorWithCallback(request *DescribeDrdsInstanceDbMonitorRequest, callback func(response *DescribeDrdsInstanceDbMonitorResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDrdsInstanceDbMonitorResponse
		var err error
		defer close(result)
		response, err = client.DescribeDrdsInstanceDbMonitor(request)
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

// DescribeDrdsInstanceDbMonitorRequest is the request struct for api DescribeDrdsInstanceDbMonitor
type DescribeDrdsInstanceDbMonitorRequest struct {
	*requests.RpcRequest
	DbName         string           `position:"Query" name:"DbName"`
	EndTime        requests.Integer `position:"Query" name:"EndTime"`
	StartTime      requests.Integer `position:"Query" name:"StartTime"`
	DrdsInstanceId string           `position:"Query" name:"DrdsInstanceId"`
	Key            string           `position:"Query" name:"Key"`
}

// DescribeDrdsInstanceDbMonitorResponse is the response struct for api DescribeDrdsInstanceDbMonitor
type DescribeDrdsInstanceDbMonitorResponse struct {
	*responses.BaseResponse
	RequestId string                              `json:"RequestId" xml:"RequestId"`
	Success   bool                                `json:"Success" xml:"Success"`
	Data      DataInDescribeDrdsInstanceDbMonitor `json:"Data" xml:"Data"`
}

// CreateDescribeDrdsInstanceDbMonitorRequest creates a request to invoke DescribeDrdsInstanceDbMonitor API
func CreateDescribeDrdsInstanceDbMonitorRequest() (request *DescribeDrdsInstanceDbMonitorRequest) {
	request = &DescribeDrdsInstanceDbMonitorRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Drds", "2017-10-16", "DescribeDrdsInstanceDbMonitor", "drds", "openAPI")
	return
}

// CreateDescribeDrdsInstanceDbMonitorResponse creates a response to parse from DescribeDrdsInstanceDbMonitor response
func CreateDescribeDrdsInstanceDbMonitorResponse() (response *DescribeDrdsInstanceDbMonitorResponse) {
	response = &DescribeDrdsInstanceDbMonitorResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
