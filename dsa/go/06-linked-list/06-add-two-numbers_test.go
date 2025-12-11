package linked_list

import (
    "reflect"
    "testing"
)

func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
    dummy := &ListNode{}
    cur := dummy
    carry := 0
    for l1 != nil || l2 != nil || carry != 0 {
        sum := carry
        if l1 != nil { sum += l1.Val; l1 = l1.Next }
        if l2 != nil { sum += l2.Val; l2 = l2.Next }
        cur.Next = &ListNode{Val: sum % 10}
        cur = cur.Next
        carry = sum / 10
    }
    return dummy.Next
}

func TestAddTwoNumbers(t *testing.T) {
    a := buildList([]int{2,4,3})
    b := buildList([]int{5,6,4})
    got := toSlice(addTwoNumbers(a,b))
    want := []int{7,0,8}
    if !reflect.DeepEqual(got, want) { t.Fatalf("got %v want %v", got, want) }
}

