# https://leetcode.com/problems/kth-smallest-element-in-a-bst/description/

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
    # time complexity: O(n)
    # space complexity: O(n) for the array
    def kthSmallest(self, root: Optional[TreeNode], k: int) -> int:
        self.treeValues = []
        self.dfs(root)
        return self.treeValues[k-1]
    
    def dfs(self, root: Optional[TreeNode]):
        if root is None:
            return
        
        self.dfs(root.left)
        self.treeValues.append(root.val)
        self.dfs(root.right)
        return
    
    # time complexity: O(n)
    # space complexity: No array and only recursion stack → O(h)
    def kthSmallestOptimal(self, root: Optional[TreeNode], k: int) -> int:
        self.res = 0
        self.k = k
        self.dfsOptimal(root)
        return self.res
    
    def dfsOptimal(self, root: Optional[TreeNode]):
        if root is None:
            return
        
        self.dfsOptimal(root.left)

        self.k -= 1
        if self.k == 0:
            self.res = root.val
        
        self.dfsOptimal(root.right)
        return

if __name__ == "__main__":
    solution = Solution()
    root = TreeNode.fromList([3,1,4,None,2])
    assert solution.kthSmallest(root, 1) == 1

    root = TreeNode.fromList([5,3,6,2,4,None, None,1])
    assert solution.kthSmallest(root, 3) == 3


