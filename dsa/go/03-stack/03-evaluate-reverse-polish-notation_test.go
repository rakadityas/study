package stack

import "testing"

func evalRPN(tokens []string) int {
    st := []int{}
    for _, tok := range tokens {
        switch tok {
        case "+": st[len(st)-2] = st[len(st)-2] + st[len(st)-1]; st = st[:len(st)-1]
        case "-": st[len(st)-2] = st[len(st)-2] - st[len(st)-1]; st = st[:len(st)-1]
        case "*": st[len(st)-2] = st[len(st)-2] * st[len(st)-1]; st = st[:len(st)-1]
        case "/": st[len(st)-2] = st[len(st)-2] / st[len(st)-1]; st = st[:len(st)-1]
        default:
            // parse int
            sign := 1
            s := tok
            if len(s) > 0 && s[0] == '-' { sign = -1; s = s[1:] }
            v := 0
            for i := 0; i < len(s); i++ { v = v*10 + int(s[i]-'0') }
            st = append(st, sign*v)
        }
    }
    return st[0]
}

func TestEvalRPN(t *testing.T) {
    if evalRPN([]string{"2","1","+","3","*"}) != 9 { t.Fatalf("expected 9") }
    if evalRPN([]string{"4","13","5","/","+"}) != 6 { t.Fatalf("expected 6") }
}

