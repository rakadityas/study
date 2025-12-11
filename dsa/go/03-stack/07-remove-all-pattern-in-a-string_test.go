package stack

import "testing"

func removeOccurrencesPatterns(s string) string {
    patterns := map[string]bool{"ab":true,"bc":true,"cd":true,"dc":true,"cb":true,"ba":true}
    st := make([]rune, 0, len(s))
    for _, ch := range s {
        if len(st) > 0 {
            pair := string([]rune{st[len(st)-1], ch})
            if patterns[pair] { st = st[:len(st)-1]; continue }
        }
        st = append(st, ch)
    }
    return string(st)
}

func TestRemoveOccurrencesPatterns(t *testing.T) {
    if removeOccurrencesPatterns("abccba") != "ca" { t.Fatalf("expected ca") }
    if removeOccurrencesPatterns("aabccba") != "aca" { t.Fatalf("expected aca") }
    if removeOccurrencesPatterns("") != "" { t.Fatalf("expected empty") }
}

