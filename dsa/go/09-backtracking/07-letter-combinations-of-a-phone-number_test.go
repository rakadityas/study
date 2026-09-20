package backtracking

import "testing"

// letterCombinations returns every string spellable from a phone keypad digit string.
// Approach: backtrack one digit at a time, branching over the letters mapped to that digit
// time: O(4^n), space: O(n) — up to four letters per digit; recursion depth is the digit count
func letterCombinations(digits string) []string {
    if len(digits) == 0 { return []string{} }

    dictPhone := map[byte][]byte{
        '2': []byte("abc"),
        '3': []byte("def"),
        '4': []byte("ghi"),
        '5': []byte("jkl"),
        '6': []byte("mno"),
        '7': []byte("pqrs"),
        '8': []byte("tuv"),
        '9': []byte("wxyz"),
    }

    res := []string{}

    var backtracking func(combination []byte, digitsIdx int)
    backtracking = func(combination []byte, digitsIdx int) {
        if digitsIdx == len(digits) {
            res = append(res, string(combination))
            return
        }

        letters := dictPhone[digits[digitsIdx]]
        for i := range letters {
            combination = append(combination, letters[i])
            backtracking(combination, digitsIdx+1)
            combination = combination[:len(combination)-1]
        }
    }

    backtracking([]byte{}, 0)
    return res
}

func TestLetterCombinations(t *testing.T) {
    cases := []struct {
        in   string
        want []string
    }{
        {"23", []string{"ad", "ae", "af", "bd", "be", "bf", "cd", "ce", "cf"}},
        {"", []string{}},
        {"2", []string{"a", "b", "c"}},
    }

    for _, c := range cases {
        got := letterCombinations(c.in)
        if len(got) != len(c.want) {
            t.Fatalf("letterCombinations(%q) = %v, want %v", c.in, got, c.want)
        }
        for i := range got {
            if got[i] != c.want[i] {
                t.Fatalf("letterCombinations(%q) = %v, want %v", c.in, got, c.want)
            }
        }
    }
}
