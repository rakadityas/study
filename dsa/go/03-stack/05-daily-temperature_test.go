package stack

import (
    "reflect"
    "testing"
)

// dailyTemperatures returns the number of days to wait for a warmer temperature.
func dailyTemperatures(temperatures []int) []int {
    n := len(temperatures)
    res := make([]int, n)
    stack := make([]int, 0, n) // indices, monotonic decreasing by temperature
    for i := 0; i < n; i++ {
        for len(stack) > 0 && temperatures[i] > temperatures[stack[len(stack)-1]] {
            idx := stack[len(stack)-1]
            stack = stack[:len(stack)-1]
            res[idx] = i - idx
        }
        stack = append(stack, i)
    }
    return res
}

func TestDailyTemperatures(t *testing.T) {
    cases := []struct{ temps []int; want []int }{
        {[]int{73, 74, 75, 71, 69, 72, 76, 73}, []int{1, 1, 4, 2, 1, 1, 0, 0}},
        {[]int{30, 40, 50, 60}, []int{1, 1, 1, 0}},
        {[]int{30, 60, 90}, []int{1, 1, 0}},
    }
    for _, c := range cases {
        got := dailyTemperatures(c.temps)
        if !reflect.DeepEqual(got, c.want) {
            t.Fatalf("dailyTemperatures(%v) = %v; want %v", c.temps, got, c.want)
        }
    }
}

