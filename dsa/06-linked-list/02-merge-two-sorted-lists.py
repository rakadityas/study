# https://leetcode.com/problems/merge-two-sorted-lists/description/

from typing import Optional

class ListNode:
    def __init__(self, val=0, next=None):
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
    def mergeTwoLists(self, list1: Optional[ListNode], list2: Optional[ListNode]) -> Optional[ListNode]:
        
        respList = ListNode()
        curr = respList

        while list1 != None and list2 != None:
            if list1.val < list2.val:
                curr.next = list1
                list1 = list1.next
            else:
                curr.next = list2
                list2 = list2.next
                
            curr = curr.next
        
        if list1 != None:
            curr.next = list1
        elif list2 != None:
            curr.next = list2
        
        return respList.next

if __name__ == "__main__":
    solution = Solution()
    caseOneList = solution.mergeTwoLists(ListNode.fromList([1,2,4]), ListNode.fromList([1,3,4]))
    assert caseOneList.toList() == [1,1,2,3,4,4]

    caseTwoList = solution.mergeTwoLists(ListNode.fromList([1,2,4]), ListNode.fromList([]))
    assert caseTwoList.toList() == [1,2,4]

    caseThreeList = solution.mergeTwoLists(ListNode.fromList([]), ListNode.fromList([]))
    assert caseThreeList == None
    



