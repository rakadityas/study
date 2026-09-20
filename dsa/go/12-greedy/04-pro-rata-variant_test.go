package greedy

import (
    "math"
    "testing"
)

// Share is one account's prorated cut of a distribution.
type Share struct {
    Account string
    Amount  float64
}

// releaseProportional splits the available balance across claimants in proportion to what each
// is owed: their_share / total_owed * available_balance.
// Approach: one pass to total the claims, then one pass to scale each claim by the same ratio.
// Unlike the FIFO holding account, a short balance is split fairly rather than blocking.
// time: O(n), space: O(n) — two linear passes; the result holds one entry per claimant
func releaseProportional(balance int, pairs []Payout) (int, []Share) {
    totalOwed := 0
    for _, p := range pairs {
        totalOwed += p.Amount
    }

    if totalOwed == 0 { return balance, []Share{} }

    toDistribute := balance
    if totalOwed < toDistribute { toDistribute = totalOwed }

    result := []Share{}
    for _, p := range pairs {
        proratedAmt := (float64(p.Amount) / float64(totalOwed)) * float64(toDistribute)
        result = append(result, Share{p.Account, proratedAmt})
    }

    return balance - toDistribute, result
}

func equalShares(a, b []Share) bool {
    if len(a) != len(b) { return false }
    for i := range a {
        if a[i].Account != b[i].Account || math.Abs(a[i].Amount-b[i].Amount) > 0.005 { return false }
    }
    return true
}

func TestReleaseProportional(t *testing.T) {
    pairs := []Payout{{"A", 600}, {"B", 400}}

    cases := []struct {
        balance     int
        wantBalance int
        wantShares  []Share
    }{
        {1000, 0, []Share{{"A", 600}, {"B", 400}}},   // exact cover
        {500, 0, []Share{{"A", 300}, {"B", 200}}},    // half short, split 60/40
        {100, 0, []Share{{"A", 60}, {"B", 40}}},      // deeply short, same ratio
        {1500, 500, []Share{{"A", 600}, {"B", 400}}}, // surplus stays in the account
    }

    for _, c := range cases {
        gotBalance, gotShares := releaseProportional(c.balance, pairs)
        if gotBalance != c.wantBalance || !equalShares(gotShares, c.wantShares) {
            t.Fatalf("releaseProportional(%d, %v) = (%d, %v), want (%d, %v)",
                c.balance, pairs, gotBalance, gotShares, c.wantBalance, c.wantShares)
        }
    }

    // no claims at all leaves the balance untouched
    if gotBalance, gotShares := releaseProportional(250, []Payout{}); gotBalance != 250 || len(gotShares) != 0 {
        t.Fatalf("releaseProportional(250, []) = (%d, %v), want (250, [])", gotBalance, gotShares)
    }
}
