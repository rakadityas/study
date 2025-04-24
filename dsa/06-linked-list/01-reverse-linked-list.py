from typing import Optional

class ListNode:
    def __init__(self, val: int = 0, next: Optional["ListNode"] = None):
        self.val = val
        self.next = next

    @classmethod
    def fromList(cls, values: list[int]) -> Optional["ListNode"]:
        if not values:
            return None

        head = cls(values[0])
        current = head
        for val in values[1:]:
            current.next = cls(val)
            current = current.next
        return head

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
    def reverseList(self, head: Optional[ListNode]) -> Optional[ListNode]:
        prev = None
        curr = head
        while curr:
            next_temp = curr.next
            curr.next = prev
            prev = curr
            curr = next_temp
        return prev
    
if __name__ == "__main__":
    # case one
    originalCase = [1, 2, 3, 4, 5]
    expectedReversedCase = [5, 4, 3, 2, 1]

    head = ListNode.fromList(originalCase)
    solution = Solution()
    reversedHead = solution.reverseList(head)

    assert reversedHead.toList() == expectedReversedCase

    # case two
    originalCase = [1, 2]
    expectedReversedCase = [2, 1]

    head = ListNode.fromList(originalCase)
    solution = Solution()
    reversedHead = solution.reverseList(head)

    assert reversedHead.toList() == expectedReversedCase

