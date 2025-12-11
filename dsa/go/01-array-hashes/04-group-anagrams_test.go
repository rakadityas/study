package array_hashes

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/rakadityas/study/dsa/go/common"
)

func groupAnagrams(strs []string) [][]string {
	mapRespGroup := make(map[string][]string)
	res := [][]string{}

	for _, word := range strs {
		// count characters
		arrWordCounter := make([]int, 26)
		for i := 0; i < len(word); i++ {
			arrWordCounter[word[i]-'a']++
		}

		// convert to string key
		key := fmt.Sprint(arrWordCounter)

		// group by key
		mapRespGroup[key] = append(mapRespGroup[key], word)
	}

	// collect results
	for _, group := range mapRespGroup {
		res = append(res, group)
	}

	return res
}

func TestGroupAnagrams(t *testing.T) {
	testCases := []struct {
		name     string
		strs     []string
		expected [][]string
	}{
		{
			name:     "group anagrams",
			strs:     []string{"eat", "tea", "tan", "ate", "nat", "bat"},
			expected: [][]string{{"bat"}, {"nat", "tan"}, {"ate", "eat", "tea"}},
		},
		{
			name:     "group anagrams",
			strs:     []string{""},
			expected: [][]string{{""}},
		},
		{
			name:     "group anagrams",
			strs:     []string{"a"},
			expected: [][]string{{"a"}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := groupAnagrams(tc.strs)

			common.Sort2DStringSlice(result)
			common.Sort2DStringSlice(tc.expected)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("groupAnagrams(%v) = %v; want %v", tc.strs, result, tc.expected)
			}
		})
	}
}
