//////////////////////////////////////////////////////////////////////////
//
//  (c) TMList 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)
//  Copyright under MIT Lincese. Please see LICENSE file for details
//
//  list.go - various functions to manipulate exclusion lists
//
//////////////////////////////////////////////////////////////////////////

package process

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/mpkondrashin/tmlist/pkg/c1ews"
	"golang.org/x/exp/maps"
)

// Equal - return true if there is a difference between lists.
// We are checking only fileds that could be changed by TMList
func Equal(a, b *c1ews.ListResponse) bool {
	return a.Description == b.Description && reflect.DeepEqual(a.Items, b.Items)
}

// Includes - list all includes of the list
func Includes(l *c1ews.ListResponse) (result []string) {
	for _, line := range strings.Split(l.Description, "\n") {
		colon := strings.Index(line, ":")
		if colon == -1 {
			continue
		}
		include := line[:colon]
		if strings.ToLower(strings.TrimSpace(include)) != "include" {
			continue
		}
		name := strings.TrimSpace(line[colon+1:])
		result = append(result, name)
	}
	return
}

// HasIncludes - return true if list has includes
func HasIncludes(l *c1ews.ListResponse) bool {
	return len(Includes(l)) > 0
}

// Cleanup - remove all items if list has includes
func Cleanup(l *c1ews.ListResponse) {
	if HasIncludes(l) {
		l.Items = []string{}
	}
}

// AddToTheList - add items to the list avoiding duplicates and sort them
func AddToTheList(l *c1ews.ListResponse, items []string) {
	l.Items = RemoveDuplicates(append(l.Items, items...))
}

// ClearDependence - remove dependence lines from description if exist any
func ClearDependence(l *c1ews.ListResponse) {
	result := []string{}
	for _, line := range strings.Split(l.Description, "\n") {
		if strings.HasPrefix(line, DependencePrefix) {
			continue
		}
		result = append(result, line)
	}
	l.Description = strings.Join(result, "\n")
}

// ListDependencies - return all dependencies for given list
func ListDependencies(l *c1ews.ListResponse) (result []string) {
	for _, line := range strings.Split(l.Description, "\n") {
		if !strings.HasPrefix(line, DependencePrefix) {
			continue
		}
		line = line[len(DependencePrefix):]
		for _, name := range strings.Split(line, ",") {
			name = strings.TrimSpace(name)
			result = append(result, name)
		}
	}
	return
}

// AddDependences - add dependence line or modify one of it already exists
func AddDependences(l *c1ews.ListResponse, name ...string) {
	if len(name) == 0 {
		return
	}
	//fmt.Println("XXX AddDependences", strings.Join(name, "|"))
	deps := ListDependencies(l)
	ClearDependence(l)
	deps = RemoveDuplicates(append(deps, name...))
	dependence := fmt.Sprintf("\n%s %s", DependencePrefix, strings.Join(deps, ", "))
	l.Description = l.Description + dependence
}

// RemoveDuplicates - remove duplicates from string slice. Return in sorted order
func RemoveDuplicates(names []string) (result []string) {
	m := make(map[string]struct{})
	for _, each := range names {
		m[each] = struct{}{}
	}
	result = maps.Keys(m)
	sort.Strings(result)
	return
}

// ListFromResponse - return List struct build from ListResponse struct
func ListFromResponse(response *c1ews.ListResponse) *c1ews.List {
	return &c1ews.List{
		Name:        response.Name,
		Description: response.Description,
		Items:       response.Items,
	}
}
