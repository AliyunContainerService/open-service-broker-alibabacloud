package qualitycheck

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

// DeleteDataSet invokes the qualitycheck.DeleteDataSet API synchronously
// api document: https://help.aliyun.com/api/qualitycheck/deletedataset.html
func (client *Client) DeleteDataSet(request *DeleteDataSetRequest) (response *DeleteDataSetResponse, err error) {
	response = CreateDeleteDataSetResponse()
	err = client.DoAction(request, response)
	return
}

// DeleteDataSetWithChan invokes the qualitycheck.DeleteDataSet API asynchronously
// api document: https://help.aliyun.com/api/qualitycheck/deletedataset.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteDataSetWithChan(request *DeleteDataSetRequest) (<-chan *DeleteDataSetResponse, <-chan error) {
	responseChan := make(chan *DeleteDataSetResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DeleteDataSet(request)
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

// DeleteDataSetWithCallback invokes the qualitycheck.DeleteDataSet API asynchronously
// api document: https://help.aliyun.com/api/qualitycheck/deletedataset.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteDataSetWithCallback(request *DeleteDataSetRequest, callback func(response *DeleteDataSetResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DeleteDataSetResponse
		var err error
		defer close(result)
		response, err = client.DeleteDataSet(request)
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

// DeleteDataSetRequest is the request struct for api DeleteDataSet
type DeleteDataSetRequest struct {
	*requests.RpcRequest
	ResourceOwnerId requests.Integer `position:"Query" name:"ResourceOwnerId"`
	JsonStr         string           `position:"Query" name:"JsonStr"`
}

// DeleteDataSetResponse is the response struct for api DeleteDataSet
type DeleteDataSetResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	Success   bool   `json:"Success" xml:"Success"`
	Code      string `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
}

// CreateDeleteDataSetRequest creates a request to invoke DeleteDataSet API
func CreateDeleteDataSetRequest() (request *DeleteDataSetRequest) {
	request = &DeleteDataSetRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Qualitycheck", "2019-01-15", "DeleteDataSet", "", "")
	return
}

// CreateDeleteDataSetResponse creates a response to parse from DeleteDataSet response
func CreateDeleteDataSetResponse() (response *DeleteDataSetResponse) {
	response = &DeleteDataSetResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
