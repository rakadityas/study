class Node:
    def __init__(self, key=None, value=None, prev=None, next=None):
        self.key = key
        self.value = value
        self.prev = prev
        self.next = next

class LRU:
    # Time: O(1), Space: O(size) for the map + doubly linked list nodes
    def __init__(self, size: int):
        self.size = size
        self.mapData = {}
        self.head = Node()
        self.tail = Node()
        self.head.next = self.tail
        self.tail.prev = self.head

    # Time: O(1) - dict lookup + O(1) list ops, Space: O(1)
    def Set(self, key: str, val: str):
        if key in self.mapData:
            node = self.mapData[key]
            node.value = val
            self.rearrangeNodes(node)
        else:
            newNode = Node(key, val)
            self.mapData[key] = newNode
            self.addHead(newNode)
            self.cleanupLRU()

    # Time: O(1) - dict lookup + O(1) list ops, Space: O(1)
    def Get(self, key: str):
        if key not in self.mapData:
            return ""
        node = self.mapData[key]
        self.rearrangeNodes(node)
        return node.value

    # Time: O(1), Space: O(1)
    def addHead(self, node: Node):
        nextNode = self.head.next
        self.head.next = node
        node.prev = self.head
        node.next = nextNode
        nextNode.prev = node

    # Time: O(1), Space: O(1)
    def detach(self, node: Node):
        node.prev.next = node.next
        node.next.prev = node.prev

    # Time: O(1), Space: O(1)
    def rearrangeNodes(self, node: Node):
        self.detach(node)
        self.addHead(node)

    # Time: O(1), Space: O(1)
    def cleanupLRU(self):
        if len(self.mapData) <= self.size:
            return
        lru = self.tail.prev          # actual least-recently-used node
        self.detach(lru)
        del self.mapData[lru.key]