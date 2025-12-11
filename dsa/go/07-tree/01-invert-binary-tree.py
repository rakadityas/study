# https://leetcode.com/problems/invert-binary-tree/description/
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
    def invertTree(self, root: Optional[TreeNode]) -> Optional[TreeNode]:
        if root is None:
            return None
        
        temp = root.left
        root.left = root.right
        root.right = temp

        self.invertTree(root.left)
        self.invertTree(root.right)
        return root

if __name__ == "__main__":
    solution = Solution()

    # Case 1: Regular binary tree
    caseOne = TreeNode.fromList([4,2,7,1,3,6,9])
    assert solution.invertTree(caseOne).toList() == [4,7,2,9,6,3,1]

    # Case 2: Small tree
    caseTwo = TreeNode.fromList([2,1,3])
    assert solution.invertTree(caseTwo).toList() == [2,3,1]

    # Case 3: Empty tree
    caseThree = TreeNode.fromList([])
    assert solution.invertTree(caseThree) == None

    # Case 4: Single node
    caseFour = TreeNode.fromList([1])
    assert solution.invertTree(caseFour).toList() == [1]