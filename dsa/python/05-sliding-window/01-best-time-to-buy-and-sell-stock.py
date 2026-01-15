# https://leetcode.com/problems/best-time-to-buy-and-sell-stock/description/
# time complexity: O(n)
# space complexity: O(1)

class Solution:
    def maxProfit(self, prices: List[int]) -> int:
        l = 0
        maxAmt = 0
        for i in range(len(prices)):
            if prices[l] <= prices[i]:
                maxAmt = max(maxAmt, prices[i]-prices[l])
            else:
                l = i
        return maxAmt

if __name__ == "__main__":
    solution = Solution()

    assert solution.maxProfit([7,1,5,3,6,4]) == 5
    assert solution.maxProfit([7,6,4,3,1]) == 0
