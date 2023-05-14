//////////////////////////////////////////////////////////////////////////
//
//  (c) TMList 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)
//  Copyright under MIT Lincese. Please see LICENSE file for details
//
//  levenshtein_test.go - levenshtein distance test
//
//////////////////////////////////////////////////////////////////////////

package levenshtein

import (
	"fmt"
	"testing"
)

func TestSame(t *testing.T) {
	testCases := []struct {
		a string
		b string
	}{
		{"a", "a"},
		{"ab", "ab"},
		{"abc", "abc"},
		{"12344321", "12344321"},
	}
	for _, each := range testCases {
		t.Run(fmt.Sprintf("%s_vs_%s", each.a, each.b), func(t *testing.T) {
			expected := 0
			actual := Distance(each.a, each.b)
			if actual != expected {
				t.Errorf("a = \"%s\" b = \"%s\", but %d != %d", each.a, each.b, actual, expected)
			}

		})
	}
}

func TestWithEmpty(t *testing.T) {
	testCases := []struct {
		a string
		b string
	}{
		{"a", ""},
		{"", "ab"},
		{"abc", ""},
		{"", "12344321"},
	}
	for _, each := range testCases {
		t.Run(fmt.Sprintf("%s_vs_%s", each.a, each.b), func(t *testing.T) {
			expected := len(each.a) + len(each.b)
			actual := Distance(each.a, each.b)
			if actual != expected {
				t.Errorf("a = \"%s\" b = \"%s\", but %d != %d", each.a, each.b, actual, expected)
			}

		})
	}
}

func TestOne(t *testing.T) {
	testCases := []struct {
		a        string
		b        string
		expected int
	}{
		{"a", "", 1},
		{"a", "ab", 1},
		{"b", "ab", 1},
		{"abc", "ac", 1},
		{"abc", "bc", 1},
		{"abc", "ab", 1},
	}
	for _, each := range testCases {
		t.Run(fmt.Sprintf("%s_vs_%s", each.a, each.b), func(t *testing.T) {
			actual := Distance(each.a, each.b)
			if actual != each.expected {
				t.Errorf("a = \"%s\" b = \"%s\", but %d != %d", each.a, each.b, actual, each.expected)
			}

		})
	}
}

func TestTwone(t *testing.T) {
	testCases := []struct {
		a        string
		b        string
		expected int
	}{
		{"ab", "", 2},
		{"", "ab", 2},
		{"b", "abb", 2},
		{"abc", "c", 2},
		{"abc", "qc", 2},
		{"ab", "cd", 2},
	}
	for _, each := range testCases {
		t.Run(fmt.Sprintf("%s_vs_%s", each.a, each.b), func(t *testing.T) {
			actual := Distance(each.a, each.b)
			if actual != each.expected {
				t.Errorf("a = \"%s\" b = \"%s\", but %d != %d", each.a, each.b, actual, each.expected)
			}

		})
	}
}
