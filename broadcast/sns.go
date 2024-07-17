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

package broadcast

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/coinbase-samples/prime-activity-listener-go/config"
	prime "github.com/coinbase-samples/prime-sdk-go"
)

const eol = "\n"

func SnsPublishActivityMessage(
	ctx context.Context,
	app *config.AppConfig,
	activity *prime.Activity,
) error {

	v, err := json.Marshal(activity)
	if err != nil {
		return fmt.Errorf("unable to marshal activity %w", err)
	}

	input := &sns.PublishInput{
		Message:  aws.String(string(v) + eol),
		TopicArn: aws.String(app.ActivityTopicArn),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"category": {
				DataType:    aws.String("String"),
				StringValue: aws.String(activity.Category),
			},
			"primaryType": {
				DataType:    aws.String("String"),
				StringValue: aws.String(activity.PrimaryType),
			},
			"secondaryType": {
				DataType:    aws.String("String"),
				StringValue: aws.String(activity.SecondaryType),
			},
		},
	}

	if _, err := app.SnsClient.Publish(ctx, input); err != nil {
		return fmt.Errorf("unable to publish sns message %w", err)
	}

	return nil
}
