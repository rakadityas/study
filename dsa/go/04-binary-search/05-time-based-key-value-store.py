# https://leetcode.com/problems/time-based-key-value-store/description/
# time complexity: O(log(n))
# space complexity: O(n)

from collections import defaultdict

class TimeMap:
    def __init__(self):
        self.timeMap = defaultdict(list)

    def set(self, key: str, value: str, timestamp: int) -> None:
        if not self.timeMap[key]:
            self.timeMap[key] = []

        self.timeMap[key].append([value, timestamp])
        
    def get(self, key: str, timestamp: int) -> str:
        if key not in self.timeMap:
            return ""
        
        timeMapArr = self.timeMap[key]
        l, r = 0, len(timeMapArr) - 1

        latestData = timeMapArr[-1]

        while l <= r:
            mid = (l+r)//2
            midValue, midTimestamp = timeMapArr[mid]
            if midTimestamp == timestamp:
                return midValue
            elif timeMapArr[l][1] > midTimestamp:
                l = mid+1
            else:
                r = mid-1
            
        return latestData[0]


# Your TimeMap object will be instantiated and called as such:
# obj = TimeMap()
# obj.set(key,value,timestamp)
# param_2 = obj.get(key,timestamp)

if __name__ == "__main__":
    timeMap = TimeMap()
    
    timeMap.set("foo", "bar", 1)
    assert timeMap.get("foo", 1) == "bar"
    assert timeMap.get("foo", 3) == "bar"
    timeMap.set("foo", "bar2", 4)
    assert timeMap.get("foo", 4) == "bar2"
    assert timeMap.get("foo", 5) == "bar2"
    