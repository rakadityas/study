# https://leetcode.com/problems/add-two-numbers/description/
# time complexity: O(n)
# space complexity: O(n)

from typing import Optional

class ListNode:
    def __init__(self, val: int = 0, next: Optional["ListNode"] = None):
        self.val = val
        self.next = next

    @classmethod
    def fromList(cls, values: list[int]) -> Optional["ListNode"]:
        if len(values) == 0:
            return None

        dummy = cls()
        current = dummy
        for i in range(len(values)):
            current.next = cls(values[i])
            current = current.next

        return dummy.next

    def toList(self) -> list[int]:
        result = []
        current = self
        while current is not None:
            result.append(current.val)
            current = current.next
        return result

    def printNode(self):
        values = self.toList()
        print(" -> ".join(map(str, values)) + " -> None")


class Solution:
    def addTwoNumbers(self, l1: Optional[ListNode], l2: Optional[ListNode]) -> Optional[ListNode]:
        res = ListNode()
        curr = res
        remainder = 0
        while l1 is not None or l2 is not None or remainder > 0:
            l1Val = 0
            if l1 is not None:
                l1Val = l1.val
                l1 = l1.next
            
            l2Val = 0
            if l2 is not None:
                l2Val = l2.val
                l2 = l2.next

            sumAmt = l1Val + l2Val + remainder
            remainder = sumAmt // 10
            valAmt = sumAmt % 10
            curr.next = ListNode(valAmt)
            curr = curr.next
        return res.next

if __name__ == "__main__":
    solution = Solution()

    caseOneA = ListNode.fromList([2, 4, 3])
    caseOneB = ListNode.fromList([5, 6, 4])
    resultOne = solution.addTwoNumbers(caseOneA, caseOneB)
    assert (resultOne.toList() if resultOne else []) == [7, 0, 8]

    caseTwoA = ListNode.fromList([0])
    caseTwoB = ListNode.fromList([0])
    resultTwo = solution.addTwoNumbers(caseTwoA, caseTwoB)
    assert (resultTwo.toList() if resultTwo else []) == [0]

    caseThreeA = ListNode.fromList([9, 9, 9, 9, 9, 9, 9])
    caseThreeB = ListNode.fromList([9, 9, 9, 9])
    resultThree = solution.addTwoNumbers(caseThreeA, caseThreeB)
    assert (resultThree.toList() if resultThree else []) == [8, 9, 9, 9, 0, 0, 0, 1]