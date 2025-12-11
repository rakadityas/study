# https://leetcode.com/problems/diameter-of-binary-tree/description/
# time complexity: O(n)
# space complexity: O(n)

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
    def diameterOfBinaryTree(self, root: Optional[TreeNode]) -> int:
        self.maxDiameter = 0
        self.dfs(root)
        return self.maxDiameter
    
    def dfs(self, root: Optional[TreeNode]) -> int:
        if root is None:
            return 0
            
        leftVal = self.dfs(root.left)
        rightVal = self.dfs(root.right)

        self.maxDiameter = max(leftVal+rightVal, self.maxDiameter)

        return 1 + max(leftVal, rightVal) # just need to pick one of it.

if __name__ == "__main__":
    solution = Solution()
    root = TreeNode.fromList([1,2,3,4,5])
    assert solution.diameterOfBinaryTree(root) == 3

    root = TreeNode.fromList([1,2])
    assert solution.diameterOfBinaryTree(root) == 1

    root = TreeNode.fromList([1,2,3,4,5,6,7])
    assert solution.diameterOfBinaryTree(root) == 4