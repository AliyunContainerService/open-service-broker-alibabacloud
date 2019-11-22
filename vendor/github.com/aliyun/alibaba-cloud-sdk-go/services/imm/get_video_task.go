package imm

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

// GetVideoTask invokes the imm.GetVideoTask API synchronously
// api document: https://help.aliyun.com/api/imm/getvideotask.html
func (client *Client) GetVideoTask(request *GetVideoTaskRequest) (response *GetVideoTaskResponse, err error) {
	response = CreateGetVideoTaskResponse()
	err = client.DoAction(request, response)
	return
}

// GetVideoTaskWithChan invokes the imm.GetVideoTask API asynchronously
// api document: https://help.aliyun.com/api/imm/getvideotask.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) GetVideoTaskWithChan(request *GetVideoTaskRequest) (<-chan *GetVideoTaskResponse, <-chan error) {
	responseChan := make(chan *GetVideoTaskResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.GetVideoTask(request)
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

// GetVideoTaskWithCallback invokes the imm.GetVideoTask API asynchronously
// api document: https://help.aliyun.com/api/imm/getvideotask.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) GetVideoTaskWithCallback(request *GetVideoTaskRequest, callback func(response *GetVideoTaskResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *GetVideoTaskResponse
		var err error
		defer close(result)
		response, err = client.GetVideoTask(request)
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

// GetVideoTaskRequest is the request struct for api GetVideoTask
type GetVideoTaskRequest struct {
	*requests.RpcRequest
	Project  string `position:"Query" name:"Project"`
	TaskId   string `position:"Query" name:"TaskId"`
	TaskType string `position:"Query" name:"TaskType"`
}

// GetVideoTaskResponse is the response struct for api GetVideoTask
type GetVideoTaskResponse struct {
	*responses.BaseResponse
	RequestId       string `json:"RequestId" xml:"RequestId"`
	TaskId          string `json:"TaskId" xml:"TaskId"`
	TaskType        string `json:"TaskType" xml:"TaskType"`
	Parameters      string `json:"Parameters" xml:"Parameters"`
	Result          string `json:"Result" xml:"Result"`
	Status          string `json:"Status" xml:"Status"`
	StartTime       string `json:"StartTime" xml:"StartTime"`
	EndTime         string `json:"EndTime" xml:"EndTime"`
	ErrorMessage    string `json:"ErrorMessage" xml:"ErrorMessage"`
	NotifyEndpoint  string `json:"NotifyEndpoint" xml:"NotifyEndpoint"`
	NotifyTopicName string `json:"NotifyTopicName" xml:"NotifyTopicName"`
	Progress        int    `json:"Progress" xml:"Progress"`
}

// CreateGetVideoTaskRequest creates a request to invoke GetVideoTask API
func CreateGetVideoTaskRequest() (request *GetVideoTaskRequest) {
	request = &GetVideoTaskRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("imm", "2017-09-06", "GetVideoTask", "imm", "openAPI")
	return
}

// CreateGetVideoTaskResponse creates a response to parse from GetVideoTask response
func CreateGetVideoTaskResponse() (response *GetVideoTaskResponse) {
	response = &GetVideoTaskResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
