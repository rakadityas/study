package backtracking

import "testing"

// exist reports whether word can be spelled by walking adjacent cells without reusing a cell.
// Approach: DFS from every cell, marking the current cell with a sentinel while descending so the
// path cannot revisit it, then restoring it on the way out
// time: O(m*n*4^len(word)), space: O(len(word)) — recursion depth is the word length
func exist(board [][]byte, word string) bool {
    if len(board) == 0 { return false }

    rows, cols := len(board), len(board[0])

    var backtracking func(r, c, i int) bool
    backtracking = func(r, c, i int) bool {
        if i == len(word) { return true }

        if r < 0 || r >= rows || c < 0 || c >= cols || board[r][c] != word[i] { return false }

        temp := board[r][c]
        board[r][c] = '#'

        found := backtracking(r+1, c, i+1) || backtracking(r-1, c, i+1) ||
            backtracking(r, c+1, i+1) || backtracking(r, c-1, i+1)

        board[r][c] = temp

        return found
    }

    for r := 0; r < rows; r++ {
        for c := 0; c < cols; c++ {
            if backtracking(r, c, 0) { return true }
        }
    }

    return false
}

func TestExist(t *testing.T) {
    newBoard := func() [][]byte {
        return [][]byte{
            []byte("ABCE"),
            []byte("SFCS"),
            []byte("ADEE"),
        }
    }

    cases := []struct {
        word string
        want bool
    }{
        {"ABCCED", true},
        {"SEE", true},
        {"ABCB", false},
    }

    for _, c := range cases {
        if got := exist(newBoard(), c.word); got != c.want {
            t.Fatalf("exist(board, %q) = %v, want %v", c.word, got, c.want)
        }
    }
}
