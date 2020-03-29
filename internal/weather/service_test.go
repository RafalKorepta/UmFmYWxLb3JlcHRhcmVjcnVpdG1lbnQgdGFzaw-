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

package weather_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/rafalkorepta/UmFmYWxLb3JlcHRhcmVjcnVpdG1lbnQgdGFzaw-/internal/api"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/rafalkorepta/UmFmYWxLb3JlcHRhcmVjcnVpdG1lbnQgdGFzaw-/internal/weather"
)

type mockExternalWeatherHandler struct {
	numberCalls int
	err         error
	m           sync.Mutex
}

func (m *mockExternalWeatherHandler) Weather(ctx context.Context, cityName string) (*api.Response, error) {
	m.m.Lock()
	defer m.m.Unlock()

	m.numberCalls++
	if m.err != nil {
		return nil, m.err
	}

	return &api.Response{}, nil
}

func TestService_WeatherHandler(t *testing.T) { // nolint: funlen
	t.Run("return error when missing query parameters", func(t *testing.T) {
		// Given Weather service
		s := weather.Service{}
		resp := httptest.NewRecorder()

		// And request without query params
		req := &http.Request{
			Method: "GET",
			URL: &url.URL{
				Path:     "/api/v1alpha1/weather",
				RawQuery: "",
			},
		}

		// When Weather handler is invoke
		s.WeatherHandler(resp, req)

		// Then bad request HTTP status code is returned
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.NotEmpty(t, resp.Body.String())
	})

	t.Run("return list of weather", func(t *testing.T) {
		tcs := []struct {
			name          string
			queryParam    string
			numberOfCalls int
		}{
			{
				name:          "one city",
				queryParam:    "city=Some",
				numberOfCalls: 1, // nolint: gomnd
			},
			{
				name:          "many cities",
				queryParam:    "city=Some&city=Other&city=NewCity",
				numberOfCalls: 3, // nolint: gomnd
			},
		}
		for _, tc := range tcs {
			qp := tc.queryParam
			numCalls := tc.numberOfCalls
			t.Run(tc.name, func(t *testing.T) {
				// Given Weather service
				mock := mockExternalWeatherHandler{}
				s, err := weather.NewWeatherService(&mock)
				require.NoError(t, err)

				resp := httptest.NewRecorder()

				// And request with query params
				req := &http.Request{
					Method: "GET",
					URL: &url.URL{
						Path:     "/api/v1alpha1/weather",
						RawQuery: qp,
					},
				}

				// When Weather handler is invoke
				s.WeatherHandler(resp, req)

				// Then all works correctly
				assert.Equal(t, http.StatusOK, resp.Code)
				assert.NotEmpty(t, resp.Body.String())
				assert.Equal(t, numCalls, mock.numberCalls)
			})
		}
	})

	t.Run("Error when external handler fail", func(t *testing.T) {
		// Given Weather service
		mock := mockExternalWeatherHandler{
			err: errors.New("some error"),
		}
		s, err := weather.NewWeatherService(&mock)
		require.NoError(t, err)

		resp := httptest.NewRecorder()

		// And request with query params
		req := &http.Request{
			Method: "GET",
			URL: &url.URL{
				Path:     "/api/v1alpha1/weather",
				RawQuery: "city=Some",
			},
		}

		// When Weather handler is invoke
		s.WeatherHandler(resp, req)

		// Then bad request HTTP status code is returned
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Empty(t, resp.Body.String())
	})
}

func TestNewWeatherService(t *testing.T) {
	t.Run("Missing required external weather handler", func(t *testing.T) {
		_, err := weather.NewWeatherService(nil)
		assert.Error(t, err)
	})
}
