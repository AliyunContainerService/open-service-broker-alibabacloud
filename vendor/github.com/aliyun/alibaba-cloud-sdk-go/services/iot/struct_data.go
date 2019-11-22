package iot

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

// Data is a nested struct in iot response
type Data struct {
	DataFormat             int                   `json:"DataFormat" xml:"DataFormat"`
	BizEnable              bool                  `json:"BizEnable" xml:"BizEnable"`
	LatestDeploymentStatus int                   `json:"LatestDeploymentStatus" xml:"LatestDeploymentStatus"`
	RoleAttachTime         string                `json:"RoleAttachTime" xml:"RoleAttachTime"`
	RequestProtocol        string                `json:"RequestProtocol" xml:"RequestProtocol"`
	GmtCompleted           string                `json:"GmtCompleted" xml:"GmtCompleted"`
	UtcCreate              string                `json:"UtcCreate" xml:"UtcCreate"`
	RoleName               string                `json:"RoleName" xml:"RoleName"`
	Spec                   int                   `json:"Spec" xml:"Spec"`
	DeviceActive           int                   `json:"DeviceActive" xml:"DeviceActive"`
	RequestMethod          string                `json:"RequestMethod" xml:"RequestMethod"`
	Nickname               string                `json:"Nickname" xml:"Nickname"`
	PageNo                 int                   `json:"PageNo" xml:"PageNo"`
	DevEui                 string                `json:"DevEui" xml:"DevEui"`
	GroupId                string                `json:"GroupId" xml:"GroupId"`
	LatestDeploymentType   string                `json:"LatestDeploymentType" xml:"LatestDeploymentType"`
	RoleArn                string                `json:"RoleArn" xml:"RoleArn"`
	Type                   string                `json:"Type" xml:"Type"`
	FileId                 string                `json:"FileId" xml:"FileId"`
	LastUpdateTime         int64                 `json:"LastUpdateTime" xml:"LastUpdateTime"`
	Tags                   string                `json:"Tags" xml:"Tags"`
	Versions               string                `json:"Versions" xml:"Versions"`
	AliyunCommodityCode    string                `json:"AliyunCommodityCode" xml:"AliyunCommodityCode"`
	ApplyId                int64                 `json:"ApplyId" xml:"ApplyId"`
	UtcCreatedOn           string                `json:"UtcCreatedOn" xml:"UtcCreatedOn"`
	MessageId              string                `json:"MessageId" xml:"MessageId"`
	DeviceName             string                `json:"DeviceName" xml:"DeviceName"`
	PageCount              int64                 `json:"PageCount" xml:"PageCount"`
	Size                   string                `json:"Size" xml:"Size"`
	Id2                    bool                  `json:"Id2" xml:"Id2"`
	NodeType               int                   `json:"NodeType" xml:"NodeType"`
	ApiSrn                 string                `json:"ApiSrn" xml:"ApiSrn"`
	ProductName            string                `json:"ProductName" xml:"ProductName"`
	Name                   string                `json:"Name" xml:"Name"`
	GroupName              string                `json:"GroupName" xml:"GroupName"`
	DownloadUrl            string                `json:"DownloadUrl" xml:"DownloadUrl"`
	CreateTime             int64                 `json:"CreateTime" xml:"CreateTime"`
	DeploymentId           string                `json:"DeploymentId" xml:"DeploymentId"`
	PageSize               int                   `json:"PageSize" xml:"PageSize"`
	GmtCreate              string                `json:"GmtCreate" xml:"GmtCreate"`
	InstanceId             string                `json:"InstanceId" xml:"InstanceId"`
	Description            string                `json:"Description" xml:"Description"`
	DateFormat             string                `json:"DateFormat" xml:"DateFormat"`
	ApiPath                string                `json:"ApiPath" xml:"ApiPath"`
	DeviceOnline           int                   `json:"DeviceOnline" xml:"DeviceOnline"`
	Status                 int                   `json:"Status" xml:"Status"`
	Result                 string                `json:"Result" xml:"Result"`
	DeviceSecret           string                `json:"DeviceSecret" xml:"DeviceSecret"`
	ProductKey             string                `json:"ProductKey" xml:"ProductKey"`
	GmtModified            string                `json:"GmtModified" xml:"GmtModified"`
	DisplayName            string                `json:"DisplayName" xml:"DisplayName"`
	JoinEui                string                `json:"JoinEui" xml:"JoinEui"`
	CurrentPage            int                   `json:"CurrentPage" xml:"CurrentPage"`
	IotId                  string                `json:"IotId" xml:"IotId"`
	GroupDesc              string                `json:"GroupDesc" xml:"GroupDesc"`
	DeviceCount            int                   `json:"DeviceCount" xml:"DeviceCount"`
	ProtocolType           string                `json:"ProtocolType" xml:"ProtocolType"`
	AuthType               string                `json:"AuthType" xml:"AuthType"`
	Total                  int64                 `json:"Total" xml:"Total"`
	FieldNameList          FieldNameList         `json:"FieldNameList" xml:"FieldNameList"`
	InvalidDeviceNameList  InvalidDeviceNameList `json:"InvalidDeviceNameList" xml:"InvalidDeviceNameList"`
	ResultList             ResultList            `json:"ResultList" xml:"ResultList"`
	SqlTemplateDTO         SqlTemplateDTO        `json:"SqlTemplateDTO" xml:"SqlTemplateDTO"`
	TaskList               []Task                `json:"TaskList" xml:"TaskList"`
	List                   ListInGetThingTopo    `json:"List" xml:"List"`
}
