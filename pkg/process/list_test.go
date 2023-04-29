package process

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/mpkondrashin/tmlist/pkg/c1ews"
)

func TestIncludes(t *testing.T) {
	list := &c1ews.ListResponse{
		Name:        "nameC",
		Description: "desc C",
		Items:       []string{"7", "8", "9"},
	}
	actual := Includes(list)
	expected := []string{}
	if len(expected) > 0 {
		t.Errorf("%v is not empty", actual)
	}
	list = &c1ews.ListResponse{
		Name:        "nameC",
		Description: "desc C\ninclude: nameA\nmore lines\n  Include:\t nameB   \n some test",
		Items:       []string{"7", "8", "9"},
	}
	actual = Includes(list)
	expected = []string{"nameA", "nameB"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%v is not equal to %v", actual, expected)
	}
}

func TestListCleanup(t *testing.T) {
	in := c1ews.ListResponse{
		Name:        "nameC",
		Description: "desc C",
		Items:       []string{"7", "8", "9"},
	}
	out := in
	Cleanup(&out)
	if !reflect.DeepEqual(in, out) {
		t.Errorf("%v is not equal to %v after Cleanup", in, out)
	}
	in = c1ews.ListResponse{
		Name:        "nameC",
		Description: "desc C\ninclude: nameA",
		Items:       []string{"7", "8", "9"},
	}
	out = in
	Cleanup(&out)
	if reflect.DeepEqual(in, out) {
		t.Errorf("%v is equal to %v after Cleanup", in, out)
	}
}

func TestEqual(t *testing.T) {
	a := &c1ews.ListResponse{
		ID:          1,
		Name:        "nameA",
		Description: "desc A",
		Items:       []string{"7", "8", "9"},
	}
	b := &c1ews.ListResponse{
		ID:          2,
		Name:        "nameA",
		Description: "desc A",
		Items:       []string{"7", "8", "9"},
	}
	if !Equal(a, b) {
		t.Errorf("%v and %v are not equal", a, b)
	}
	c := b
	c.ID = 3
	if !Equal(a, c) {
		t.Errorf("%v and %v are not equal", a, c)
	}
	c = b
	c.Name = "C"
	if !Equal(a, c) {
		t.Errorf("%v and %v are not equal", a, c)
	}
	c = b
	c.Description = "desc C"
	if Equal(a, c) {
		t.Errorf("%v and %v are equal", a, c)
	}
	c = b
	c.Items[0] = "changed"
	if Equal(a, c) {
		t.Errorf("%v and %v are equal", a, c)
	}
}

func TestHasIncludes(t *testing.T) {
	a := &c1ews.ListResponse{
		ID:          1,
		Name:        "nameA",
		Description: "desc A",
		Items:       []string{"7", "8", "9"},
	}
	if HasIncludes(a) {
		t.Errorf("%v should not have includes", a)
	}
	a = &c1ews.ListResponse{
		ID:          1,
		Name:        "nameA",
		Description: "desc A\ninclude: dddX",
		Items:       []string{"7", "8", "9"},
	}
	if !HasIncludes(a) {
		t.Errorf("%v does not have includes", a)
	}
}

func TestAddToTheList(t *testing.T) {
	a := &c1ews.ListResponse{
		ID:          1,
		Name:        "nameA",
		Description: "desc A\ninclude: dddX",
		Items:       []string{"5", "4", "3"},
	}
	add := []string{"1", "2", "3"}
	AddToTheList(a, add)
	actual := a.Items
	expected := []string{"1", "2", "3", "4", "5"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("[%v] is not equal to [%v] after AddToTheList", actual, expected)
	}
}

