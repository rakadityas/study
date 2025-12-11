package array_hashes

import "testing"

func isAnagram(s string, t string) bool {
	if len(s) != len(t) {
		return false
	}
	mapS := make(map[rune]int)
	mapT := make(map[rune]int)

	for _, char := range s {
		mapS[char]++
	}
	for _, char := range t {
		mapT[char]++
	}

	for key := range mapS {
		if mapS[key] != mapT[key] {
			return false
		}
	}

	return true
}

func TestIsAnagram(t *testing.T) {
	testCases := []struct {
		name     string
		s        string
		t        string
		expected bool
	}{
		{
			name:     "anagram",
			s:        "anagram",
			t:        "nagaram",
			expected: true,
		},
		{
			name:     "not anagram",
			s:        "rat",
			t:        "car",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isAnagram(tc.s, tc.t)
			if result != tc.expected {
				t.Errorf("isAnagram(%v, %v) = %v; want %v", tc.s, tc.t, result, tc.expected)
			}
		})
	}
}
