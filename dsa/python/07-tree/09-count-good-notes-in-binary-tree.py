# https://leetcode.com/problems/count-good-nodes-in-binary-tree/description/
# Time complexity is O(N) because each node is visited once.
# Space complexity is O(H) due to recursion stack, where H is the tree height (worst case O(N)).

from typing import Optional, List

class TreeNode:
    def __init__(self, val=0, left: Optional["TreeNode"] = None, right: Optional["TreeNode"] = None):
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
    def goodNodes(self, root: TreeNode) -> int:
        self.nodesCount = 0
        self.dfs(root, root.val)
        return self.nodesCount
    
    def dfs(self, root: TreeNode, maxVal: int) -> int:
        if root is None:
            return
        
        if root.val >= maxVal:
            self.nodesCount += 1
            maxVal = root.val
        
        self.dfs(root.left, maxVal)
        self.dfs(root.right, maxVal)
        
        return

if __name__ == "__main__":
    solution = Solution()
    
    root = TreeNode.fromList([3,1,4,3,None,1,5])
    assert solution.goodNodes(root) == 4

    root = TreeNode.fromList([3,3,None,4,2])
    assert solution.goodNodes(root) == 3

    root = TreeNode.fromList([1])
    assert solution.goodNodes(root) == 1

    root = TreeNode.fromList([-1,5,-2,4,4,2,-2,None,None,-4,None,-2,3,None,-2,0,None,-1,None,-3,None,-4,-3,3,None,None,None,None,None,None,None,3,-3])
    assert solution.goodNodes(root) == 5

    



