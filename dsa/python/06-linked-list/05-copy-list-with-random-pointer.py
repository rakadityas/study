# https://leetcode.com/problems/copy-list-with-random-pointer/description/
# time complexity: O(n)
# space complexity: O(n)

from typing import Optional
from collections import defaultdict

class Node:
    def __init__(self, val: int = 0, next: Optional['Node'] = None, random: Optional['Node'] = None):
        self.val = val
        self.next = next
        self.random = random

    def __repr__(self):
        # Helpful for debugging
        next_val = self.next.val if self.next else None
        random_val = self.random.val if self.random else None
        return f"Node(val={self.val}, next={next_val}, random={random_val})"


class Solution:
    def copyRandomList(self, head: Optional[Node]) -> Optional[Node]:
        if head is None:
            return None

        node_map = {}

        current = head
        while current:
            copy_node = Node(current.val)
            node_map[current] = copy_node
            current = current.next

        current = head
        while current:
            if current.next:
                node_map[current].next = node_map[current.next]
            else:
                node_map[current].next = None

            if current.random:
                node_map[current].random = node_map[current.random]
            else:
                node_map[current].random = None

            current = current.next

        return node_map[head]