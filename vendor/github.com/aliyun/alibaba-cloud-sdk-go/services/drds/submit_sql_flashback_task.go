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

// SubmitSqlFlashbackTask invokes the drds.SubmitSqlFlashbackTask API synchronously
// api document: https://help.aliyun.com/api/drds/submitsqlflashbacktask.html
func (client *Client) SubmitSqlFlashbackTask(request *SubmitSqlFlashbackTaskRequest) (response *SubmitSqlFlashbackTaskResponse, err error) {
	response = CreateSubmitSqlFlashbackTaskResponse()
	err = client.DoAction(request, response)
	return
}

// SubmitSqlFlashbackTaskWithChan invokes the drds.SubmitSqlFlashbackTask API asynchronously
// api document: https://help.aliyun.com/api/drds/submitsqlflashbacktask.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SubmitSqlFlashbackTaskWithChan(request *SubmitSqlFlashbackTaskRequest) (<-chan *SubmitSqlFlashbackTaskResponse, <-chan error) {
	responseChan := make(chan *SubmitSqlFlashbackTaskResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.SubmitSqlFlashbackTask(request)
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

// SubmitSqlFlashbackTaskWithCallback invokes the drds.SubmitSqlFlashbackTask API asynchronously
// api document: https://help.aliyun.com/api/drds/submitsqlflashbacktask.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SubmitSqlFlashbackTaskWithCallback(request *SubmitSqlFlashbackTaskRequest, callback func(response *SubmitSqlFlashbackTaskResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *SubmitSqlFlashbackTaskResponse
		var err error
		defer close(result)
		response, err = client.SubmitSqlFlashbackTask(request)
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

// SubmitSqlFlashbackTaskRequest is the request struct for api SubmitSqlFlashbackTask
type SubmitSqlFlashbackTaskRequest struct {
	*requests.RpcRequest
	TraceId           string           `position:"Query" name:"TraceId"`
	SqlType           string           `position:"Query" name:"SqlType"`
	SqlPk             string           `position:"Query" name:"SqlPk"`
	RecallRestoreType requests.Integer `position:"Query" name:"RecallRestoreType"`
	RecallType        requests.Integer `position:"Query" name:"RecallType"`
	DbName            string           `position:"Query" name:"DbName"`
	EndTime           string           `position:"Query" name:"EndTime"`
	StartTime         string           `position:"Query" name:"StartTime"`
	TableName         string           `position:"Query" name:"TableName"`
	DrdsInstanceId    string           `position:"Query" name:"DrdsInstanceId"`
}

// SubmitSqlFlashbackTaskResponse is the response struct for api SubmitSqlFlashbackTask
type SubmitSqlFlashbackTaskResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	Success   bool   `json:"Success" xml:"Success"`
	TaskId    int64  `json:"TaskId" xml:"TaskId"`
}

// CreateSubmitSqlFlashbackTaskRequest creates a request to invoke SubmitSqlFlashbackTask API
func CreateSubmitSqlFlashbackTaskRequest() (request *SubmitSqlFlashbackTaskRequest) {
	request = &SubmitSqlFlashbackTaskRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Drds", "2019-01-23", "SubmitSqlFlashbackTask", "drds", "openAPI")
	return
}

// CreateSubmitSqlFlashbackTaskResponse creates a response to parse from SubmitSqlFlashbackTask response
func CreateSubmitSqlFlashbackTaskResponse() (response *SubmitSqlFlashbackTaskResponse) {
	response = &SubmitSqlFlashbackTaskResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
