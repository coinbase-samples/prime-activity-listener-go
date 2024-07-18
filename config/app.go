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

package config

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/coinbase-samples/prime-activity-listener-go/util"
	"github.com/coinbase-samples/prime-sdk-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type AppConfig struct {
	PrimeClient                    *prime.Client
	LruCacheSizeInItems            string `mapstructure:"LRU_CACHE_SIZE"`
	ActivityPollFrequencyInSeconds string `mapstructure:"ACTIVITY_POLL_FREQUENCY"`
	ListActivitiesTimeoutInSeconds string `mapstructure:"LIST_ACTIVITIES_TIMEOUT"`
	SnsPublishTimeoutInSeconds     string `mapstructure:"SNS_PUBLISH_TIMEOUT"`
	HttpConnectTimeoutInSeconds    string `mapstructure:"HTTP_CONNECT_TIMEOUT"`
	HttpConnKeepAliveInSeconds     string `mapstructure:"HTTP_CONN_KEEP_ALIVE"`
	HttpExpectContinueInSeconds    string `mapstructure:"HTTP_EXPECT_CONTINUE"`
	HttpIdleConnInSeconds          string `mapstructure:"HTTP_IDLE_CONN"`
	HttpMaxAllIdleConnsCount       string `mapstructure:"HTTP_MAX_ALL_IDLE_CONNS"`
	HttpMaxHostIdleConnsCount      string `mapstructure:"HTTP_MAX_HOST_IDLE_CONNS"`
	HttpResponseHeaderInSeconds    string `mapstructure:"HTTP_RESPONSE_HEADER"`
	HttpTLSHandshakeInSeconds      string `mapstructure:"HTTP_TLS_HANDSHAKE"`
	EnvName                        string `mapstructure:"ENV_NAME"`
	AwsRegion                      string `mapstructure:"AWS_REGION"`
	ActivityTopicArn               string `mapstructure:"ACTIVITY_TOPIC_ARN"`
	AwsConfig                      aws.Config
	SnsClient                      *sns.Client
	HttpClient                     *http.Client
}

func (a AppConfig) IsLocalEnv() bool {
	return a.EnvName == "local"
}

func SetupAppConfig(app *AppConfig) error {

	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)

	viper.SetDefault("ENV_NAME", "local")
	viper.SetDefault("AWS_REGION", "us-east-1")
	viper.SetDefault("LRU_CACHE_SIZE", "100000")
	viper.SetDefault("ACTIVITY_POLL_FREQUENCY", "5")
	viper.SetDefault("LIST_ACTIVITIES_TIMEOUT", "10")
	viper.SetDefault("SNS_PUBLISH_TIMEOUT", "10")
	viper.SetDefault("HTTP_CONNECT_TIMEOUT", "5")
	viper.SetDefault("HTTP_CONN_KEEP_ALIVE", "30")
	viper.SetDefault("HTTP_EXPECT_CONTINUE", "1")
	viper.SetDefault("HTTP_IDLE_CONN", "90")
	viper.SetDefault("HTTP_MAX_ALL_IDLE_CONNS", "10")
	viper.SetDefault("HTTP_MAX_HOST_IDLE_CONNS", "5")
	viper.SetDefault("HTTP_RESPONSE_HEADER", "5")
	viper.SetDefault("HTTP_TLS_HANDSHAKE", "5")
	viper.SetDefault("ACTIVITY_TOPIC_ARN", "NOTSET")

	viper.ReadInConfig()

	if err := viper.Unmarshal(&app); err != nil {
		zap.L().Debug("cannot parse env file", zap.Error(err))
	}

	httpClient, err := InitHttpClient(app)
	if err != nil {
		return fmt.Errorf("cannot init the http client %w", err)
	}

	app.HttpClient = httpClient

	cfg, err := awsConfig.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(app.AwsRegion),
		config.WithHTTPClient(httpClient),
	)

	if err != nil {
		return fmt.Errorf("unable to load aws config %w", err)
	}

	app.AwsConfig = cfg

	app.SnsClient = sns.NewFromConfig(app.AwsConfig)

	return nil
}

func (a AppConfig) LruCacheSize() int {
	return convertStrIntOrFatal(a.LruCacheSizeInItems, "LruCacheSizeInItems")
}

func (a AppConfig) ActivityPollFrequency() time.Duration {
	return convertStrIntToDurationOrFatal(a.ActivityPollFrequencyInSeconds, "ActivityPollFrequencyInSeconds", time.Second)
}

func (a AppConfig) ListActivitiesTimeout() time.Duration {
	return convertStrIntToDurationOrFatal(a.ListActivitiesTimeoutInSeconds, "ListActivitiesTimeoutInSeconds", time.Second)
}

func (a AppConfig) SnsPublishTimeout() time.Duration {
	return convertStrIntToDurationOrFatal(a.SnsPublishTimeoutInSeconds, "SnsPublishTimeoutInSeconds", time.Second)
}

func (a AppConfig) HttpConnectTimeout() time.Duration {
	return convertStrIntToDurationOrFatal(a.HttpConnectTimeoutInSeconds, "HttpConnectTimeoutInSeconds", time.Second)
}

func (a AppConfig) HttpConnKeepAlive() time.Duration {
	return convertStrIntToDurationOrFatal(a.HttpConnKeepAliveInSeconds, "HttpConnKeepAliveInSeconds", time.Second)
}

func (a AppConfig) HttpExpectContinue() time.Duration {
	return convertStrIntToDurationOrFatal(a.HttpExpectContinueInSeconds, "HttpExpectContinueInSeconds", time.Second)
}

func (a AppConfig) HttpIdleConn() time.Duration {
	return convertStrIntToDurationOrFatal(a.HttpIdleConnInSeconds, "HttpIdleConnInSeconds", time.Second)
}

func (a AppConfig) HttpResponseHeader() time.Duration {
	return convertStrIntToDurationOrFatal(a.HttpResponseHeaderInSeconds, "HttpResponseHeaderInSeconds", time.Second)
}

func (a AppConfig) HttpTLSHandshake() time.Duration {
	return convertStrIntToDurationOrFatal(a.HttpTLSHandshakeInSeconds, "HttpTLSHandshakeInSeconds", time.Second)
}

func (a AppConfig) HttpMaxAllIdleConns() int {
	return convertStrIntOrFatal(a.HttpMaxAllIdleConnsCount, "HttpMaxAllIdleConnsCount")
}

func (a AppConfig) HttpMaxHostIdleConns() int {
	return convertStrIntOrFatal(a.HttpMaxHostIdleConnsCount, "HttpMaxHostIdleConnsCount")
}

func convertStrIntToDurationOrFatal(v, n string, dt time.Duration) time.Duration {
	i, err := util.ConvertStrIntToDuration(v, dt)
	if err != nil {
		zap.L().Fatal("cannot convert string to int", zap.String("value", v), zap.String("name", n), zap.Error(err))
	}
	return i
}

func convertStrIntOrFatal(v, n string) int {
	i, err := strconv.Atoi(v)
	if err != nil {
		zap.L().Fatal("cannot convert string to int", zap.String("value", v), zap.String("name", n), zap.Error(err))
	}
	return i
}
