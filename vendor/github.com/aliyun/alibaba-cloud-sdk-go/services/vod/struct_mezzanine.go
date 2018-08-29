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

// Mezzanine is a nested struct in vod response
type Mezzanine struct {
	PreprocessStatus string        `json:"PreprocessStatus" xml:"PreprocessStatus"`
	Fps              string        `json:"Fps" xml:"Fps"`
	Bitrate          string        `json:"Bitrate" xml:"Bitrate"`
	CreationTime     string        `json:"CreationTime" xml:"CreationTime"`
	CRC64            string        `json:"CRC64" xml:"CRC64"`
	OriginalFileName string        `json:"OriginalFileName" xml:"OriginalFileName"`
	Height           int           `json:"Height" xml:"Height"`
	FileURL          string        `json:"FileURL" xml:"FileURL"`
	Duration         string        `json:"Duration" xml:"Duration"`
	Width            int           `json:"Width" xml:"Width"`
	Size             int           `json:"Size" xml:"Size"`
	Status           string        `json:"Status" xml:"Status"`
	FileName         string        `json:"FileName" xml:"FileName"`
	VideoId          string        `json:"VideoId" xml:"VideoId"`
	VideoStreamList  []VideoStream `json:"VideoStreamList" xml:"VideoStreamList"`
	AudioStreamList  []AudioStream `json:"AudioStreamList" xml:"AudioStreamList"`
}
