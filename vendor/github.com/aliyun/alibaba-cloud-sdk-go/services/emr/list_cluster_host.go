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

// ListClusterHost invokes the emr.ListClusterHost API synchronously
// api document: https://help.aliyun.com/api/emr/listclusterhost.html
func (client *Client) ListClusterHost(request *ListClusterHostRequest) (response *ListClusterHostResponse, err error) {
	response = CreateListClusterHostResponse()
	err = client.DoAction(request, response)
	return
}

// ListClusterHostWithChan invokes the emr.ListClusterHost API asynchronously
// api document: https://help.aliyun.com/api/emr/listclusterhost.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ListClusterHostWithChan(request *ListClusterHostRequest) (<-chan *ListClusterHostResponse, <-chan error) {
	responseChan := make(chan *ListClusterHostResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ListClusterHost(request)
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

// ListClusterHostWithCallback invokes the emr.ListClusterHost API asynchronously
// api document: https://help.aliyun.com/api/emr/listclusterhost.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ListClusterHostWithCallback(request *ListClusterHostRequest, callback func(response *ListClusterHostResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ListClusterHostResponse
		var err error
		defer close(result)
		response, err = client.ListClusterHost(request)
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

// ListClusterHostRequest is the request struct for api ListClusterHost
type ListClusterHostRequest struct {
	*requests.RpcRequest
	ResourceOwnerId requests.Integer `position:"Query" name:"ResourceOwnerId"`
	HostInstanceId  string           `position:"Query" name:"HostInstanceId"`
	StatusList      *[]string        `position:"Query" name:"StatusList"  type:"Repeated"`
	PrivateIp       string           `position:"Query" name:"PrivateIp"`
	ComponentName   string           `position:"Query" name:"ComponentName"`
	PublicIp        string           `position:"Query" name:"PublicIp"`
	ClusterId       string           `position:"Query" name:"ClusterId"`
	PageNumber      requests.Integer `position:"Query" name:"PageNumber"`
	HostName        string           `position:"Query" name:"HostName"`
	GroupType       string           `position:"Query" name:"GroupType"`
	HostGroupId     string           `position:"Query" name:"HostGroupId"`
	PageSize        requests.Integer `position:"Query" name:"PageSize"`
}

// ListClusterHostResponse is the response struct for api ListClusterHost
type ListClusterHostResponse struct {
	*responses.BaseResponse
	RequestId  string                    `json:"RequestId" xml:"RequestId"`
	PageNumber int                       `json:"PageNumber" xml:"PageNumber"`
	PageSize   int                       `json:"PageSize" xml:"PageSize"`
	Total      int                       `json:"Total" xml:"Total"`
	HostList   HostListInListClusterHost `json:"HostList" xml:"HostList"`
}

// CreateListClusterHostRequest creates a request to invoke ListClusterHost API
func CreateListClusterHostRequest() (request *ListClusterHostRequest) {
	request = &ListClusterHostRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Emr", "2016-04-08", "ListClusterHost", "emr", "openAPI")
	return
}

// CreateListClusterHostResponse creates a response to parse from ListClusterHost response
func CreateListClusterHostResponse() (response *ListClusterHostResponse) {
	response = &ListClusterHostResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
