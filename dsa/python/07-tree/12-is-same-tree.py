# https://leetcode.com/problems/same-tree/description/
# time complexity: O(n)
# space complexity: O(n) stack space

from typing import List, Optional

class TreeNode:
    def __init__(self, val: int = 0, left: Optional["TreeNode"] = None, right: Optional["TreeNode"] = None):
        self.val = val
        self.left = left
        self.right = right

    @classmethod
    def fromList(cls, values: List[int]) -> Optional["TreeNode"]:
        if not values:
            return None

        root = cls(values[0])
        queue = [root]
        i = 1
        while queue and i < len(values):
            node = queue.pop(0)
            if i < len(values) and values[i] is not None:
                node.left = cls(values[i])
                queue.append(node.left)
            i += 1
            if i < len(values) and values[i] is not None:
                node.right = cls(values[i])
                queue.append(node.right)
            i += 1
        return root

    def toList(self) -> List[int]:
        result = []
        queue = [self]
        while queue:
            node = queue.pop(0)
            if node:
                result.append(node.val)
                queue.append(node.left)
                queue.append(node.right)
            else:
                result.append(None)
        while result and result[-1] is None:
            result.pop()
        return result

class Solution:
    def isSameTree(self, p: Optional[TreeNode], q: Optional[TreeNode]) -> bool:
        if p is None and q is None:
            return True
        elif p is None or q is None:
            return False
        
        if p.val != q.val:
            return False
        
        isSameLeft = self.isSameTree(p.left, q.left)
        if not isSameLeft:
            return False

        isSameRight = self.isSameTree(p.right, q.right)
        if not isSameRight:
            return False
        
        return True

if "__name__" == "__main__":
    solution = Solution()

    assert solution.isSameTree(TreeNode.fromList([1,2,3]), TreeNode.fromList([1,2,3])) == True
    assert solution.isSameTree(TreeNode.fromList([1,2]), TreeNode.fromList([1,None,2])) == False
    assert solution.isSameTree(TreeNode.fromList([1,2,1]), TreeNode.fromList([1,1,2])) == False
