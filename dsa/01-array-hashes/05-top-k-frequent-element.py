from collections import defaultdict

class Solution:
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

if __name__ == "__main__":
    solution = Solution()

    assert solution.topKFrequent([1,1,1,2,2,3], 2) == [1,2]
    assert solution.topKFrequent([1], 1) == [1]
