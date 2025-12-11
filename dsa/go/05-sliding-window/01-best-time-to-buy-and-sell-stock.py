# https://leetcode.com/problems/best-time-to-buy-and-sell-stock/description/
# time complexity: O(n)
# space complexity: O(1)

class Solution:
    def maxProfit(self, prices: list[int]) -> int:
        l, r = 0, 1
        maxProfit = 0

        while r < len(prices):
            if prices[l] < prices[r]:
                maxProfit = max(prices[r] - prices[l], maxProfit)
            else:
                l = r
            r += 1
        
        return maxProfit

if __name__ == "__main__":
    solution = Solution()

    assert solution.maxProfit([7,1,5,3,6,4]) == 5
    assert solution.maxProfit([7,6,4,3,1]) == 0
