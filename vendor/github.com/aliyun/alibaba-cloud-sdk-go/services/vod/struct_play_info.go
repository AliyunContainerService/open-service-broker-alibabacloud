package vod

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

// PlayInfo is a nested struct in vod response
type PlayInfo struct {
	Region           string `json:"Region" xml:"Region"`
	JobId            string `json:"JobId" xml:"JobId"`
	Format           string `json:"Format" xml:"Format"`
	PreprocessStatus string `json:"PreprocessStatus" xml:"PreprocessStatus"`
	Fps              string `json:"Fps" xml:"Fps"`
	Bitrate          string `json:"Bitrate" xml:"Bitrate"`
	Encrypt          int    `json:"Encrypt" xml:"Encrypt"`
	Rand             string `json:"Rand" xml:"Rand"`
	StreamType       string `json:"StreamType" xml:"StreamType"`
	AccessKeyId      string `json:"AccessKeyId" xml:"AccessKeyId"`
	Height           int    `json:"Height" xml:"Height"`
	AccessKeySecret  string `json:"AccessKeySecret" xml:"AccessKeySecret"`
	PlayDomain       string `json:"PlayDomain" xml:"PlayDomain"`
	WatermarkId      string `json:"WatermarkId" xml:"WatermarkId"`
	Duration         string `json:"Duration" xml:"Duration"`
	Complexity       string `json:"Complexity" xml:"Complexity"`
	Width            int    `json:"Width" xml:"Width"`
	AuthInfo         string `json:"AuthInfo" xml:"AuthInfo"`
	Size             int    `json:"Size" xml:"Size"`
	Status           string `json:"Status" xml:"Status"`
	Definition       string `json:"Definition" xml:"Definition"`
	Plaintext        string `json:"Plaintext" xml:"Plaintext"`
	PlayURL          string `json:"PlayURL" xml:"PlayURL"`
	SecurityToken    string `json:"SecurityToken" xml:"SecurityToken"`
}
