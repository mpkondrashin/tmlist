package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mpkondrashin/tmlist/pkg/c1ews"
	"github.com/mpkondrashin/tmlist/pkg/process"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const EnvPrefix = "TMLIST"

const (
	ConfigFileName = "config"
	ConfigFileType = "yaml"
)

const (
	flagAddress = "address"
	flagAPIKey  = "api_key"
)

func Configure() {
	fs := pflag.NewFlagSet("", pflag.ExitOnError)
	fs.String(flagAddress, "", "Cloud One Woekload Security entry point URL")
	fs.String(flagAPIKey, "", "Cloud One API Key")
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

func main() {
	Configure()
	host := viper.GetString(flagAddress) // "https://workload.trend-us-1.cloudone.trendmicro.com/api"
	if host == "" {
		log.Fatal(fmt.Errorf("%s parameter is missing", flagAddress))
	}
	apikey := viper.GetString(flagAPIKey) // "tmc12OuKXO2Ji71RFrjWYD2d9KPj2GW:7wQFMF5rGHdyE5gAoDLyXZE3nEWAt6TrnDnwF88YzgJzd28YxwmzXiPoqdUTA7Rcjy"
	if apikey == "" {
		log.Fatal(fmt.Errorf("%s parameter is missing", flagAPIKey))
	}
	ws := c1ews.NewWorkloadSecurity(apikey, host)

	type List func(context.Context) ([]c1ews.ListResponse, error)
	type Modify func(context.Context, int, *c1ews.List) (*c1ews.ListResponse, error)

	queries := []struct {
		name   string
		list   List
		modify Modify
	}{
		{"file extension list", ws.ListFileExtensionLists, ws.ModifyFileExtensionList},
		{"file list", ws.ListFileLists, ws.ModifyFileList},
		{"directory list", ws.ListDirectoryLists, ws.ModifyDirectoryList},
	}
	for _, query := range queries {
		log.Printf("%s: Start", query.name)
		r, err := query.list(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		p := process.NewProcess(r)
		err = p.Process()
		if err != nil {
			log.Fatal(fmt.Errorf("%s: %w", query.name, err))
		}
		count := 0
		p.IterateChanged(func(list *c1ews.ListResponse) error {
			log.Printf("%s: modify %s", query.name, list.Name)
			l := process.ListFromResponse(list)
			_, err := query.modify(context.TODO(), list.ID, l)
			return err
		})
		if count == 0 {
			log.Printf("%s: No modifications", query.name)
		}
	}

}
