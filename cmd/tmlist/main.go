//////////////////////////////////////////////////////////////////////////
//
//  (c) TMList 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)
//  Copyright under MIT Lincese. Please see LICENSE file for details
//
//  main.go - main LMList file
//
//////////////////////////////////////////////////////////////////////////

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mpkondrashin/tmlist/pkg/c1ews"
	"github.com/mpkondrashin/tmlist/pkg/process"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	RCOther = 3 + iota
	RCAPIError
	RCCycleDependence
	RCListNotFound
)

const EnvPrefix = "TMLIST"

const (
	ConfigFileName = "config"
	ConfigFileType = "yaml"
)

const (
	flagAddress         = "address"
	flagAPIKey          = "api_key"
	flagIgnoreTLSErrors = "ignore_tls_errors"
	flagDir             = "dir"
	flagExt             = "ext"
	flagFile            = "file"
	flagDryRun          = "dry"
)

func Configure() {
	fs := pflag.NewFlagSet("", pflag.ExitOnError)
	fs.String(flagAddress, "", "Cloud One Woekload Security entry point URL")
	fs.String(flagAPIKey, "", "Cloud One API Key")
	fs.Bool(flagIgnoreTLSErrors, false, "Ignore all TLS errors")
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

func ProcessList(name string, list List, modify Modify, dryRun bool) int {
	log.Printf("%s: Start", name)
	r, err := list(context.TODO())
	if err != nil {
		log.Print(err)
		return RCAPIError
	}
	p := process.NewProcess(r)
	err = p.Process()
	if err != nil {
		log.Printf("%s: %v", name, err)
		if errors.Is(err, process.ErrListNotFound) {
			return RCListNotFound
		}
		if errors.Is(err, process.ErrCycleDependence) {
			return RCCycleDependence
		}
		return RCOther
	}
	count := 0
	err = p.IterateChanged(func(list *c1ews.ListResponse) error {
		count++
		log.Printf("%s: modify %s", name, list.Name)
		if dryRun {
			return nil
		}
		l := process.ListFromResponse(list)
		_, err := modify(context.TODO(), list.ID, l)
		return err
	})
	if err != nil {
		return RCAPIError
	}
	if count == 0 {
		log.Printf("%s: No modifications", name)
	}
	return 0
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
	ws.SetIgnoreTLSErrors(viper.GetBool(flagIgnoreTLSErrors))
	dryRun := viper.GetBool(flagDryRun)
	all := !viper.GetBool(flagDir) && !viper.GetBool(flagExt) && !viper.GetBool(flagFile)
	returnCode := 0
	if viper.GetBool(flagDir) || all {
		rc := ProcessList("Directory Lists", ws.ListDirectoryLists, ws.ModifyDirectoryList, dryRun)
		if rc > returnCode {
			returnCode = rc
		}
	}
	if viper.GetBool(flagExt) || all {
		rc := ProcessList("File Extension Lists", ws.ListFileExtensionLists, ws.ModifyFileExtensionList, dryRun)
		if rc > returnCode {
			returnCode = rc
		}
	}
	if viper.GetBool(flagFile) || all {
		rc := ProcessList("File Lists", ws.ListFileLists, ws.ModifyFileList, dryRun)
		if rc > returnCode {
			returnCode = rc
		}
	}
	os.Exit(returnCode)
}
