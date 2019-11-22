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

// ListETLJobRelease invokes the emr.ListETLJobRelease API synchronously
// api document: https://help.aliyun.com/api/emr/listetljobrelease.html
func (client *Client) ListETLJobRelease(request *ListETLJobReleaseRequest) (response *ListETLJobReleaseResponse, err error) {
	response = CreateListETLJobReleaseResponse()
	err = client.DoAction(request, response)
	return
}

// ListETLJobReleaseWithChan invokes the emr.ListETLJobRelease API asynchronously
// api document: https://help.aliyun.com/api/emr/listetljobrelease.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ListETLJobReleaseWithChan(request *ListETLJobReleaseRequest) (<-chan *ListETLJobReleaseResponse, <-chan error) {
	responseChan := make(chan *ListETLJobReleaseResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ListETLJobRelease(request)
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

// ListETLJobReleaseWithCallback invokes the emr.ListETLJobRelease API asynchronously
// api document: https://help.aliyun.com/api/emr/listetljobrelease.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ListETLJobReleaseWithCallback(request *ListETLJobReleaseRequest, callback func(response *ListETLJobReleaseResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ListETLJobReleaseResponse
		var err error
		defer close(result)
		response, err = client.ListETLJobRelease(request)
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

// ListETLJobReleaseRequest is the request struct for api ListETLJobRelease
type ListETLJobReleaseRequest struct {
	*requests.RpcRequest
	ResourceOwnerId requests.Integer `position:"Query" name:"ResourceOwnerId"`
	EtlJobId        string           `position:"Query" name:"EtlJobId"`
	ReleaseId       string           `position:"Query" name:"ReleaseId"`
	PageSize        requests.Integer `position:"Query" name:"PageSize"`
	PageNumber      requests.Integer `position:"Query" name:"PageNumber"`
	ReleaseVersion  requests.Integer `position:"Query" name:"ReleaseVersion"`
	Status          string           `position:"Query" name:"Status"`
}

// ListETLJobReleaseResponse is the response struct for api ListETLJobRelease
type ListETLJobReleaseResponse struct {
	*responses.BaseResponse
	RequestId   string      `json:"RequestId" xml:"RequestId"`
	Total       int         `json:"Total" xml:"Total"`
	PageSize    int         `json:"PageSize" xml:"PageSize"`
	PageNumber  int         `json:"PageNumber" xml:"PageNumber"`
	ReleaseList ReleaseList `json:"ReleaseList" xml:"ReleaseList"`
}

// CreateListETLJobReleaseRequest creates a request to invoke ListETLJobRelease API
func CreateListETLJobReleaseRequest() (request *ListETLJobReleaseRequest) {
	request = &ListETLJobReleaseRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Emr", "2016-04-08", "ListETLJobRelease", "emr", "openAPI")
	return
}

// CreateListETLJobReleaseResponse creates a response to parse from ListETLJobRelease response
func CreateListETLJobReleaseResponse() (response *ListETLJobReleaseResponse) {
	response = &ListETLJobReleaseResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
