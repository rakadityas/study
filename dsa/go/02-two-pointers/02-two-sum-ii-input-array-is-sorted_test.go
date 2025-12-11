package two_pointers

import (
    "reflect"
    "testing"
)

// twoSumII returns 1-based indices of two numbers that add up to target.
func twoSumII(numbers []int, target int) []int {
    i, j := 0, len(numbers)-1
    for i < j {
        sum := numbers[i] + numbers[j]
        if sum == target {
            return []int{i + 1, j + 1}
        } else if sum < target {
            i++
        } else {
            j--
        }
    }
    return nil
}

func TestTwoSumII(t *testing.T) {
    cases := []struct{ nums []int; target int; want []int }{
        {[]int{2, 7, 11, 15}, 9, []int{1, 2}},
        {[]int{2, 3, 4}, 6, []int{1, 3}},
        {[]int{-1, 0}, -1, []int{1, 2}},
    }
    for _, c := range cases {
        got := twoSumII(c.nums, c.target)
        if !reflect.DeepEqual(got, c.want) {
            t.Fatalf("twoSumII(%v, %d) = %v; want %v", c.nums, c.target, got, c.want)
        }
    }
}

