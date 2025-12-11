# https://leetcode.com/problems/subtree-of-another-tree/description/
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
    def isSubtree(self, root: Optional[TreeNode], subRoot: Optional[TreeNode]) -> bool:
        return self.dfs(root, subRoot)
    
    def dfs(self, root: Optional[TreeNode], subRoot: Optional[TreeNode]) -> bool:
        if root is None:
            return False
        
        if self.checkSubroot(root, subRoot):
            return True
        
        return self.dfs(root.left, subRoot) or self.dfs(root.right, subRoot)
    
    def checkSubroot(self, root: Optional[TreeNode], subRoot: Optional[TreeNode]) -> bool:
        if root is None and subRoot is None:
            return True
        
        if root is None or subRoot is None:
            return False
        
        if root.val != subRoot.val:
            return False
        
        return self.checkSubroot(root.left, subRoot.left) and self.checkSubroot(root.right, subRoot.right)

if __name__ == "__main__":
    solution = Solution()

    root = TreeNode.fromList([3,4,5,1,2])
    subRoot = TreeNode.fromList([4,1,2])
    assert solution.isSubtree(root, subRoot) == True

    root = TreeNode.fromList([3,4,5,1,2,None,None,None,None,0])
    subRoot = TreeNode.fromList([4,1,2])
    assert solution.isSubtree(root, subRoot) == False