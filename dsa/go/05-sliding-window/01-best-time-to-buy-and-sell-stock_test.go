package sliding_window

import "testing"

// time: O(n), space: O(1) — track running minimum price; profit = current price − min so far
func maxProfit(prices []int) int {
    minPrice := int(^uint(0) >> 1) // max int
    best := 0
    for _, p := range prices {
        if p < minPrice { minPrice = p }
        if p-minPrice > best { best = p - minPrice }
    }
    return best
}

func TestMaxProfit(t *testing.T) {
    if maxProfit([]int{7,1,5,3,6,4}) != 5 { t.Fatalf("expected 5") }
    if maxProfit([]int{7,6,4,3,1}) != 0 { t.Fatalf("expected 0") }
}

