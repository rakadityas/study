package linked_list

import (
	"reflect"
	"testing"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

// reverseList reverses a singly linked list.
func reverseList(head *ListNode) *ListNode {
	var prev *ListNode
	cur := head
	for cur != nil {
		next := cur.Next
		cur.Next = prev
		prev = cur
		cur = next
	}
	return prev
}

func buildList(vals []int) *ListNode {
	var head, tail *ListNode
	for _, v := range vals {
		node := &ListNode{Val: v}
		if head == nil {
			head = node
			tail = node
		} else {
			tail.Next = node
			tail = node
		}
	}
	return head
}

func toSlice(head *ListNode) []int {
	res := []int{}
	for head != nil {
		res = append(res, head.Val)
		head = head.Next
	}
	return res
}

func TestReverseList(t *testing.T) {
	cases := []struct {
		in   []int
		want []int
	}{
		{[]int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
		{[]int{}, []int{}},
		{[]int{1}, []int{1}},
	}
	for _, c := range cases {
		got := toSlice(reverseList(buildList(c.in)))
		if !reflect.DeepEqual(got, c.want) {
			t.Fatalf("reverseList(%v) = %v; want %v", c.in, got, c.want)
		}
	}
}
