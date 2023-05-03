//////////////////////////////////////////////////////////////////////////
//
//  (c) TMList 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)
//  Copyright under MIT Lincese. Please see LICENSE file for details
//
//  detect_region.go - detect region (Web API entripoint URL) using
//  given API Key
//
//////////////////////////////////////////////////////////////////////////

package c1ews

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var ErrNotFound = errors.New("not found")

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

func DetectRegion(ctx context.Context, APIKey string) (result string, err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	err = ErrNotFound
	var wg sync.WaitGroup
	for _, region := range RegionList {
		wg.Add(1)
		go func(region string) {
			defer wg.Done()
			ws := NewWorkloadSecurity(APIKey, EntryPoint(region))
			if _, err := ws.DescribeCurrentAPIKey(ctx); err != nil {
				return
			}
			cancel()
			result = region
			err = nil
		}(region.ID)
	}
	wg.Wait()
	return
}

func DetectEntryPoint(ctx context.Context, APIKey string) (string, error) {
	region, err := DetectRegion(ctx, APIKey)
	if err != nil {
		return "", err
	}
	return EntryPoint(region), nil
}

func EntryPoint(region string) string {
	return fmt.Sprintf("https://workload.%s.cloudone.trendmicro.com/api", region)
}
