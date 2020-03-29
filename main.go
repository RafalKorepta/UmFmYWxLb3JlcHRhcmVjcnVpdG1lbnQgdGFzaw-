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

package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/rafalkorepta/UmFmYWxLb3JlcHRhcmVjcnVpdG1lbnQgdGFzaw-/internal/cache"
	"github.com/rafalkorepta/UmFmYWxLb3JlcHRhcmVjcnVpdG1lbnQgdGFzaw-/internal/weather"
)

var (
	// Version will be populated with binary semver by the linker
	// during the build process.
	// See https://blog.cloudflare.com/setting-go-variables-at-compile-time/
	// and https://golang.org/cmd/link/ in section Flags `-X importpath.name=value`.
	Version string //nolint: gochecknoglobals

	// Commit will be populated with correct git commit id by the linker
	// during the build process.
	// See https://blog.cloudflare.com/setting-go-variables-at-compile-time/
	// and https://golang.org/cmd/link/ in section Flags `-X importpath.name=value`.
	Commit string //nolint: gochecknoglobals
)

const (
	debugFlag          = "debug"
	configPathFlag     = "configPath"
	configFileNameFlag = "config"
	apiKeyFlag         = "apiKey"
	portFlag           = "port"

	readTimeout  = 15 * time.Second
	writeTimeout = 15 * time.Second
)

func main() {
	err := initConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}

	ext, err := weather.NewExternalService(viper.GetString(apiKeyFlag))
	if err != nil {
		zap.L().Fatal("Unable to create external weather service", zap.Error(err))
	}

	c := cache.NewInMemoryCache()

	weatherService, err := weather.NewWeatherService(ext, weather.WithCache(c))
	if err != nil {
		zap.L().Fatal("Unable to create service", zap.Error(err))
	}

	r := mux.NewRouter()
	r.Schemes("http").
		Methods("GET").
		Path("/api/v1alpha1/weather").
		HandlerFunc(weatherService.WeatherHandler)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:" + strconv.Itoa(viper.GetInt(portFlag)),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
	}

	err = srv.ListenAndServe()
	if err != nil {
		zap.L().Fatal("Server failed", zap.Error(err))
	}
}

func initConfig() error {
	pflag.Int(portFlag, 8000, "Number of the port that service will listen on")
	pflag.String(configFileNameFlag, "config.yaml", "Name of the config file")
	pflag.String(configPathFlag, ".", "Relative path where config resides")
	pflag.String(apiKeyFlag, "", "ApiKey for open weather map")
	pflag.Bool(debugFlag, false, "setup logger for debug log level and prettify logs")
	pflag.Parse()

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return errors.Wrap(err, "binding pflags to viper")
	}

	viper.SetConfigName(viper.GetString(configFileNameFlag)) // name of config file (without extension)
	viper.AddConfigPath(viper.GetString(configPathFlag))
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err = viper.ReadInConfig(); err != nil {
		zap.S().Errorw("Failed to read from config file",
			"configFile", viper.ConfigFileUsed(),
			"error", err)
	}

	var cfg zap.Config
	if viper.GetBool(debugFlag) {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	newLogger, err := cfg.Build(zap.AddStacktrace(zap.ErrorLevel),
		zap.Fields(
			zap.Field{
				Key:    "commit",
				Type:   zapcore.StringType,
				String: Commit,
			},
			zap.Field{
				Key:    "version",
				Type:   zapcore.StringType,
				String: Version,
			},
		))
	if err != nil {
		log.Fatalf("Unable to create logger. Error: %v", err)
	}

	zap.ReplaceGlobals(newLogger)

	return nil
}
