package greedy

import "testing"

// Account is a named bank account with a balance.
type Account struct {
    Name    string
    Balance int
}

// Transfer moves Amount from one account to another.
type Transfer struct {
    From   string
    To     string
    Amount int
}

// determineTransfers works out the transfers needed so every account holds at least 100.
// Approach: split accounts into surplus sources and under-funded targets, then walk both lists
// with a moving source cursor — each target drains sources greedily until its gap is closed,
// and a source is retired the moment it drops to exactly 100
// time: O(n), space: O(n) — one pass to split, then a two-pointer style traversal
func determineTransfers(accounts []Account) []Transfer {
    sources := []Account{}
    targets := []Account{}

    for _, acc := range accounts {
        if acc.Balance > 100 {
            sources = append(sources, acc)
        } else {
            targets = append(targets, acc)
        }
    }

    res := []Transfer{}
    sourceIdx := 0

    for _, target := range targets {
        gap := 100 - target.Balance

        for gap > 0 && sourceIdx < len(sources) {
            surplus := sources[sourceIdx].Balance - 100

            trfAmt := surplus
            if gap < trfAmt { trfAmt = gap }
            gap -= trfAmt

            res = append(res, Transfer{sources[sourceIdx].Name, target.Name, trfAmt})

            sources[sourceIdx].Balance -= trfAmt
            if sources[sourceIdx].Balance == 100 { sourceIdx++ }
        }
    }

    return res
}

func equalTransfers(a, b []Transfer) bool {
    if len(a) != len(b) { return false }
    for i := range a {
        if a[i] != b[i] { return false }
    }
    return true
}

func TestDetermineTransfers(t *testing.T) {
    cases := []struct {
        in   []Account
        want []Transfer
    }{
        { // two accounts, a single deduction
            []Account{{"a", 80}, {"b", 120}},
            []Transfer{{"b", "a", 20}},
        },
        { // one source covers more than one target
            []Account{{"a", 140}, {"b", 80}, {"c", 80}},
            []Transfer{{"a", "b", 20}, {"a", "c", 20}},
        },
        { // a target drains two sources
            []Account{{"a", 110}, {"b", 95}, {"c", 90}, {"d", 105}},
            []Transfer{{"a", "b", 5}, {"a", "c", 5}, {"d", "c", 5}},
        },
        {
            []Account{{"a", 200}, {"b", 95}, {"c", 90}, {"d", 50}},
            []Transfer{{"a", "b", 5}, {"a", "c", 10}, {"a", "d", 50}},
        },
        { // no surplus anywhere
            []Account{{"a", 50}, {"b", 95}, {"c", 90}, {"d", 50}},
            []Transfer{},
        },
    }

    for _, c := range cases {
        if got := determineTransfers(c.in); !equalTransfers(got, c.want) {
            t.Fatalf("determineTransfers(%v) = %v, want %v", c.in, got, c.want)
        }
    }
}
