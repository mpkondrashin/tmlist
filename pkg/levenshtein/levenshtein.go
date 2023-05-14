//////////////////////////////////////////////////////////////////////////
//
//  (c) TMList 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)
//  Copyright under MIT Lincese. Please see LICENSE file for details
//
//  levenshtein.go - calculate levenshtein distance
//
//////////////////////////////////////////////////////////////////////////

package levenshtein

func min(a, b, c int) int {
	if a < b {
		if c < a {
			return c
		}
		return a
	}
	if c < b {
		return c
	}
	return b
}

func Distance(a, b string) int {
	N := len(a)
	M := len(b)
	s := make([]int, N+1)
	t := make([]int, N+1)
	for i := 1; i <= N; i++ {
		s[i] = i
	}
	for j := 1; j <= M; j++ {
		t[0] = j
		for i := 1; i <= N; i++ {
			p := s[i-1]
			if a[i-1] != b[j-1] {
				p++
			}
			t[i] = min(p, t[i-1]+1, s[i]+1)
		}
		s, t = t, s
	}
	return s[N]
}
