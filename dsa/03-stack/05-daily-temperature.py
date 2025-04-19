class Solution:
    def dailyTemperatures(self, temperatures: list[int]) -> list[int]:
        res = [0] * len(temperatures)
        stack = []

        for i in range(len(temperatures)):
            while stack and temperatures[i] > temperatures[stack[-1]]:
                prevIdx = stack.pop()
                res[prevIdx] = i - prevIdx
            
            stack.append(i)
        
        return res
    
if __name__ == "__main__":
    solution = Solution()

    assert solution.dailyTemperatures([73,74,75,71,69,72,76,73]) == [1,1,4,2,1,1,0,0]
    assert solution.dailyTemperatures([30,40,50,60]) == [1,1,1,0]
    assert solution.dailyTemperatures([30,60,90]) == [1,1,0]