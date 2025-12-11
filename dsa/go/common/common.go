package common

import (
	"sort"
	"strings"
)

func Sort2DStringSlice(s [][]string) {
	// sort each inner slice
	for i := range s {
		sort.Strings(s[i])
	}

	// sort the outer slice based on joined string
	sort.Slice(s, func(i, j int) bool {
		return strings.Join(s[i], "") < strings.Join(s[j], "")
	})
}
