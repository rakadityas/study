# time complexity: O(n) - each character in the formula is visited exactly once
# space complexity: O(d) stack entries, where d = max nesting depth (d <= n)

# it is known that
# C = 12
# H = 1

ATOMIC_WEIGHTS = {
    "C": 12,
    "H": 1,
}


class Solution:
    def countAtoms(self, formula: str) -> int:
        stack = [0]  # running weight total at each nesting level
        i, n = 0, len(formula)

        while i < n:
            char = formula[i]

            if char == "(":
                stack.append(0)
                i += 1

            elif char == ")":
                count, i = self._read_count(formula, i + 1)
                group_weight = stack.pop() * count
                stack[-1] += group_weight

            else:
                weight = ATOMIC_WEIGHTS[char]
                count, i = self._read_count(formula, i + 1)
                stack[-1] += weight * count

        return stack[0]

    @staticmethod
    def _read_count(formula, i):
        start = i
        while i < len(formula) and formula[i].isdigit():
            i += 1
        count = int(formula[start:i]) if i > start else 1
        return count, i


if __name__ == "__main__":
    solution = Solution()
    assert solution.countAtoms("CH4") == 16
    assert solution.countAtoms("C(H4)H4") == 20
    assert solution.countAtoms("H10") == 10  # multi-digit count, was broken before
