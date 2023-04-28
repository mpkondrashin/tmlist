package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mpkondrashin/tmlist/pkg/c1ews"
	"github.com/mpkondrashin/tmlist/pkg/process"
)

/*
Get List of Lists
for each list:
   for each include:


*/

func main() {
	host := "https://workload.trend-us-1.cloudone.trendmicro.com/api"
	apikey := "tmc12OuKXO2Ji71RFrjWYD2d9KPj2GW:7wQFMF5rGHdyE5gAoDLyXZE3nEWAt6TrnDnwF88YzgJzd28YxwmzXiPoqdUTA7Rcjy"
	ws := c1ews.NewWorkloadSecurity(apikey, host)
	type List func(context.Context) ([]c1ews.ListResponse, error)
	type Modify func(context.Context, int, *c1ews.List) (*c1ews.ListResponse, error)

	queries := []struct {
		name   string
		list   List
		modify Modify
	}{
		{"directory list", ws.ListDirectoryLists, ws.ModifyDirectoryList},
		{"file extension list", ws.ListFileExtensionLists, ws.ModifyFileExtensionList},
		{"file list", ws.ListFileLists, ws.ModifyFileList},
	}
	for _, query := range queries {
		log.Printf("Query: %s", query.name)
		r, err := query.list(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		p := process.NewProcess(r)
		err = p.Process()
		if err != nil {
			log.Fatal(fmt.Errorf("%s: %w", query.name, err))
		}
		p.IterateChanged(func(list *c1ews.ListResponse) error {
			log.Printf("%s: modify %s", query.name, list.Name)
			l := process.ListFromResponse(list)
			_, err := query.modify(context.TODO(), list.ID, l)
			return err
		})
	}

}
