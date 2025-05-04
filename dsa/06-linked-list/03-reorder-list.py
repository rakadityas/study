# https://leetcode.com/problems/reorder-list/description/

from typing import Optional

class ListNode:
    def __init__(self, val: int = 0, next: Optional["ListNode"] = None):
        self.val = val
        self.next = next

    @classmethod
    def fromList(cls, values: list[int]) -> Optional["ListNode"]:
        if not values:
            return None

        head = cls()
        current = head
        for val in range(len(values)):
            current.next = cls(values[val])
            current = current.next
        return head.next

    def toList(self) -> list[int]:
        result = []
        curr = self
        while curr:
            result.append(curr.val)
            curr = curr.next
        return result

    def printNode(self):
        print(" -> ".join(map(str, self.toList())) + " -> None")

class Solution:
    def reorderList(self, head: Optional[ListNode]) -> None:
        nodesArr = []
        curr = head
        while curr:
            nodesArr.append(curr)
            curr = curr.next
        
        left, right = 0, len(nodesArr)-1
        while left < right:
            nodesArr[left].next = nodesArr[right]
            left += 1

            if left >= right:
                break

            nodesArr[right].next = nodesArr[left]
            right -= 1
        
        # crucial to stop the list from cyclic
        nodesArr[left].next = None

        return

if __name__ == "__main__":
    solution = Solution()
    caseOne = ListNode.fromList([1,2,3,4])
    solution.reorderList(caseOne)
    assert caseOne.toList() == [1,4,2,3]

    caseTwo = ListNode.fromList([1,2,3,4,5])
    solution.reorderList(caseTwo)
    assert caseTwo.toList() == [1,5,2,4,3]