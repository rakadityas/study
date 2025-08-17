# https://leetcode.com/problems/binary-tree-level-order-traversal/

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
    def levelOrder(self, root: Optional[TreeNode]) -> List[List[int]]:
        if not root:
            return []
        
        self.res = []
        self.dfs(root, 0)
        return self.res

    def dfs(self, root: Optional[TreeNode], depth: int):
        if root is None:
            return
        
        # Ensure the result list has enough sublists for the current depth
        if depth >= len(self.res):
            self.res.append([])
        
        self.res[depth].append(root.val)
        
        self.dfs(root.left, depth + 1)
        self.dfs(root.right, depth + 1)


if __name__ == "__main__":
    solution = Solution()

    root1 = TreeNode.fromList([3, 9, 20, None, None, 15, 7])
    result1 = solution.levelOrder(root1)
    assert result1 == [[3], [9, 20], [15, 7]]
    
    root2 = TreeNode.fromList([1])
    result2 = solution.levelOrder(root2)
    assert result2 == [[1]]
    
    root3 = TreeNode.fromList([])
    result3 = solution.levelOrder(root3)
    assert result3 == []