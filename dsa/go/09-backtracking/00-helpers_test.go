package backtracking

import (
    "sort"
    "strconv"
    "strings"
)

// normalize2D sorts a slice of int slices into a canonical order so results can be
// compared without depending on the order backtracking happens to emit them in.
func normalize2D(s [][]int) []string {
    keys := make([]string, 0, len(s))
    for _, inner := range s {
        parts := make([]string, 0, len(inner))
        for _, v := range inner {
            parts = append(parts, strconv.Itoa(v))
        }
        keys = append(keys, strings.Join(parts, ","))
    }
    sort.Strings(keys)
    return keys
}

// equalIgnoringOrder compares two result sets ignoring the order of the outer slice.
func equalIgnoringOrder(got, want [][]int) bool {
    a, b := normalize2D(got), normalize2D(want)
    if len(a) != len(b) { return false }
    for i := range a {
        if a[i] != b[i] { return false }
    }
    return true
}
