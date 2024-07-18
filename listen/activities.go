/**
 * Copyright 2024-present Coinbase Global, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package listen

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/coinbase-samples/prime-activity-listener-go/config"
	prime "github.com/coinbase-samples/prime-sdk-go"
	lru "github.com/hashicorp/golang-lru/v2"
	"go.uber.org/zap"
)

type ActivityListener struct {
	config        *config.AppConfig
	cache         *lru.Cache[string, bool]
	stopWaitGroup sync.WaitGroup
	running       atomic.Bool
}

func (l ActivityListener) run() {

	defer l.stopWaitGroup.Done()

	for {

		if !l.running.Load() {
			break
		}

		l.poll()

		time.Sleep(l.config.ActivityPollFrequency())
	}

}

func (l ActivityListener) poll() bool {

	var nextCursor string

	for {

		response, err := l.listActivities(nextCursor)
		if err != nil {
			zap.L().Error("unable to call ListActivities", zap.Error(err))
			return false
		}

		l.handleActivities(response.Activities)

		if response.Pagination != nil {
			if !response.Pagination.HasNext {
				break
			}

			nextCursor = response.Pagination.NextCursor
		}
	}

	return true
}

func (l ActivityListener) listActivities(nextCursor string) (*prime.ListActivitiesResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), l.config.ListActivitiesTimeout())
	defer cancel()

	now := time.Now().UTC()

	pagination := &prime.PaginationParams{Limit: "100"}

	if len(nextCursor) > 0 {
		pagination.Cursor = nextCursor
	}

	response, err := l.config.PrimeClient.ListActivities(
		ctx,
		&prime.ListActivitiesRequest{
			PortfolioId: l.config.PrimeClient.Credentials.PortfolioId,
			Start:       now.Add(l.config.ActivityPollFrequency() * time.Duration(-400)),
			End:         now,
			Pagination:  pagination,
		},
	)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func StartActivityListener(config *config.AppConfig) (*ActivityListener, error) {

	cache, err := lru.New[string, bool](config.LruCacheSize())

	if err != nil {
		return nil, err
	}

	listener := &ActivityListener{
		config: config,
		cache:  cache,
	}

	listener.stopWaitGroup.Add(1)

	listener.running.Store(true)

	go listener.run()

	return listener, nil
}

func StopActivityListener(al *ActivityListener) error {

	al.running.Store(false)

	al.stopWaitGroup.Wait()

	return nil
}
