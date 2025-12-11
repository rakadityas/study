package stack

import "testing"

type MinStack struct {
    s   []int
    min []int
}

func ConstructorMinStack() MinStack { return MinStack{s: []int{}, min: []int{}} }
func (ms *MinStack) Push(x int) { ms.s = append(ms.s, x); if len(ms.min)==0 || x <= ms.min[len(ms.min)-1] { ms.min = append(ms.min, x) } }
func (ms *MinStack) Pop() { v := ms.s[len(ms.s)-1]; ms.s = ms.s[:len(ms.s)-1]; if v == ms.min[len(ms.min)-1] { ms.min = ms.min[:len(ms.min)-1] } }
func (ms *MinStack) Top() int { return ms.s[len(ms.s)-1] }
func (ms *MinStack) GetMin() int { return ms.min[len(ms.min)-1] }

func TestMinStack(t *testing.T) {
    ms := ConstructorMinStack()
    ms.Push(-2); ms.Push(0); ms.Push(-3)
    if ms.GetMin() != -3 { t.Fatalf("min should be -3") }
    ms.Pop()
    if ms.Top() != 0 { t.Fatalf("top should be 0") }
    if ms.GetMin() != -2 { t.Fatalf("min should be -2") }
}

