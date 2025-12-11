package stack

import (
    "reflect"
    "testing"
)

func generateParenthesis(n int) []string {
    res := []string{}
    var dfs func(open, close int, cur []byte)
    dfs = func(open, close int, cur []byte) {
        if len(cur) == 2*n { res = append(res, string(cur)); return }
        if open < n {
            dfs(open+1, close, append(cur, '('))
        }
        if close < open {
            dfs(open, close+1, append(cur, ')'))
        }
    }
    dfs(0, 0, []byte{})
    return res
}

func TestGenerateParenthesis(t *testing.T) {
    got := generateParenthesis(3)
    want := []string{"((()))","(()())","(())()","()(())","()()()"}
    if len(got) != len(want) { t.Fatalf("unexpected length: %d", len(got)) }
    // Simple membership check
    m := map[string]bool{}
    for _, w := range want { m[w] = true }
    for _, g := range got { if !m[g] { t.Fatalf("missing %s", g) } }
    // Ensure all unique
    seen := map[string]bool{}
    for _, g := range got { if seen[g] { t.Fatalf("duplicate %s", g) } ; seen[g]=true }
    _ = reflect.DeepEqual
}

