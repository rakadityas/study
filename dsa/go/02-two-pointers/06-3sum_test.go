package two_pointers

import (
    "reflect"
    "sort"
    "testing"
)

// threeSum finds unique triplets that sum to zero.
func threeSum(nums []int) [][]int {
    sort.Ints(nums)
    res := [][]int{}
    for i := 0; i < len(nums); i++ {
        if i > 0 && nums[i] == nums[i-1] { continue }
        l, r := i+1, len(nums)-1
        for l < r {
            s := nums[i] + nums[l] + nums[r]
            if s == 0 {
                res = append(res, []int{nums[i], nums[l], nums[r]})
                l++
                for l < r && nums[l] == nums[l-1] { l++ }
                r--
                for l < r && nums[r] == nums[r+1] { r-- }
            } else if s < 0 { l++ } else { r-- }
        }
    }
    return res
}

func TestThreeSum(t *testing.T) {
    got := threeSum([]int{-1,0,1,2,-1,-4})
    want := [][]int{{-1,-1,2},{-1,0,1}}
    sort.Slice(got, func(i, j int) bool {
        a, b := got[i], got[j]
        for k := 0; k < 3; k++ { if a[k] != b[k] { return a[k] < b[k] } }
        return false
    })
    sort.Slice(want, func(i, j int) bool {
        a, b := want[i], want[j]
        for k := 0; k < 3; k++ { if a[k] != b[k] { return a[k] < b[k] } }
        return false
    })
    if !reflect.DeepEqual(got, want) { t.Fatalf("threeSum got %v want %v", got, want) }
}

