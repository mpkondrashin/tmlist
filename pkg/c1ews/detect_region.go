package c1ews

import (
	"context"
	"fmt"
	"sync"
)

type Region struct {
	Name string
	ID   string
}

// Taken from https://cloudone.trendmicro.com/docs/identity-and-account-management/c1-regions/
var RegionList = []Region{
	{"US", "us-1"},
	{"India", "in-1"},
	{"UK", "gb-1"},
	{"Japan", "jp-1"},
	{"Germany", "de-1"},
	{"Australia", "au-1"},
	{"Canada", "ca-1"},
	{"Singapore", "sg-1"},
	{"Trend US", "trend-us-1"},
}

func DetectEntryPoint(ctx context.Context, APIKey string) (result string) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var wg sync.WaitGroup
	for _, region := range RegionList {
		host := fmt.Sprintf("https://workload.%s.cloudone.trendmicro.com/api", region.ID)
		wg.Add(1)
		go func(h string) {
			defer wg.Done()
			ws := NewWorkloadSecurity(APIKey, h)
			_, err := ws.DescribeCurrentAPIKey(ctx)
			if err != nil {
				return
			}
			cancel()
			result = h
		}(host)
	}
	wg.Wait()
	return
}