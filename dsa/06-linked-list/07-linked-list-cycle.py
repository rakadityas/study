# https://leetcode.com/problems/linked-list-cycle/

from typing import Optional

class ListNode:
    def __init__(self, val: int = 0, next: Optional["ListNode"] = None):
        self.val = val
        self.next = next

    @classmethod
    def fromList(cls, values: list[int], pos: int = -1) -> Optional["ListNode"]:
        if len(values) == 0:
            return None

        nodes = []
        for val in values:
            nodes.append(cls(val))

        for i in range(len(nodes) - 1):
            nodes[i].next = nodes[i + 1]

        if pos != -1:
            nodes[-1].next = nodes[pos]

        return nodes[0]

    def toList(self, maxNodes: int = 100) -> list[int]:
        result = []
        current = self
        count = 0
        visited = set()
        while current is not None and count < maxNodes:
            if id(current) in visited:
                result.append(f"(cycle to {current.val})")
                break
            visited.add(id(current))
            result.append(current.val)
            current = current.next
            count += 1
        return result


class Solution:
    def hasCycle(self, head: Optional[ListNode]) -> bool:
        curr = head
        mapHistory = {}
        while curr is not None:
            if curr in mapHistory:
                return True
            mapHistory[curr] = True
            curr = curr.next
        
        return False

if __name__ == "__main__":
    solution = Solution()

    # Case 1: cycle at position 1
    caseOne = ListNode.fromList([3, 2, 0, -4], pos=1)
    assert solution.hasCycle(caseOne) == True

    # Case 2: cycle at position 0
    caseTwo = ListNode.fromList([1, 2], pos=0)
    assert solution.hasCycle(caseTwo) == True

    # Case 3: no cycle
    caseThree = ListNode.fromList([1], pos=-1)
    assert solution.hasCycle(caseThree) == False

    # Case 4: empty list
    caseFour = ListNode.fromList([], pos=-1)
    assert solution.hasCycle(caseFour) == False