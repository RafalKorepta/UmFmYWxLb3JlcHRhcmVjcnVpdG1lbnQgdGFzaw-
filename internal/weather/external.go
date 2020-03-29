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

package weather

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/rafalkorepta/UmFmYWxLb3JlcHRhcmVjcnVpdG1lbnQgdGFzaw-/internal/api"

	"github.com/pkg/errors"
)

type ExternalService struct {
	ExternalWeatherHandler
	netClient *http.Client
	apiKey    string
}

const (
	openWeatherMapAPI  = "http://api.openweathermap.org/data/2.5/weather?q="
	apiParameterPrefix = "&apiKey="

	transportDialTimeout = 5 * time.Second
)

func NewExternalService(apiKey string) (*ExternalService, error) {
	if apiKey == "" {
		return nil, errors.New("cannot create external weather service without apiKey")
	}

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: transportDialTimeout,
		}).DialContext,
	}

	client := &http.Client{Transport: transport}

	return &ExternalService{
		netClient: client,
		apiKey:    apiParameterPrefix + apiKey,
	}, nil
}

func (s *ExternalService) Weather(ctx context.Context, cityName string) (*api.Response, error) {
	r := &api.Response{}

	weatherURL := openWeatherMapAPI + cityName + s.apiKey

	req, err := http.NewRequestWithContext(ctx, "GET", weatherURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create request")
	}

	resp, err := s.netClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "unable to send the request")
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read to body")
	}

	err = json.Unmarshal(b, &r)
	if err != nil {
		return nil, errors.Wrap(err, "unable to marshal weather response")
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, errors.Wrap(err, "unable to close the body")
	}

	return r, nil
}
