package linked_list

import (
    "reflect"
    "testing"
)

func mergeTwoLists(l1 *ListNode, l2 *ListNode) *ListNode {
    dummy := &ListNode{}
    cur := dummy
    for l1 != nil && l2 != nil {
        if l1.Val < l2.Val { cur.Next = l1; l1 = l1.Next } else { cur.Next = l2; l2 = l2.Next }
        cur = cur.Next
    }
    if l1 != nil { cur.Next = l1 } else { cur.Next = l2 }
    return dummy.Next
}

func TestMergeTwoLists(t *testing.T) {
    a := buildList([]int{1,2,4})
    b := buildList([]int{1,3,4})
    got := toSlice(mergeTwoLists(a,b))
    want := []int{1,1,2,3,4,4}
    if !reflect.DeepEqual(got, want) { t.Fatalf("got %v want %v", got, want) }
}

