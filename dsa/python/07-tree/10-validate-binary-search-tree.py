# https://leetcode.com/problems/validate-binary-search-tree/description/
# time complexity: O(n)
# space complexity: O(n)

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
    def isValidBST(self, root: Optional[TreeNode]) -> bool:
        return self.dfs(root, float('-inf'), float('inf'))
    
    def dfs(self, root: Optional[TreeNode], min_val: float, max_val: float) -> bool:
        if root is None:
            return True
        
        if root.val <= min_val or root.val >= max_val:
            return False
        
        return (self.dfs(root.left, min_val, root.val) and 
                self.dfs(root.right, root.val, max_val)) 

if __name__ == "__main__":
    solution = Solution()
    
    root1 = TreeNode.fromList([2, 1, 3])
    result1 = solution.isValidBST(root1)
    assert result1 == True  
    
 
    root2 = TreeNode.fromList([5, 1, 4, None, None, 3, 6])
    result2 = solution.isValidBST(root2)
    assert result2 == False
    
    root3 = TreeNode.fromList([1])
    result3 = solution.isValidBST(root3)
    assert result3 == True
    
    root4 = TreeNode.fromList([2, 2, 2])
    result4 = solution.isValidBST(root4)
    assert result4 == False
