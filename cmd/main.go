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

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/coinbase-samples/prime-activity-listener-go/config"
	"github.com/coinbase-samples/prime-activity-listener-go/listen"
	prime "github.com/coinbase-samples/prime-sdk-go"
	"go.uber.org/zap"
)

func main() {

	run := make(chan os.Signal, 1)
	signal.Notify(run, os.Interrupt, syscall.SIGTERM)

	log := config.LogInit("prime-activities")
	zap.ReplaceGlobals(log)
	defer log.Sync()

	log.Info("prime-activities", zap.String("state", "starting"))

	if err := os.Setenv("TZ", "UTC"); err != nil {
		log.Fatal("cannot set time zone: UTC", zap.Error(err))
	}

	credentials, err := prime.ReadEnvCredentials("PRIME_CREDENTIALS")
	if err != nil {
		log.Fatal("cannot init the prime credentials", zap.Error(err))
	}

	appConfig := &config.AppConfig{}

	if err := config.SetupAppConfig(appConfig); err != nil {
		log.Fatal("cannot setup app config", zap.Error(err))
	}

	appConfig.PrimeClient = prime.NewClient(credentials, *appConfig.HttpClient)

	daemon, err := listen.StartActivityListener(appConfig)
	if err != nil {
		log.Fatal("cannot start the activity listener", zap.Error(err))
	}

	log.Info("prime-activities", zap.String("state", "started"))

	<-run

	log.Info("prime-activities", zap.String("state", "stopping"))

	if err := listen.StopActivityListener(daemon); err != nil {
		log.Error("process did not stop cleanly", zap.Error(err))
	}

	log.Info("prime-activities", zap.String("state", "stopped"))
}
