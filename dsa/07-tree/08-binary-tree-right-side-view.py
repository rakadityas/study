# https://leetcode.com/problems/binary-tree-right-side-view/description/

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
    # Time: O(n), Space: O(n)
    def rightSideViewQueue(self, root: Optional[TreeNode]) -> List[int]:
        self.listTree = [] # creates a queue of each tree depth
        self.dfsRightSideViewQueue(root, 0)

        res = []
        for i in range(len(self.listTree)):
            if len(self.listTree) > 0:
                res.append(self.listTree[i][0])
        
        return res
    
    def dfsRightSideViewQueue(self, root: Optional[TreeNode], depth: int):
        if root is None:
            return 
        
        if depth >= len(self.listTree):
            self.listTree.append([])
        
        self.listTree[depth].append(root.val)

        self.dfsRightSideViewQueue(root.right, depth+1)
        self.dfsRightSideViewQueue(root.left, depth+1)

        return
    
    # Time: O(n), Space: O(n)
    def rightSideViewList(self, root: Optional[TreeNode]) -> List[int]:
        return self.traverseTreeList(root, [], 0)

    def traverseTreeList(self, node: Optional[TreeNode], resp: List[int], depth: int) -> List[int]:
        if node is None:
            return resp
        
        if depth == len(resp):
            resp.append(node.val)
        
        resp = self.traverseTreeList(node.right, resp, depth+1)
        resp = self.traverseTreeList(node.left, resp, depth+1)

        return resp

if __name__ == "__main__":
    solution = Solution()
    
    root1 = TreeNode.fromList([1,2,3,None,5,None,4])
    resp = solution.rightSideViewQueue(root1)
    assert resp == [1,3,4]

    root2 = TreeNode.fromList([])
    resp = solution.rightSideViewQueue(root2)
    assert resp == []

    root3 = TreeNode.fromList([1,2,3,4,None,None,None,5])
    resp = solution.rightSideViewQueue(root3)
    assert resp == [1,3,4,5]

    root4 = TreeNode.fromList([1,None,3])
    resp = solution.rightSideViewQueue(root4)
    assert resp == [1,3]

    #####################

    root1 = TreeNode.fromList([1,2,3,None,5,None,4])
    resp = solution.rightSideViewList(root1)
    assert resp == [1,3,4]

    root2 = TreeNode.fromList([])
    resp = solution.rightSideViewList(root2)
    assert resp == []

    root3 = TreeNode.fromList([1,2,3,4,None,None,None,5])
    resp = solution.rightSideViewList(root3)
    assert resp == [1,3,4,5]

    root4 = TreeNode.fromList([1,None,3])
    resp = solution.rightSideViewList(root4)
    assert resp == [1,3]