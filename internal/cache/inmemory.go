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

package cache

import (
	"errors"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/rafalkorepta/UmFmYWxLb3JlcHRhcmVjcnVpdG1lbnQgdGFzaw-/internal/api"
)

const (
	defaultTTL = 5 * time.Minute
)

type InMemory struct {
	ttl   time.Duration
	cache *cache.Cache
}

func NewInMemoryCache() *InMemory {
	return &InMemory{
		cache: cache.New(defaultTTL, time.Minute+defaultTTL),
		ttl:   defaultTTL,
	}
}

func (in *InMemory) Check(cityName string) (*api.Response, error) {
	item, exist := in.cache.Get(cityName)
	if !exist {
		return nil, errors.New("city does not exist")
	}

	return item.(*api.Response), nil
}

func (in *InMemory) Save(resp *api.Response) {
	in.cache.Set(resp.Name, resp, cache.DefaultExpiration)
}
