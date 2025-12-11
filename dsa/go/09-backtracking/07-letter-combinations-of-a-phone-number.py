# https://leetcode.com/problems/letter-combinations-of-a-phone-number/description/
# time complexity: O(4^n)
# space complexity: O(n)

from typing import List

class Solution:
    def letterCombinations(self, digits: str) -> List[str]:
        self.lenDigits = len(digits)

        if self.lenDigits == 0:
            return []
        
        self.res = []
        self.digits = digits
        self.dictPhone = {
            "2": ["a", "b", "c"],
            "3": ["d", "e", "f"],
            "4": ["g", "h", "i"],
            "5": ["j", "k", "l"],
            "6": ["m", "n", "o"],
            "7": ["p", "q", "r", "s"],
            "8": ["t", "u", "v"],
            "9": ["w", "x", "y", "z"],
        }

        self.backtracking([], 0)



        return self.res
    
    def backtracking(self, combination: List[int], digitsIdx: int):
        if digitsIdx == self.lenDigits:
            self.res.append("".join(combination))
            return
        
        for i in range(len(self.dictPhone[self.digits[digitsIdx]])):
            combination.append(self.dictPhone[self.digits[digitsIdx]][i])
            self.backtracking(combination, digitsIdx+1)
            combination.pop()
        
        return

if __name__ == "__main__":
    s = Solution()
    assert s.letterCombinations("23") == ["ad", "ae", "af", "bd", "be", "bf", "cd", "ce", "cf"]
    assert s.letterCombinations("") == []
    assert s.letterCombinations("2") == ["a", "b", "c"]