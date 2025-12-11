# https://leetcode.com/problems/remove-nth-node-from-end-of-list/description/
# time complexity: O(n)
# space complexity: O(n)

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
    
    def removeNthFromEnd(self, head: Optional[ListNode], n: int) -> Optional[ListNode]:
        nodeArr = []

        curr = head
        while curr:
            nodeArr.append(curr)
            curr = curr.next
        
        idx = len(nodeArr) - n

        if idx == 0: # no idx-1
            return head.next

        nodeArr[idx-1].next = nodeArr[idx].next
        
        return head

if __name__ == "__main__":
    solution = Solution()
    caseOne = ListNode.fromList([1,2,3,4,5])
    caseOne = solution.removeNthFromEnd(caseOne, 2)
    assert caseOne.toList() == [1,2,3,5]

    caseTwo = ListNode.fromList([1])
    caseTwo = solution.removeNthFromEnd(caseTwo, 1)
    assert (caseTwo.toList() if caseTwo else []) == []

    caseThree = ListNode.fromList([1,2])
    caseThree = solution.removeNthFromEnd(caseThree, 1)
    assert caseThree.toList() == [1]