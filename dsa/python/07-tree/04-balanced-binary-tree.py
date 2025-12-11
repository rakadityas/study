# https://leetcode.com/problems/balanced-binary-tree/description/
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
    def isBalanced(self, root: Optional[TreeNode]) -> bool:
        resp = self.dfs(root)
        return resp[1]
    
    def dfs(self, root: Optional[TreeNode]) -> list:
        if root is None:
            return 0, True
        
        leftChild, isValid = self.dfs(root.left)
        if isValid == False:
            return [0, False]

        rightChild, isValid = self.dfs(root.right)
        if isValid == False:
            return [0, False]

        if abs(leftChild-rightChild) > 1: # means not balance
            return [0, False]
        
        return [max(leftChild+1, rightChild+1), True]

if __name__ == "__main__":
    solution = Solution()

    root = TreeNode.fromList([3,9,20,None,None,15,7])
    assert solution.isBalanced(root) == True

    root2 = TreeNode.fromList([1,2,2,3,3,None,None,4,4])
    assert solution.isBalanced(root2) == False

    root3 = TreeNode.fromList([])
    assert solution.isBalanced(root3) == True
