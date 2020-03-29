package cache

import (
	"github.com/pkg/errors"

	"github.com/rafalkorepta/UmFmYWxLb3JlcHRhcmVjcnVpdG1lbnQgdGFzaw-/internal/api"
)

type NoOpCache struct{}

func (in *NoOpCache) Check(cityName string) (*api.Response, error) {
	return nil, errors.New("city does not exist")
}

func (in *NoOpCache) Save(resp *api.Response) {}
