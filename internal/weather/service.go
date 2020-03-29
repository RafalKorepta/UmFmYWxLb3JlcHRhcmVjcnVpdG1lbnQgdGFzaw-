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
	"errors"
	"net/http"
	"time"

	"github.com/rafalkorepta/UmFmYWxLb3JlcHRhcmVjcnVpdG1lbnQgdGFzaw-/internal/cache"

	"github.com/rafalkorepta/UmFmYWxLb3JlcHRhcmVjcnVpdG1lbnQgdGFzaw-/internal/api"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	maxClientTimeout = 14 * time.Second
)

type Service struct {
	external         ExternalWeatherHandler
	cache            Cache
	clientMaxTimeout time.Duration
}

type Cache interface {
	Check(cityName string) (*api.Response, error)
	Save(resp *api.Response)
}

type ExternalWeatherHandler interface {
	Weather(ctx context.Context, cityName string) (*api.Response, error)
}

func NewWeatherService(external ExternalWeatherHandler, opts ...func(*Service)) (*Service, error) {
	if external == nil {
		return nil, errors.New("cannot create service without external Weather handler")
	}

	srv := Service{
		external:         external,
		cache:            &cache.NoOpCache{},
		clientMaxTimeout: maxClientTimeout,
	}

	for _, option := range opts {
		option(&srv)
	}

	return &srv, nil
}

func WithCache(cache Cache) func(*Service) {
	return func(srv *Service) {
		srv.cache = cache
	}
}

func (s *Service) WeatherHandler(resp http.ResponseWriter, req *http.Request) {
	cities, found := req.URL.Query()["city"]
	if !found {
		resp.WriteHeader(http.StatusBadRequest)

		_, err := resp.Write([]byte("Unable to find required city query parameters"))
		if err != nil {
			zap.L().Error("Unable to write response body", zap.Error(err))
		}

		return
	}

	zap.L().Debug("Requested Cities", zap.Strings("city", cities))

	weather := make([]*api.Response, len(cities))
	g, ctx := errgroup.WithContext(req.Context())

	for iteration, city := range cities {
		i, c := iteration, city // https://golang.org/doc/faq#closures_and_goroutines

		g.Go(func() (err error) {
			r, err := s.cache.Check(c)
			if err == nil {
				weather[i] = r
				return nil
			}
			weather[i], err = s.external.Weather(ctx, c)
			s.cache.Save(weather[i])
			return err
		})
	}

	if err := g.Wait(); err != nil {
		zap.L().Error("Waiting for the error group to end the retrieval", zap.Error(err))
		resp.WriteHeader(http.StatusInternalServerError)

		return
	}

	b, err := json.Marshal(weather)
	if err != nil {
		zap.L().Error("Marshal of service response failed", zap.Error(err))
		resp.WriteHeader(http.StatusInternalServerError)

		return
	}

	resp.Header().Add("Content-Type", "application/json")

	_, err = resp.Write(b)
	if err != nil {
		zap.L().Error("Marshal of service response failed", zap.Error(err))
		resp.WriteHeader(http.StatusInternalServerError)

		return
	}
}
