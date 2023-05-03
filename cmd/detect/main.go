//////////////////////////////////////////////////////////////////////////
//
//  (c) TMList 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)
//  Copyright under MIT Lincese. Please see LICENSE file for details
//
//  main.go - detection of API entry point URL using API Key
//
//////////////////////////////////////////////////////////////////////////

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mpkondrashin/tmlist/pkg/c1ews"
	"github.com/mpkondrashin/tmlist/pkg/cone"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	RCOther    = 0
	RCNotFound = 1
)

const EnvPrefix = "DETECT"

const (
	ConfigFileName = "config"
	ConfigFileType = "yaml"
)

const (
	flagAPIKey          = "api_key"
	flagIgnoreTLSErrors = "ignore_tls_errors"
)

func Configure() {
	fs := pflag.NewFlagSet("", pflag.ExitOnError)
	fs.String(flagAPIKey, "", "Cloud One API Key")
	fs.Bool(flagIgnoreTLSErrors, false, "Ignore all TLS errors")
	err := fs.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	if err := viper.BindPFlags(fs); err != nil {
		log.Fatal(err)
	}
	viper.SetEnvPrefix(EnvPrefix)
	viper.AutomaticEnv()

	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType(ConfigFileType)
	path, err := os.Executable()
	if err == nil {
		dir := filepath.Dir(path)
		viper.AddConfigPath(dir)
	}
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if !ok {
			log.Fatal(err)
		}
	}
}

func DetectEntryPoint(apikey string) (string, error) {
	region, err := cone.NewClient(apikey).CurrentAPIKeyRegion(context.TODO())
	if err == nil {
		return c1ews.EntryPoint(region), nil
	}
	return c1ews.DetectEntryPoint(context.TODO(), apikey)
}

func main() {
	Configure()
	apikey := viper.GetString(flagAPIKey)
	if apikey == "" {
		log.Fatal(fmt.Errorf("%s parameter is missing", flagAPIKey))
	}
	entryPoint, err := DetectEntryPoint(apikey)
	if err != nil {
		log.Fatal("Not detected")
	}
	fmt.Println(entryPoint)
}
