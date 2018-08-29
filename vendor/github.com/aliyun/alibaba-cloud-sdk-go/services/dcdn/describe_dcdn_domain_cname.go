package dcdn

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

// DescribeDcdnDomainCname invokes the dcdn.DescribeDcdnDomainCname API synchronously
// api document: https://help.aliyun.com/api/dcdn/describedcdndomaincname.html
func (client *Client) DescribeDcdnDomainCname(request *DescribeDcdnDomainCnameRequest) (response *DescribeDcdnDomainCnameResponse, err error) {
	response = CreateDescribeDcdnDomainCnameResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDcdnDomainCnameWithChan invokes the dcdn.DescribeDcdnDomainCname API asynchronously
// api document: https://help.aliyun.com/api/dcdn/describedcdndomaincname.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDcdnDomainCnameWithChan(request *DescribeDcdnDomainCnameRequest) (<-chan *DescribeDcdnDomainCnameResponse, <-chan error) {
	responseChan := make(chan *DescribeDcdnDomainCnameResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDcdnDomainCname(request)
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

// DescribeDcdnDomainCnameWithCallback invokes the dcdn.DescribeDcdnDomainCname API asynchronously
// api document: https://help.aliyun.com/api/dcdn/describedcdndomaincname.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDcdnDomainCnameWithCallback(request *DescribeDcdnDomainCnameRequest, callback func(response *DescribeDcdnDomainCnameResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDcdnDomainCnameResponse
		var err error
		defer close(result)
		response, err = client.DescribeDcdnDomainCname(request)
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

// DescribeDcdnDomainCnameRequest is the request struct for api DescribeDcdnDomainCname
type DescribeDcdnDomainCnameRequest struct {
	*requests.RpcRequest
}

// DescribeDcdnDomainCnameResponse is the response struct for api DescribeDcdnDomainCname
type DescribeDcdnDomainCnameResponse struct {
	*responses.BaseResponse
	RequestId  string     `json:"RequestId" xml:"RequestId"`
	CnameDatas CnameDatas `json:"CnameDatas" xml:"CnameDatas"`
}

// CreateDescribeDcdnDomainCnameRequest creates a request to invoke DescribeDcdnDomainCname API
func CreateDescribeDcdnDomainCnameRequest() (request *DescribeDcdnDomainCnameRequest) {
	request = &DescribeDcdnDomainCnameRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("dcdn", "2018-01-15", "DescribeDcdnDomainCname", "dcdn", "openAPI")
	return
}

// CreateDescribeDcdnDomainCnameResponse creates a response to parse from DescribeDcdnDomainCname response
func CreateDescribeDcdnDomainCnameResponse() (response *DescribeDcdnDomainCnameResponse) {
	response = &DescribeDcdnDomainCnameResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
