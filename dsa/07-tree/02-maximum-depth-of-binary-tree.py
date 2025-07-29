# https://leetcode.com/problems/maximum-depth-of-binary-tree/

from typing import Optional, List

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
    def maxDepth(self, root: Optional[TreeNode]) -> int:
        return self.dfs(root, 0)
        
    def dfs(self, root: Optional[TreeNode], counter: int) -> int:
        if root is None:
            return counter
        
        maxCounter = max(self.dfs(root.left, counter+1), self.dfs(root.right, counter+1))
        return maxCounter

if __name__ == "__main__":
    solution = Solution()
    # Case 1: Regular binary tree
    root1 = TreeNode.fromList([3,9,20,None, None,15,7])
    assert solution.maxDepth(root1) == 3

    # Case 2: Empty tree
    root2 = TreeNode.fromList([1,None,2]) 
    assert solution.maxDepth(root2) == 2