func TestClearDependence(t *testing.T) {
	{
		a := &c1ews.ListResponse{
			ID:          1,
			Name:        "nameA",
			Description: "desc A\ninclude: dddX",
			Items:       []string{"5", "4", "3"},
		}
		expected := a.Description
		ClearDependence(a)
		actual := a.Description
		if expected != actual {
			t.Errorf("[%v] is not equal to [%v] after ClearDependence", actual, expected)
		}
	}
	{
		dependence := fmt.Sprintf("\n%s eeeY fffZ", DependencePrefix)
		descriptionFmt := "desc A\ninclude: dddX%s\nother line"
		description := fmt.Sprintf(descriptionFmt, dependence)
		expected := fmt.Sprintf(descriptionFmt, "")
		a := &c1ews.ListResponse{
			ID:          1,
			Name:        "nameA",
			Description: description,
			Items:       []string{"5", "4", "3"},
		}
		ClearDependence(a)
		actual := a.Description
		if expected != actual {
			t.Errorf("[%v] is not equal to [%v] after ClearDependence", actual, expected)
		}
	}
}

func TestListDependencies(t *testing.T) {
	dependenceA := []string{"eeeY", "fffZ"}
	dependenceB := []string{"dddX"}
	description := fmt.Sprintf("desc A\ninclude: dddX\n%s %s\nSome list\n%s %s\n",
		DependencePrefix, strings.Join(dependenceA, ", "),
		DependencePrefix, strings.Join(dependenceB, ", "))
	expected := append(dependenceA, dependenceB...)
	a := &c1ews.ListResponse{
		ID:          1,
		Name:        "nameA",
		Description: description,
		Items:       []string{"5", "4", "3"},
	}
	actual := ListDependencies(a)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("[%v] is not equal to [%v] after ListDependencies", actual, expected)
	}
}

func TestAddDependence(t *testing.T) {
	{
		name := "eeeY"
		description := "desc A\ninclude: dddX"
		dependence := fmt.Sprintf("%s %s", DependencePrefix, name)
		expected := fmt.Sprintf("%s\n%s\n", description, dependence)
		a := &c1ews.ListResponse{
			ID:          1,
			Name:        "nameA",
			Description: description,
			Items:       []string{"5", "4", "3"},
		}
		AddDependences(a, name)
		actual := a.Description
		if expected != actual {
			t.Errorf("[%v] is not equal to [%v] after ClearDependence", actual, expected)
		}
	}
}

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "no duplicates",
			input:    []string{"Bob", "Charlie", "Alice"},
			expected: []string{"Alice", "Bob", "Charlie"},
		},
		{
			name:     "some duplicates",
			input:    []string{"Charlie", "Alice", "Bob", "Bob", "Charlie"},
			expected: []string{"Alice", "Bob", "Charlie"},
		},
		{
			name:     "all duplicates",
			input:    []string{"Alice", "Bob", "Bob", "Alice"},
			expected: []string{"Alice", "Bob"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := RemoveDuplicates(test.input)
			// Check that the output is correct
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("RemoveDuplicates(%v) = %v, expected %v", test.input, result, test.expected)
			}
			// Check that the output has no duplicates
			for i := 1; i < len(result); i++ {
				if result[i] == result[i-1] {
					t.Errorf("RemoveDuplicates(%v) returned duplicate elements: %v", test.input, result)
					break
				}
			}
		})
	}
}

func TestListFromResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    *c1ews.ListResponse
		expected *c1ews.List
	}{
		{
			name: "empty response",
			input: &c1ews.ListResponse{
				Name:        "",
				Description: "",
				Items:       []string{},
			},
			expected: &c1ews.List{
				Name:        "",
				Description: "",
				Items:       []string{},
			},
		},
		{
			name: "response with items",
			input: &c1ews.ListResponse{
				Name:        "MyList",
				Description: "A test list",
				Items:       []string{"item1", "item2", "item3"},
			},
			expected: &c1ews.List{
				Name:        "MyList",
				Description: "A test list",
				Items:       []string{"item1", "item2", "item3"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ListFromResponse(test.input)
			if result.Name != test.expected.Name ||
				result.Description != test.expected.Description ||
				!reflect.DeepEqual(result.Items, test.expected.Items) {
				t.Errorf("ListFromResponse(%v) = %v, expected %v", test.input, result, test.expected)
			}
		})
	}
}
