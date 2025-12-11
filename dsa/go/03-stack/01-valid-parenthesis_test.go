package stack

import "testing"

// isValidParentheses checks if parentheses are valid using a stack.
func isValidParentheses(s string) bool {
    stack := make([]rune, 0, len(s))
    pairs := map[rune]rune{')': '(', ']': '[', '}': '{'}
    for _, ch := range s {
        switch ch {
        case '(', '[', '{':
            stack = append(stack, ch)
        case ')', ']', '}':
            if len(stack) == 0 || stack[len(stack)-1] != pairs[ch] {
                return false
            }
            stack = stack[:len(stack)-1]
        }
    }
    return len(stack) == 0
}

func TestIsValidParentheses(t *testing.T) {
    cases := []struct{ s string; want bool }{
        {"()", true},
        {"()[]{}", true},
        {"(]", false},
        {"([)]", false},
        {"{[]}", true},
    }
    for _, c := range cases {
        if got := isValidParentheses(c.s); got != c.want {
            t.Fatalf("isValidParentheses(%q) = %v; want %v", c.s, got, c.want)
        }
    }
}

