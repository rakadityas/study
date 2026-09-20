package graph

import "testing"

// TicTacToe is a two-player 3x3 game. Players alternate marking X or O, and the first to
// align three of their symbols horizontally, vertically or diagonally wins.
//
// time: O(1), space: O(1) — the board is a fixed 3x3, so every check is a constant number of lookups
type TicTacToe struct {
    board [3][3]string
    moves int
}

func ConstructorTicTacToe() *TicTacToe { return &TicTacToe{} }

// ExecuteMove applies a move and reports whether it was accepted plus the resulting game state.
// time: O(1), space: O(1)
func (g *TicTacToe) ExecuteMove(symbol string, x, y int) (bool, string) {
    if !g.updateBoard(symbol, x, y) { return false, "Ongoing" }

    g.moves++

    if g.boardCheck(symbol, x, y) { return true, symbol + " wins" }

    if g.moves == 9 { return true, "Draw" }

    return true, "Ongoing"
}

// updateBoard rejects out-of-range coordinates and already-filled cells.
// time: O(1), space: O(1)
func (g *TicTacToe) updateBoard(symbol string, x, y int) bool {
    if x < 0 || x > 2 || y < 0 || y > 2 { return false }

    if g.board[x][y] != "" { return false }

    g.board[x][y] = symbol
    return true
}

// boardCheck tests only the four lines the last move could have completed.
// time: O(1), space: O(1)
func (g *TicTacToe) boardCheck(symbol string, x, y int) bool {
    lines := [4][3][2]int{
        {{x, 0}, {x, 1}, {x, 2}}, // the played row
        {{0, y}, {1, y}, {2, y}}, // the played column
        {{0, 0}, {1, 1}, {2, 2}}, // main diagonal
        {{0, 2}, {1, 1}, {2, 0}}, // anti-diagonal
    }

    for _, line := range lines {
        if g.board[line[0][0]][line[0][1]] == symbol &&
            g.board[line[1][0]][line[1][1]] == symbol &&
            g.board[line[2][0]][line[2][1]] == symbol {
            return true
        }
    }

    return false
}

func TestTicTacToe(t *testing.T) {
    type move struct {
        symbol    string
        x, y      int
        wantOK    bool
        wantState string
    }

    games := [][]move{
        { // O wins on the middle column
            {"X", 0, 0, true, "Ongoing"},
            {"O", 1, 1, true, "Ongoing"},
            {"X", 1, 0, true, "Ongoing"},
            {"O", 2, 0, true, "Ongoing"},
            {"X", 0, 2, true, "Ongoing"},
            {"O", 0, 1, true, "Ongoing"},
            {"X", 0, 1, false, "Ongoing"}, // already filled
            {"X", 2, 2, true, "Ongoing"},
            {"O", 2, 1, true, "O wins"},
        },
        { // X wins on the main diagonal
            {"X", 1, 1, true, "Ongoing"},
            {"O", 0, 1, true, "Ongoing"},
            {"X", 1, 0, true, "Ongoing"},
            {"O", 1, 2, true, "Ongoing"},
            {"X", 0, 0, true, "Ongoing"},
            {"O", 2, 0, true, "Ongoing"},
            {"X", 2, 2, true, "X wins"},
        },
        { // full board, nobody aligned three
            {"X", 0, 0, true, "Ongoing"},
            {"O", 0, 1, true, "Ongoing"},
            {"X", 0, 2, true, "Ongoing"},
            {"O", 1, 1, true, "Ongoing"},
            {"X", 1, 0, true, "Ongoing"},
            {"O", 1, 2, true, "Ongoing"},
            {"X", 2, 1, true, "Ongoing"},
            {"O", 2, 0, true, "Ongoing"},
            {"X", 2, 2, true, "Draw"},
        },
        { // every coordinate out of range
            {"X", -1, 0, false, "Ongoing"},
            {"X", 0, -1, false, "Ongoing"},
            {"X", 3, 0, false, "Ongoing"},
            {"X", 0, 3, false, "Ongoing"},
            {"X", 3, 3, false, "Ongoing"},
        },
    }

    for gi, moves := range games {
        game := ConstructorTicTacToe()
        for mi, m := range moves {
            gotOK, gotState := game.ExecuteMove(m.symbol, m.x, m.y)
            if gotOK != m.wantOK || gotState != m.wantState {
                t.Fatalf("game %d move %d: ExecuteMove(%q, %d, %d) = (%v, %q), want (%v, %q)",
                    gi, mi, m.symbol, m.x, m.y, gotOK, gotState, m.wantOK, m.wantState)
            }
        }
    }
}
