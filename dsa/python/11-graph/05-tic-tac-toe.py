"""
Tic Tac Toe is a simple two-player game where players take turns marking X or O on a 3x3 grid.
The objective is to be the first to align three of your symbols horizontally, vertically, or
diagonally, upon which you win the game.
"""

from typing import Literal, Tuple, List

class TicTacToe:
    def __init__(self) -> None:
        self.board: List[List[str]] = [
            ["", "", ""],
            ["", "", ""],
            ["", "", ""],
        ]
        self.moves = 0

    def execute_move(
        self, symbol: Literal["X", "O"], x: int, y: int
    ) -> Tuple[bool, Literal["X wins", "O wins", "Draw", "Ongoing"]]:
        if not self.update_board(symbol, x, y):
            return False, "Ongoing"

        self.moves += 1

        if self.board_check(symbol, x, y):
            return True, f"{symbol} wins"

        if self.moves == 9:
            return True, "Draw"

        return True, "Ongoing"

    def update_board(self, symbol: Literal["X", "O"], x: int, y: int) -> bool:
        if x < 0 or x > 2 or y < 0 or y > 2:
            print("invalid coordinates")
            return False

        if self.board[x][y] != "":
            print("its already filled")
            return False

        self.board[x][y] = symbol
        return True

    def board_check(self, symbol: str, x: int, y: int) -> bool:
        lines = [
            (x, 0, x, 1, x, 2),
            (0, y, 1, y, 2, y),
            (0, 0, 1, 1, 2, 2),
            (0, 2, 1, 1, 2, 0),
        ]

        for line in lines:
            if self.board[line[0]][line[1]] == symbol and self.board[line[2]][line[3]] == symbol and self.board[line[4]][line[5]] == symbol:
                return True

        return False

game = TicTacToe()
assert game.execute_move("X", 0, 0) == (True, "Ongoing")
assert game.execute_move("O", 1, 1) == (True, "Ongoing")
assert game.execute_move("X", 1, 0) == (True, "Ongoing")
assert game.execute_move("O", 2, 0) == (True, "Ongoing")
assert game.execute_move("X", 0, 2) == (True, "Ongoing")
assert game.execute_move("O", 0, 1) == (True, "Ongoing")
assert game.execute_move("X", 0, 1) == (False, "Ongoing")
assert game.execute_move("X", 2, 2) == (True, "Ongoing")
assert game.execute_move("O", 2, 1) == (True, "O wins")

game = TicTacToe()
assert game.execute_move("X", 1, 1) == (True, "Ongoing")
assert game.execute_move("O", 0, 1) == (True, "Ongoing")
assert game.execute_move("X", 1, 0) == (True, "Ongoing")
assert game.execute_move("O", 1, 2) == (True, "Ongoing")
assert game.execute_move("X", 0, 0) == (True, "Ongoing")
assert game.execute_move("O", 2, 0) == (True, "Ongoing")
assert game.execute_move("X", 2, 2) == (True, "X wins")

game = TicTacToe()
assert game.execute_move("X", 0, 0) == (True, "Ongoing")
assert game.execute_move("O", 0, 1) == (True, "Ongoing")
assert game.execute_move("X", 0, 2) == (True, "Ongoing")
assert game.execute_move("O", 1, 1) == (True, "Ongoing")
assert game.execute_move("X", 1, 0) == (True, "Ongoing")
assert game.execute_move("O", 1, 2) == (True, "Ongoing")
assert game.execute_move("X", 2, 1) == (True, "Ongoing")
assert game.execute_move("O", 2, 0) == (True, "Ongoing")
assert game.execute_move("X", 2, 2) == (True, "Draw")

game = TicTacToe()
assert game.execute_move("X", -1, 0) == (False, "Ongoing")
assert game.execute_move("X", 0, -1) == (False, "Ongoing")
assert game.execute_move("X", 3, 0) == (False, "Ongoing")
assert game.execute_move("X", 0, 3) == (False, "Ongoing")
assert game.execute_move("X", 3, 3) == (False, "Ongoing")
