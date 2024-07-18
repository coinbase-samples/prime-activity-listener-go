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
	"github.com/coinbase-samples/prime-activity-listener-go/util"
	prime "github.com/coinbase-samples/prime-sdk-go"
	"go.uber.org/zap"
)

func (l ActivityListener) handleActivities(activities []*prime.Activity) {

	zap.L().Debug("ListActivities", zap.Int("found", len(activities)))

	for _, activity := range activities {
		l.handleActivity(activity)
	}
}

func (l ActivityListener) handleActivity(activity *prime.Activity) {

	hash, err := util.Fingerprint(activity)
	if err != nil {
		zap.L().Error(
			"cannot fingerprint Activity", zap.Error(err),
			zap.String("activityId", activity.Id),
			zap.String("referenceId", activity.ReferenceId),
		)
		return
	}

	if found := l.cache.Contains(hash); found {
		return
	}

	zap.L().Debug("activity", zap.Any("doc", activity))

	if l.publishActivity(activity) {
		l.cache.Add(hash, true)
	}
}
