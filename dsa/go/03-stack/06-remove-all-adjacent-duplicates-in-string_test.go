package stack

import "testing"

func removeDuplicatesString(s string) string {
    st := make([]rune, 0, len(s))
    for _, ch := range s {
        if len(st) > 0 && st[len(st)-1] == ch { st = st[:len(st)-1] } else { st = append(st, ch) }
    }
    return string(st)
}

func TestRemoveDuplicatesString(t *testing.T) {
    if removeDuplicatesString("abbaca") != "ca" { t.Fatalf("expected ca") }
    if removeDuplicatesString("") != "" { t.Fatalf("expected empty") }
}

