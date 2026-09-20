package greedy

import (
    "math"
    "sort"
    "testing"
)

// Provider quotes a rate (EUR received per USD sent) and the maximum USD volume it can absorb.
type Provider struct {
    Name    string
    Rate    float64
    MaxVol int
}

// Fill is one provider's slice of a conversion.
type Fill struct {
    Name        string
    USDUsed     int
    EURReceived float64
}

// bestExecution fills a USD->EUR conversion across providers at the best achievable total.
// Approach: greedy — sort providers by rate descending and take as much as each can absorb
// until the requested amount is filled. Because the rate is fixed per provider, taking the
// best rate first is always optimal. If total liquidity falls short, it fills what it can.
// time: O(n log n), space: O(n) — dominated by the sort
func bestExecution(amountUSD int, providers []Provider) (float64, []Fill) {
    sortedProviders := append([]Provider{}, providers...)
    sort.SliceStable(sortedProviders, func(i, j int) bool {
        return sortedProviders[i].Rate > sortedProviders[j].Rate
    })

    total := 0.0
    breakdown := []Fill{}
    remaining := amountUSD

    for _, p := range sortedProviders {
        if remaining <= 0 { break }

        usdUsed := remaining
        if p.MaxVol < usdUsed { usdUsed = p.MaxVol }
        if usdUsed == 0 { continue }

        eurReceived := math.Round(float64(usdUsed)*p.Rate*100) / 100
        total += eurReceived

        breakdown = append(breakdown, Fill{p.Name, usdUsed, eurReceived})
        remaining -= usdUsed
    }

    return total, breakdown
}

func equalFills(a, b []Fill) bool {
    if len(a) != len(b) { return false }
    for i := range a {
        if a[i].Name != b[i].Name || a[i].USDUsed != b[i].USDUsed ||
            math.Abs(a[i].EURReceived-b[i].EURReceived) > 0.005 {
            return false
        }
    }
    return true
}

func TestBestExecution(t *testing.T) {
    providers := []Provider{
        {"Bank A", 1.08, 10_000},
        {"Bank B", 1.07, 5_000},
        {"Bank C", 1.09, 3_000},
    }

    cases := []struct {
        amount    int
        wantTotal float64
        wantFills []Fill
    }{
        { // Bank C first (best rate, 3k), then Bank A (10k), then Bank B (2k)
            15_000, 16_210.0,
            []Fill{{"Bank C", 3_000, 3_270.0}, {"Bank A", 10_000, 10_800.0}, {"Bank B", 2_000, 2_140.0}},
        },
        { // exact fill from the single best provider
            3_000, 3_270.0,
            []Fill{{"Bank C", 3_000, 3_270.0}},
        },
        { // insufficient liquidity — 18k available against a 20k ask
            20_000, 19_420.0,
            []Fill{{"Bank C", 3_000, 3_270.0}, {"Bank A", 10_000, 10_800.0}, {"Bank B", 5_000, 5_350.0}},
        },
    }

    for _, c := range cases {
        total, fills := bestExecution(c.amount, providers)
        if !equalFills(fills, c.wantFills) {
            t.Fatalf("bestExecution(%d) fills = %v, want %v", c.amount, fills, c.wantFills)
        }
        if math.Abs(total-c.wantTotal) > 0.005 {
            t.Fatalf("bestExecution(%d) total = %.2f, want %.2f", c.amount, total, c.wantTotal)
        }
    }
}
