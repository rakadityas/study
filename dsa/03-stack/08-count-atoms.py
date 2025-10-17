# it is known that
# C = 12
# H = 1

from typing import List

class Solution:
    def countAtoms(self, formula: str) -> int:
        mapFormula = {
            "C": 12,
            "H": 1,
        }
        openBracket = {
            "(": True,
        }
        closeBracket = {
            ")": True,
        }

        stack = []
        
        for i in range(len(formula)):
            if formula[i] in mapFormula:
                stack.append(mapFormula[formula[i]])
            
            elif formula[i] in openBracket:
                stack.append(formula[i])

            elif formula[i] in closeBracket:
                hasFound = False
            
                while hasFound is False:
                    prev = stack.pop()
                    prevPrevVal = stack.pop()

                    if prevPrevVal in openBracket:
                        hasFound = True
                        stack.append(prev)
                    else:
                        stack.append(prev+prevPrevVal)
            else:
                prevVal = stack.pop()
                stack.append(prevVal * int(formula[i]))

        while len(stack) > 1:
            stack.append(stack.pop() + stack.pop())
        
        return stack.pop()

if __name__ == "__main__":
    solution = Solution()
    assert solution.countAtoms("CH4") == 16
    assert solution.countAtoms("C(H4)H4") == 20