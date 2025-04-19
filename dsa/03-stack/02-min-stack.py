
class MinStack:
    def __init__(self):
        self.stack = []
        self.minStack = []
    
    def push(self, val: int) -> None:
        self.stack.append(val)
        lowest = val
        if self.minStack and self.minStack[-1] < val:
            lowest = self.minStack[-1]
        self.minStack.append(lowest)
            
    def pop(self) -> None:
        self.stack.pop()
        self.minStack.pop()

    def top(self) -> int:
        return self.stack[-1]

    def getMin(self) -> int:
        return self.minStack[-1]

if __name__ == "__main__":
    exampleOne = MinStack()
    exampleOne.push(-2)
    exampleOne.push(0)
    exampleOne.push(-3)
    assert exampleOne.getMin() == -3
    exampleOne.pop()
    assert exampleOne.top() == 0
    assert exampleOne.getMin() == -2
    