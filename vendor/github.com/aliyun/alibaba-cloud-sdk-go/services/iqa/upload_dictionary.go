package iqa

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

// UploadDictionary invokes the iqa.UploadDictionary API synchronously
// api document: https://help.aliyun.com/api/iqa/uploaddictionary.html
func (client *Client) UploadDictionary(request *UploadDictionaryRequest) (response *UploadDictionaryResponse, err error) {
	response = CreateUploadDictionaryResponse()
	err = client.DoAction(request, response)
	return
}

// UploadDictionaryWithChan invokes the iqa.UploadDictionary API asynchronously
// api document: https://help.aliyun.com/api/iqa/uploaddictionary.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) UploadDictionaryWithChan(request *UploadDictionaryRequest) (<-chan *UploadDictionaryResponse, <-chan error) {
	responseChan := make(chan *UploadDictionaryResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.UploadDictionary(request)
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

// UploadDictionaryWithCallback invokes the iqa.UploadDictionary API asynchronously
// api document: https://help.aliyun.com/api/iqa/uploaddictionary.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) UploadDictionaryWithCallback(request *UploadDictionaryRequest, callback func(response *UploadDictionaryResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *UploadDictionaryResponse
		var err error
		defer close(result)
		response, err = client.UploadDictionary(request)
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

// UploadDictionaryRequest is the request struct for api UploadDictionary
type UploadDictionaryRequest struct {
	*requests.RpcRequest
	DictionaryFileUrl string `position:"Body" name:"DictionaryFileUrl"`
	ProjectId         string `position:"Body" name:"ProjectId"`
	DictionaryData    string `position:"Body" name:"DictionaryData"`
}

// UploadDictionaryResponse is the response struct for api UploadDictionary
type UploadDictionaryResponse struct {
	*responses.BaseResponse
	RequestId     string `json:"RequestId" xml:"RequestId"`
	ProjectId     string `json:"ProjectId" xml:"ProjectId"`
	TotalCount    int    `json:"TotalCount" xml:"TotalCount"`
	FileDataCount int    `json:"FileDataCount" xml:"FileDataCount"`
	JsonDataCount int    `json:"JsonDataCount" xml:"JsonDataCount"`
}

// CreateUploadDictionaryRequest creates a request to invoke UploadDictionary API
func CreateUploadDictionaryRequest() (request *UploadDictionaryRequest) {
	request = &UploadDictionaryRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("iqa", "2019-08-13", "UploadDictionary", "iqa", "openAPI")
	return
}

// CreateUploadDictionaryResponse creates a response to parse from UploadDictionary response
func CreateUploadDictionaryResponse() (response *UploadDictionaryResponse) {
	response = &UploadDictionaryResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
