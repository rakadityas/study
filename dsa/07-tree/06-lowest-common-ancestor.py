# https://leetcode.com/problems/lowest-common-ancestor-of-a-binary-search-tree/
from typing import Optional

class TreeNode:
    def __init__(self, val: int = 0, left: Optional["TreeNode"] = None, right: Optional["TreeNode"] = None):
        self.val = val
        self.left = left
        self.right = right
    
    @classmethod
    def fromList(cls, values: list[int]) -> Optional["TreeNode"]:
        if not values or values[0] is None:
            return None
            
        root = cls(values[0])
        queue = [root]
        i = 1
        
        while queue and i < len(values):
            node = queue.pop(0)
            
            # Left child
            if i < len(values) and values[i] is not None:
                node.left = cls(values[i])
                queue.append(node.left)
            i += 1
            
            # Right child
            if i < len(values) and values[i] is not None:
                node.right = cls(values[i])
                queue.append(node.right)
            i += 1
            
        return root
    
    def toList(self) -> list[int]:
        if not self:
            return []
            
        result = []
        queue = [self]
        
        while queue:
            node = queue.pop(0)
            result.append(node.val)
            
            if node.left:
                queue.append(node.left)
            if node.right:
                queue.append(node.right)
                
        return result
        
class Solution:
    def lowestCommonAncestor(self, root: "TreeNode", p: "TreeNode", q: "TreeNode") -> "TreeNode":
        curr = root
        while curr:
            if p.val < curr.val and q.val < curr.val:
                curr = curr.left
            elif p.val > curr.val and q.val > curr.val:
                curr = curr.right
            else:
                return curr

if __name__ == "__main__":
    solution = Solution()

    # case one
    resp = solution.lowestCommonAncestor(TreeNode.fromList([6,2,8,0,4,7,9,None,None,3,5]), 
    TreeNode(2, None, None), 
    TreeNode(8, None, None))
    assert resp.val == 6

    # case two
    resp = solution.lowestCommonAncestor(TreeNode.fromList([6,2,8,0,4,7,9,None,None,3,5]), 
    TreeNode(2, None, None), 
    TreeNode(4, None, None))
    assert resp.val == 2

    # case three
    resp = solution.lowestCommonAncestor(TreeNode.fromList([2,1]), 
    TreeNode(2, None, None), 
    TreeNode(1, None, None))
    assert resp.val == 2