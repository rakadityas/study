package sliding_window

import "testing"

// Port of SolutionOne.countOfSubstring from Python: counts specific metric per given logic.
func countOfSubstring(s string) int {
    l := 0
    mapHistory := map[byte]bool{}
    numSubstring := 0
    if len(s) == 0 { return 0 }
    for i := 0; i < len(s); i++ {
        if mapHistory[s[i]] {
            numSubstring++
            for l != i {
                delete(mapHistory, s[l])
                l++
            }
        }
        mapHistory[s[i]] = true
    }
    return numSubstring + 1
}

func TestCountOfSubstring(t *testing.T) {
    if countOfSubstring("abcabcbb") != 4 { t.Fatalf("expected 4") }
    if countOfSubstring("bbbbb") != 5 { t.Fatalf("expected 5") }
    if countOfSubstring("") != 0 { t.Fatalf("expected 0") }
}

