/*
Copyright 2021 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package redis

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"

	"github.com/dapr/components-contrib/bindings"
	rediscomponent "github.com/dapr/components-contrib/internal/component/redis"
	"github.com/dapr/kit/logger"
)

// Redis is a redis output binding.
type Redis struct {
	client         redis.UniversalClient
	clientSettings *rediscomponent.Settings
	logger         logger.Logger

	ctx    context.Context
	cancel context.CancelFunc
}

// NewRedis returns a new redis bindings instance.
func NewRedis(logger logger.Logger) *Redis {
	return &Redis{logger: logger}
}

// Init performs metadata parsing and connection creation.
func (r *Redis) Init(meta bindings.Metadata) (err error) {
	r.client, r.clientSettings, err = rediscomponent.ParseClientFromProperties(meta.Properties, nil)
	if err != nil {
		return err
	}

	r.ctx, r.cancel = context.WithCancel(context.Background())

	_, err = r.client.Ping(r.ctx).Result()
	if err != nil {
		return fmt.Errorf("redis binding: error connecting to redis at %s: %s", r.clientSettings.Host, err)
	}

	return err
}

func (r *Redis) Operations() []bindings.OperationKind {
	return []bindings.OperationKind{bindings.CreateOperation}
}

func (r *Redis) Invoke(ctx context.Context, req *bindings.InvokeRequest) (*bindings.InvokeResponse, error) {
	if val, ok := req.Metadata["key"]; ok && val != "" {
		key := val
		_, err := r.client.Do(ctx, "SET", key, req.Data).Result()
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	return nil, errors.New("redis binding: missing key on write request metadata")
}

func (r *Redis) Close() error {
	r.cancel()

	return r.client.Close()
}
