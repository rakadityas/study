package stack

import "testing"

// removeOccurrencesPart removes all occurrences of part from s.
func removeOccurrencesPart(s, part string) string {
	st := []rune{}
	pr := []rune(part)
	plen := len(pr)
	for _, ch := range s {
		st = append(st, ch)
		if len(st) >= plen {
			match := true
			for i := 0; i < plen; i++ {
				if st[len(st)-plen+i] != pr[i] {
					match = false
					break
				}
			}
			if match {
				st = st[:len(st)-plen]
			}
		}
	}
	return string(st)
}

func TestRemoveOccurrencesPart(t *testing.T) {
	if removeOccurrencesPart("daabcbaabcbc", "abc") != "dab" {
		t.Fatalf("expected dab")
	}
}
