package linked_list

import (
    "reflect"
    "testing"
)

func removeNthFromEnd(head *ListNode, n int) *ListNode {
    dummy := &ListNode{Next: head}
    fast, slow := dummy, dummy
    for i := 0; i < n; i++ { fast = fast.Next }
    for fast.Next != nil { fast = fast.Next; slow = slow.Next }
    slow.Next = slow.Next.Next
    return dummy.Next
}

func TestRemoveNthFromEnd(t *testing.T) {
    got := toSlice(removeNthFromEnd(buildList([]int{1,2,3,4,5}), 2))
    want := []int{1,2,3,5}
    if !reflect.DeepEqual(got, want) { t.Fatalf("got %v want %v", got, want) }
}

