// Copyright [2020] [Rafał Korepta]
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

type Response struct {
	Coord      Coord        `json:"coord,omitempty"`
	Weather    []Weather    `json:"weather,omitempty"`
	Main       MainMeasures `json:"main,omitempty"`
	Visibility int32        `json:"visibility,omitempty"`
	Wind       Wind         `json:"wind,omitempty"`
	Clouds     Clouds       `json:"clouds,omitempty"`
	Dt         int32        `json:"dt,omitempty"`
	Sys        Sys          `json:"sys,omitempty"`
	Timezone   int32        `json:"timezone,omitempty"`
	ID         int32        `json:"id,omitempty"`
	Name       string       `json:"name,omitempty"`
	Cod        int32        `json:"cod,omitempty"`
}

type Coord struct {
	Lon float32 `json:"lon,omitempty"`
	Lat float32 `json:"lat,omitempty"`
}

type Weather struct {
	ID          int32  `json:"id,omitempty"`
	Main        string `json:"main,omitempty"`
	Description string `json:"description,omitempty"`
}

type MainMeasures struct {
	Temp      float32 `json:"temp,omitempty"`
	FeelsLike float32 `json:"feels_like,omitempty"`
	TempMin   float32 `json:"temp_min,omitempty"`
	TempMax   float32 `json:"temp_max,omitempty"`
	Pressure  int32   `json:"pressure,omitempty"`
	Humidity  int32   `json:"humidity,omitempty"`
}

type Wind struct {
	Speed float32 `json:"wind,omitempty"`
	Deg   int32   `json:"deg,omitempty"`
}

type Clouds struct {
	All int32 `json:"all,omitempty"`
}

type Sys struct {
	Type    int32  `json:"type,omitempty"`
	ID      int32  `json:"id,omitempty"`
	Country string `json:"country,omitempty"`
	Sunrise int32  `json:"sunrise,omitempty"`
	Sunset  int32  `json:"sunset,omitempty"`
}
