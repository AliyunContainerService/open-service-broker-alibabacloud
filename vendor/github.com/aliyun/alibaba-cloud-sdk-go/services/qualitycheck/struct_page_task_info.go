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

// PageTaskInfo is a nested struct in qualitycheck response
type PageTaskInfo struct {
	JobName          string           `json:"JobName" xml:"JobName"`
	ScheduleRatio    float64          `json:"ScheduleRatio" xml:"ScheduleRatio"`
	TaskId           string           `json:"TaskId" xml:"TaskId"`
	TaskComplete     bool             `json:"TaskComplete" xml:"TaskComplete"`
	Status           int              `json:"Status" xml:"Status"`
	IsTaskComplete   bool             `json:"IsTaskComplete" xml:"IsTaskComplete"`
	RuleSize         int              `json:"RuleSize" xml:"RuleSize"`
	DataSetSize      int              `json:"DataSetSize" xml:"DataSetSize"`
	DataSets         DataSets         `json:"DataSets" xml:"DataSets"`
	RuleNameInfoList RuleNameInfoList `json:"RuleNameInfoList" xml:"RuleNameInfoList"`
}
