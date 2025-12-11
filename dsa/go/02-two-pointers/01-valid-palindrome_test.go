package two_pointers

import (
	"strings"
	"testing"
)

// isPalindrome checks if s is a palindrome considering only alphanumeric characters and ignoring cases.
func isPalindrome(s string) bool {
	i, j := 0, len(s)-1
	for i < j {
		for i < j && !isAlphaNum(s[i]) {
			i++
		}
		for i < j && !isAlphaNum(s[j]) {
			j--
		}
		if i < j {
			if toLower(s[i]) != toLower(s[j]) {
				return false
			}
			i++
			j--
		}
	}
	return true
}

func isAlphaNum(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}

func toLower(b byte) byte { return strings.ToLower(string(b))[0] }

func TestIsPalindrome(t *testing.T) {
	cases := []struct {
		s    string
		want bool
	}{
		{"A man, a plan, a canal: Panama", true},
		{"race a car", false},
		{" ", true},
		{"0P", false},
		{"abba", true},
	}
	for _, c := range cases {
		if got := isPalindrome(c.s); got != c.want {
			t.Fatalf("isPalindrome(%q) = %v; want %v", c.s, got, c.want)
		}
	}
}
