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
	flagDir     = "dir"
	flagExt     = "ext"
	flagFile    = "file"
	flagDryRun  = "dry"
)

func Configure() {
	fs := pflag.NewFlagSet("", pflag.ExitOnError)
	fs.String(flagAddress, "", "Cloud One Woekload Security entry point URL")
	fs.String(flagAPIKey, "", "Cloud One API Key")
	fs.Bool(flagDir, false, "Process directory lists")
	fs.Bool(flagExt, false, "Process file extension lists")
	fs.Bool(flagFile, false, "Process file lists")
	fs.Bool(flagDryRun, false, "Dyr run - do not modify existing lists")
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

type (
	List   func(context.Context) ([]c1ews.ListResponse, error)
	Modify func(context.Context, int, *c1ews.List) (*c1ews.ListResponse, error)
)

func ProcessList(name string, list List, modify Modify, dryRun bool) {
	log.Printf("%s: Start", name)
	r, err := list(context.TODO())
	if err != nil {
		log.Print(err)
		return
	}
	p := process.NewProcess(r)
	err = p.Process()
	if err != nil {
		log.Print(fmt.Errorf("%s: %w", name, err))
		return
	}
	count := 0
	p.IterateChanged(func(list *c1ews.ListResponse) error {
		count++
		log.Printf("%s: modify %s", name, list.Name)
		if dryRun {
			return nil
		}
		l := process.ListFromResponse(list)
		_, err := modify(context.TODO(), list.ID, l)
		return err
	})
	if count == 0 {
		log.Printf("%s: No modifications", name)
	}
}

func main() {
	Configure()
	host := viper.GetString(flagAddress)
	if host == "" {
		log.Fatal(fmt.Errorf("%s parameter is missing", flagAddress))
	}
	apikey := viper.GetString(flagAPIKey)
	if apikey == "" {
		log.Fatal(fmt.Errorf("%s parameter is missing", flagAPIKey))
	}
	ws := c1ews.NewWorkloadSecurity(apikey, host)
	dryRun := viper.GetBool(flagDryRun)
	all := !viper.GetBool(flagDir) && !viper.GetBool(flagExt) && !viper.GetBool(flagFile)
	if viper.GetBool(flagDir) || all {
		ProcessList("Directory Lists", ws.ListDirectoryLists, ws.ModifyDirectoryList, dryRun)
	}
	if viper.GetBool(flagExt) || all {
		ProcessList("File Extension Lists", ws.ListFileExtensionLists, ws.ModifyFileExtensionList, dryRun)
	}
	if viper.GetBool(flagFile) || all {
		ProcessList("File Lists", ws.ListFileLists, ws.ModifyFileList, dryRun)
	}
}
