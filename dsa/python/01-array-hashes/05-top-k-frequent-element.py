# https://leetcode.com/problems/top-k-frequent-elements/description/
# time complexity: O(n)
# space complexity: O(n)

from collections import defaultdict

class Solution:
    # Approach: count frequencies, then repeatedly call max() to extract top-k groups
    # time complexity: O(n·k) worst case — max() over distinct keys runs up to k times
    # space complexity: O(n)
    def topKFrequent(self, nums: list[int], k: int) -> list[int]:

        mapNumsHistory = defaultdict(int)
        for i in range(len(nums)):
            mapNumsHistory[nums[i]] += 1
        
        mapGroupNums = defaultdict(list)
        for key, value in mapNumsHistory.items():
            mapGroupNums[value].append(key)

        res = []
        while k > 0:
            maxKey = max(mapGroupNums.keys())
            for val in mapGroupNums[maxKey]:
                res.append(val)
                k -= 1
                if k == 0:
                    return res
            
            del mapGroupNums[maxKey]
        
        return res
    
    # Approach: bucket sort — index bucket[freq] holds all numbers with that frequency, then scan from high to low
    # time complexity: O(n) — no comparison-based sort; frequency range is bounded by n
    # space complexity: O(n)
    def topKFrequentBucketSort(self, nums: list[int], k: int) -> list[int]:
        lenNums = len(nums)
        mapCounter = defaultdict(int)
        for i in range(lenNums):
            mapCounter[nums[i]] += 1
        
        bucket = []
        for i in range(lenNums+1):
            bucket.append([])
                
        for key, val in mapCounter.items():
            bucket[val].append(key)
        
        res = []

        for freq in range(lenNums, 0, -1):
            if bucket[freq]:
                for i in range(len(bucket[freq])):
                    res.append(bucket[freq][i])
                    if len(res) == k:
                        return res
        return res


if __name__ == "__main__":
    solution = Solution()

    assert solution.topKFrequent([1,1,1,2,2,3], 2) == [1,2]
    assert solution.topKFrequent([1], 1) == [1]

    assert solution.topKFrequentBucketSort([1,1,1,2,2,3], 2) == [1,2]
    assert solution.topKFrequentBucketSort([1], 1) == [1]
