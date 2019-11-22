package dyvmsapi

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

// UploadRobotTaskCalledFile invokes the dyvmsapi.UploadRobotTaskCalledFile API synchronously
// api document: https://help.aliyun.com/api/dyvmsapi/uploadrobottaskcalledfile.html
func (client *Client) UploadRobotTaskCalledFile(request *UploadRobotTaskCalledFileRequest) (response *UploadRobotTaskCalledFileResponse, err error) {
	response = CreateUploadRobotTaskCalledFileResponse()
	err = client.DoAction(request, response)
	return
}

// UploadRobotTaskCalledFileWithChan invokes the dyvmsapi.UploadRobotTaskCalledFile API asynchronously
// api document: https://help.aliyun.com/api/dyvmsapi/uploadrobottaskcalledfile.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) UploadRobotTaskCalledFileWithChan(request *UploadRobotTaskCalledFileRequest) (<-chan *UploadRobotTaskCalledFileResponse, <-chan error) {
	responseChan := make(chan *UploadRobotTaskCalledFileResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.UploadRobotTaskCalledFile(request)
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

// UploadRobotTaskCalledFileWithCallback invokes the dyvmsapi.UploadRobotTaskCalledFile API asynchronously
// api document: https://help.aliyun.com/api/dyvmsapi/uploadrobottaskcalledfile.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) UploadRobotTaskCalledFileWithCallback(request *UploadRobotTaskCalledFileRequest, callback func(response *UploadRobotTaskCalledFileResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *UploadRobotTaskCalledFileResponse
		var err error
		defer close(result)
		response, err = client.UploadRobotTaskCalledFile(request)
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

// UploadRobotTaskCalledFileRequest is the request struct for api UploadRobotTaskCalledFile
type UploadRobotTaskCalledFileRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	TtsParamHead         string           `position:"Query" name:"TtsParamHead"`
	TtsParam             string           `position:"Query" name:"TtsParam"`
	CalledNumber         string           `position:"Query" name:"CalledNumber"`
	Id                   requests.Integer `position:"Query" name:"Id"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
}

// UploadRobotTaskCalledFileResponse is the response struct for api UploadRobotTaskCalledFile
type UploadRobotTaskCalledFileResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	Data      string `json:"Data" xml:"Data"`
	Code      string `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
}

// CreateUploadRobotTaskCalledFileRequest creates a request to invoke UploadRobotTaskCalledFile API
func CreateUploadRobotTaskCalledFileRequest() (request *UploadRobotTaskCalledFileRequest) {
	request = &UploadRobotTaskCalledFileRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Dyvmsapi", "2017-05-25", "UploadRobotTaskCalledFile", "", "")
	return
}

// CreateUploadRobotTaskCalledFileResponse creates a response to parse from UploadRobotTaskCalledFile response
func CreateUploadRobotTaskCalledFileResponse() (response *UploadRobotTaskCalledFileResponse) {
	response = &UploadRobotTaskCalledFileResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
