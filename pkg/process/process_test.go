//////////////////////////////////////////////////////////////////////////
//
//  (c) TMList 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)
//  Copyright under MIT Lincese. Please see LICENSE file for details
//
//  process_test.go - test functions in process.go
//
//////////////////////////////////////////////////////////////////////////

package process

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/mpkondrashin/tmlist/pkg/c1ews"
)

func TestNoIncludesIsChanged(t *testing.T) {
	in := []c1ews.ListResponse{
		{Name: "nameA",
			Description: "desc A",
			Items:       []string{"1", "2", "3"},
		},
		{Name: "nameB",
			Description: "desc B",
			Items:       []string{"4", "5", "6"},
		},
		{Name: "nameC",
			Description: "desc C",
			Items:       []string{"7", "8", "9"},
		},
	}
	p := NewProcess(in)
	p.Process()
	_ = p.IterateChanged(func(l *c1ews.ListResponse) error {
		t.Errorf("One list without includes was changed: %s", l.Name)
		return nil
	})
}

func TestProcess(t *testing.T) {
	in := []c1ews.ListResponse{
		{Name: "nameA",
			Description: "desc A\n",
			Items:       []string{"3", "2", "1"},
		},
		{Name: "nameB",
			Description: "desc B\ninclude: nameA",
			Items:       []string{"4", "5", "6"},
		},
		{Name: "nameC",
			Description: "desc B\ninclude: nameB",
			Items:       []string{"7", "8", "9"},
		},
		{Name: "nameD",
			Description: "desc B\ninclude: nameA\ninclude: nameC",
			Items:       []string{"7", "8", "9"},
		},
	}
	p := NewProcess(in)
	p.Process()
	count := 0
	_ = p.IterateChanged(func(l *c1ews.ListResponse) error {
		//t.Logf("Changed %s (%d)", l.Name, l.ID)
		//t.Logf("Description %s", l.Description)
		//t.Logf("Items %s", l.Items)
		count++
		switch l.Name {
		case "nameA":
			{
				expect := in[1].Name
				if !strings.Contains(l.Description, expect) {
					t.Errorf("%s not found in %s description", expect, l.Name)
				}
			}
			{
				expect := in[2].Name
				if !strings.Contains(l.Description, expect) {
					t.Errorf("%s not found in %s description", expect, l.Name)
				}
			}
			{
				expect := in[3].Name
				if !strings.Contains(l.Description, expect) {
					t.Errorf("%s not found in %s description", expect, l.Name)
				}
			}
		case "nameB":
			{
				expect := in[2].Name
				if !strings.Contains(l.Description, expect) {
					t.Errorf("%s not found in %s description", expect, l.Name)
				}
			}
			{
				expect := []string{"1", "2", "3"}
				if !reflect.DeepEqual(l.Items, expect) {
					t.Errorf("Items are [%v] and not [%v]", l.Items, expect)
				}
			}
		case "nameC":
			{
				expect := in[3].Name
				if !strings.Contains(l.Description, expect) {
					t.Errorf("%s not found in %s description", expect, l.Name)
				}
			}
			{
				expect := []string{"1", "2", "3"}
				if !reflect.DeepEqual(l.Items, expect) {
					t.Errorf("Items are [%v] and not [%v]", l.Items, expect)
				}
			}
		case "nameD":
			{
				expect := []string{"1", "2", "3"}
				if !reflect.DeepEqual(l.Items, expect) {
					t.Errorf("Items are [%v] and not [%v]", l.Items, expect)
				}
			}
		}
		return nil
	})
	expected := len(in)
	if count != expected {
		t.Errorf("Changed wrong number of lists: %d (expected %d)", count, expected)
	}
}

func TestLoop(t *testing.T) {
	tests := []struct {
		name  string
		input []c1ews.ListResponse
	}{
		{"A->A",
			[]c1ews.ListResponse{
				{Name: "nameA",
					Description: "include: nameA",
					Items:       []string{"1", "2", "3"},
				},
			}},
		{"A->B,B->A",
			[]c1ews.ListResponse{
				{Name: "nameA",
					Description: "include: nameB",
					Items:       []string{"1", "2", "3"},
				},
				{Name: "nameB",
					Description: "Include: nameA",
					Items:       []string{"4", "5", "6"},
				},
			}},
		{"A->B,B->C,C->A",
			[]c1ews.ListResponse{
				{Name: "nameA",
					Description: "include: nameB",
					Items:       []string{"1", "2", "3"},
				},
				{Name: "nameB",
					Description: "Include: nameC",
					Items:       []string{"4", "5", "6"},
				},
				{Name: "nameC",
					Description: "Include: nameA",
					Items:       []string{"7", "8", "9"},
				},
			}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := NewProcess(test.input)
			err := p.Process()
			if !errors.Is(err, ErrCycleDependence) {
				t.Errorf("Cycle %s not detected", test.name)
			}
		})
	}
}
