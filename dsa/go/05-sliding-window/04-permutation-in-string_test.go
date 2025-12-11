package sliding_window

import "testing"

func checkInclusion(s1, s2 string) bool {
    if len(s1) > len(s2) { return false }
    cnt1 := make([]int, 26)
    cnt2 := make([]int, 26)
    for i := 0; i < len(s1); i++ { cnt1[s1[i]-'a']++ ; cnt2[s2[i]-'a']++ }
    matches := 0
    for i := 0; i < 26; i++ { if cnt1[i] == cnt2[i] { matches++ } }
    for i := len(s1); i < len(s2); i++ {
        if matches == 26 { return true }
        r := s2[i] - 'a'
        l := s2[i-len(s1)] - 'a'
        cnt2[r]++
        if cnt2[r] == cnt1[r] { matches++ } else if cnt2[r]-1 == cnt1[r] { matches-- }
        cnt2[l]--
        if cnt2[l] == cnt1[l] { matches++ } else if cnt2[l]+1 == cnt1[l] { matches-- }
    }
    return matches == 26
}

func TestCheckInclusion(t *testing.T) {
    if !checkInclusion("ab", "eidbaooo") { t.Fatalf("expected true") }
    if checkInclusion("ab", "eidboaoo") { t.Fatalf("expected false") }
}

