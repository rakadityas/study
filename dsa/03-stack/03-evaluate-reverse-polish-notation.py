
class Solution:
    def evalRPN(self, tokens: list[str]) -> int:
        stack = []
        for i in range(len(tokens)):
            if tokens[i] == "+":
                val = stack.pop() + stack.pop()
                stack.append(val)
            elif tokens[i] == "*":
                val = stack.pop()* stack.pop()
                stack.append(val)
            elif tokens[i] == "-":
                valA = stack.pop()
                valB = stack.pop()
                res = valB-valA
                stack.append(res)
            elif tokens[i] == "/":
                valA = stack.pop()
                valB = stack.pop()
                res = valB/valA
                stack.append(int(res))
            else:
                stack.append(int(tokens[i]))
        
        return stack.pop()

if __name__ == "__main__":
    solution = Solution()
    
    assert solution.evalRPN(["2","1","+","3","*"]) == 9
    assert solution.evalRPN(["10","6","9","3","+","-11","*","/","*","17","+","5","+"]) == 22
    assert solution.evalRPN(["4","13","5","/","+"]) == 6

