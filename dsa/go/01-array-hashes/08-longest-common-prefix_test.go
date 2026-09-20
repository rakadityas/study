package array_hashes

import "testing"

// getSameLetter returns the longest common prefix by scanning column by column.
// Approach: brute force — for every column up to the shortest word, compare every word's
// character against the first word's and stop at the first mismatch
// time: O(n*m), space: O(m) — n words, m = shortest word length
//
// Note: the Python bank also has a zip(*words) variant. In Go there is no zip idiom, and this
// loop already bails on the first mismatched column, so the two collapse into the same code.
func getSameLetter(words []string) string {
    if len(words) == 0 { return "" }

    res := []byte{}
    smallestLen := len(words[0])
    for i := range words {
        if len(words[i]) < smallestLen { smallestLen = len(words[i]) }
    }

    for i := 0; i < smallestLen; i++ {
        for j := range words {
            if words[j][i] != words[0][i] { return string(res) }
        }

        res = append(res, words[0][i])
    }

    return string(res)
}

// getSameLetterOptimal compares only the lexicographic extremes of the list.
// Approach: the common prefix of the whole list equals the common prefix of its min and max —
// any divergent word sits strictly between them, so it cannot shrink the prefix further
// time: O(n*m), space: O(m) — one O(n) pass to find each extreme, then O(m) to compare them
func getSameLetterOptimal(words []string) string {
    if len(words) == 0 { return "" }

    first, last := words[0], words[0]
    for _, w := range words {
        if w < first { first = w }
        if w > last { last = w }
    }

    smallestLen := len(first)
    if len(last) < smallestLen { smallestLen = len(last) }

    res := []byte{}
    for i := 0; i < smallestLen; i++ {
        if first[i] != last[i] { break }
        res = append(res, first[i])
    }

    return string(res)
}

func TestLongestCommonPrefix(t *testing.T) {
    cases := []struct {
        in   []string
        want string
    }{
        {[]string{"APPLE", "APP", "APPLAUD"}, "APP"},
        {[]string{"BANANA", "BAN", "BLUE"}, "B"},
        {[]string{"BANANA", "APPLE", "CAR"}, ""},
        {[]string{}, ""},
    }

    for _, c := range cases {
        if got := getSameLetter(c.in); got != c.want {
            t.Fatalf("getSameLetter(%v) = %q, want %q", c.in, got, c.want)
        }
        if got := getSameLetterOptimal(c.in); got != c.want {
            t.Fatalf("getSameLetterOptimal(%v) = %q, want %q", c.in, got, c.want)
        }
    }
}
