# https://leetcode.com/problems/clone-graph/description/
from typing import List, Optional

class Node:
    def __init__(self, val = 0, neighbors = None):
        self.val = val
        self.neighbors = neighbors if neighbors is not None else []
    
    def fromList(self, nodesArr: List[List[int]]) -> 'Node':
        if not nodesArr:
            return None
        
        nodes = [Node(i+1) for i in range(len(nodesArr))]
        
        for i, node in enumerate(nodes):
            for neighbor in nodesArr[i]:
                node.neighbors.append(nodes[neighbor-1])
        
        return nodes[0]
    
    def toList(self) -> List[List[int]]:
        if not self:
            return []
        
        result = {}
        queue = [self]
        visited = set()
        
        while queue:
            node = queue.pop(0)
            if node in visited:
                continue
            
            visited.add(node)
            result[node.val] = [neighbor.val for neighbor in node.neighbors]
            
            for neighbor in node.neighbors:
                queue.append(neighbor)
        
        # Convert dictionary to list format
        max_val = max(result.keys())
        return [result.get(i, []) for i in range(1, max_val + 1)]
    
    def print(self):
        print(self.toList())

class Solution:
    def cloneGraph(self, node: 'Node') -> 'Node':
        if node is None:
            return None
        
        self.mapVisited = {}
        return self.dfs(node)

    def dfs(self, node: Optional['Node']) -> 'Node':
        if node is None:
            return None
        
        if node in self.mapVisited:
            return self.mapVisited[node]
        
        clone = Node(node.val)
        self.mapVisited[node] = clone
        
        for neighbor in node.neighbors:
            clone.neighbors.append(self.dfs(neighbor))
        
        return clone
    
if __name__ == '__main__':
    node = Node().fromList([[2,4],[1,3],[2,4],[1,3]])
    clonedNode = Solution().cloneGraph(node)
    assert clonedNode.toList() == [[2,4],[1,3],[2,4],[1,3]]