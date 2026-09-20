package stack

import "testing"

// atomicWeights is the known weight per atom symbol.
var atomicWeights = map[byte]int{
    'C': 12,
    'H': 1,
}

// countAtoms totals the atomic weight of a formula, honouring parenthesised groups and counts.
// Approach: a stack holding the running weight of each nesting level — '(' opens a new level,
// ')' multiplies the closed level by its count and folds it into the parent
// time: O(n), space: O(d) — each character is read once; d is the max nesting depth (d <= n)
func countAtoms(formula string) int {
    stack := []int{0} // running weight total at each nesting level
    i, n := 0, len(formula)

    for i < n {
        char := formula[i]

        switch {
        case char == '(':
            stack = append(stack, 0)
            i++

        case char == ')':
            var count int
            count, i = readCount(formula, i+1)
            groupWeight := stack[len(stack)-1] * count
            stack = stack[:len(stack)-1]
            stack[len(stack)-1] += groupWeight

        default:
            weight := atomicWeights[char]
            var count int
            count, i = readCount(formula, i+1)
            stack[len(stack)-1] += weight * count
        }
    }

    return stack[0]
}

// readCount reads the digit run starting at i, defaulting to 1 when there is none.
func readCount(formula string, i int) (int, int) {
    start := i
    for i < len(formula) && formula[i] >= '0' && formula[i] <= '9' {
        i++
    }

    if i == start { return 1, i }

    count := 0
    for _, d := range formula[start:i] {
        count = count*10 + int(d-'0')
    }
    return count, i
}

func TestCountAtoms(t *testing.T) {
    cases := []struct {
        in   string
        want int
    }{
        {"CH4", 16},
        {"C(H4)H4", 20},
        {"H10", 10}, // multi-digit count
    }

    for _, c := range cases {
        if got := countAtoms(c.in); got != c.want {
            t.Fatalf("countAtoms(%q) = %d, want %d", c.in, got, c.want)
        }
    }
}
