package sliding_window

import "testing"

// lengthOfLongestSubstring returns the length of the longest substring without repeating characters.
// Approach: last-seen-index — jump the left pointer directly to last_seen[char]+1 instead of shrinking one step at a time
// time: O(n), space: O(n) — single pass; map stores at most all distinct characters
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

// lengthOfLongestSubstringShrink answers the same question one step at a time.
// Approach: shrink-window — while the incoming character is already in the window, evict the
// leftmost character until the duplicate is gone, then expand right
// time: O(n), space: O(n) — each character enters and leaves the map at most once
func lengthOfLongestSubstringShrink(s string) int {
    dictHistory := map[byte]bool{}
    res, l := 0, 0

    for i := 0; i < len(s); i++ {
        for dictHistory[s[i]] {
            delete(dictHistory, s[l])
            l++
        }

        if i-l+1 > res { res = i - l + 1 }
        dictHistory[s[i]] = true
    }

    return res
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
        if got := lengthOfLongestSubstringShrink(c.s); got != c.want {
            t.Fatalf("lengthOfLongestSubstringShrink(%q) = %d; want %d", c.s, got, c.want)
        }
    }
}

