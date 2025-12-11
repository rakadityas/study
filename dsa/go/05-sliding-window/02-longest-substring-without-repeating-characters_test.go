package sliding_window

import "testing"

// lengthOfLongestSubstring returns the length of the longest substring without repeating characters.
func lengthOfLongestSubstring(s string) int {
    last := make(map[rune]int)
    maxLen := 0
    start := 0
    for i, ch := range []rune(s) {
        if prev, ok := last[ch]; ok && prev >= start {
            start = prev + 1
        }
        last[ch] = i
        if i-start+1 > maxLen {
            maxLen = i - start + 1
        }
    }
    return maxLen
}

func TestLengthOfLongestSubstring(t *testing.T) {
    cases := []struct{ s string; want int }{
        {"abcabcbb", 3},
        {"bbbbb", 1},
        {"pwwkew", 3},
        {"", 0},
        {" ", 1},
        {"dvdf", 3},
    }
    for _, c := range cases {
        if got := lengthOfLongestSubstring(c.s); got != c.want {
            t.Fatalf("lengthOfLongestSubstring(%q) = %d; want %d", c.s, got, c.want)
        }
    }
}

