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

// DescribeFlowNodeInstanceLauncherLog invokes the emr.DescribeFlowNodeInstanceLauncherLog API synchronously
// api document: https://help.aliyun.com/api/emr/describeflownodeinstancelauncherlog.html
func (client *Client) DescribeFlowNodeInstanceLauncherLog(request *DescribeFlowNodeInstanceLauncherLogRequest) (response *DescribeFlowNodeInstanceLauncherLogResponse, err error) {
	response = CreateDescribeFlowNodeInstanceLauncherLogResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeFlowNodeInstanceLauncherLogWithChan invokes the emr.DescribeFlowNodeInstanceLauncherLog API asynchronously
// api document: https://help.aliyun.com/api/emr/describeflownodeinstancelauncherlog.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeFlowNodeInstanceLauncherLogWithChan(request *DescribeFlowNodeInstanceLauncherLogRequest) (<-chan *DescribeFlowNodeInstanceLauncherLogResponse, <-chan error) {
	responseChan := make(chan *DescribeFlowNodeInstanceLauncherLogResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeFlowNodeInstanceLauncherLog(request)
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

// DescribeFlowNodeInstanceLauncherLogWithCallback invokes the emr.DescribeFlowNodeInstanceLauncherLog API asynchronously
// api document: https://help.aliyun.com/api/emr/describeflownodeinstancelauncherlog.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeFlowNodeInstanceLauncherLogWithCallback(request *DescribeFlowNodeInstanceLauncherLogRequest, callback func(response *DescribeFlowNodeInstanceLauncherLogResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeFlowNodeInstanceLauncherLogResponse
		var err error
		defer close(result)
		response, err = client.DescribeFlowNodeInstanceLauncherLog(request)
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

// DescribeFlowNodeInstanceLauncherLogRequest is the request struct for api DescribeFlowNodeInstanceLauncherLog
type DescribeFlowNodeInstanceLauncherLogRequest struct {
	*requests.RpcRequest
	Offset         requests.Integer `position:"Query" name:"Offset"`
	Start          requests.Integer `position:"Query" name:"Start"`
	Length         requests.Integer `position:"Query" name:"Length"`
	EndTime        requests.Integer `position:"Query" name:"EndTime"`
	StartTime      requests.Integer `position:"Query" name:"StartTime"`
	Lines          requests.Integer `position:"Query" name:"Lines"`
	Reverse        requests.Boolean `position:"Query" name:"Reverse"`
	NodeInstanceId string           `position:"Query" name:"NodeInstanceId"`
	ProjectId      string           `position:"Query" name:"ProjectId"`
}

// DescribeFlowNodeInstanceLauncherLogResponse is the response struct for api DescribeFlowNodeInstanceLauncherLog
type DescribeFlowNodeInstanceLauncherLogResponse struct {
	*responses.BaseResponse
	RequestId string                                         `json:"RequestId" xml:"RequestId"`
	LogEnd    bool                                           `json:"LogEnd" xml:"LogEnd"`
	LogEntrys LogEntrysInDescribeFlowNodeInstanceLauncherLog `json:"LogEntrys" xml:"LogEntrys"`
}

// CreateDescribeFlowNodeInstanceLauncherLogRequest creates a request to invoke DescribeFlowNodeInstanceLauncherLog API
func CreateDescribeFlowNodeInstanceLauncherLogRequest() (request *DescribeFlowNodeInstanceLauncherLogRequest) {
	request = &DescribeFlowNodeInstanceLauncherLogRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Emr", "2016-04-08", "DescribeFlowNodeInstanceLauncherLog", "emr", "openAPI")
	return
}

// CreateDescribeFlowNodeInstanceLauncherLogResponse creates a response to parse from DescribeFlowNodeInstanceLauncherLog response
func CreateDescribeFlowNodeInstanceLauncherLogResponse() (response *DescribeFlowNodeInstanceLauncherLogResponse) {
	response = &DescribeFlowNodeInstanceLauncherLogResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
