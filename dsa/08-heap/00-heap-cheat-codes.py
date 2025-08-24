import heapq

#### Basic Operations
# Create a heap from a list
arr = [3, 1, 4, 1, 5, 9, 2, 6]
heapq.heapify(arr)  # O(n) - converts list to heap in-place
print(arr)  # [1, 1, 2, 3, 5, 9, 4, 6]

# Push element
heapq.heappush(arr, 0)  # O(log n)
print(arr)  # [0, 1, 2, 1, 5, 9, 4, 6, 3]

# Pop smallest element
smallest = heapq.heappop(arr)  # O(log n)
print(smallest)  # 0
print(arr)  # [1, 1, 2, 3, 5, 9, 4, 6]

# Push and pop in one operation
result = heapq.heappushpop(arr, 7)  # Push 7, then pop smallest
print(result)  # 1

# Pop and push in one operation
result = heapq.heapreplace(arr, 8)  # Pop smallest, then push 8
print(result)  # 1

#### Advanced Functions
# Get n largest/smallest elements
arr = [3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5]

# N largest elements
largest_3 = heapq.nlargest(3, arr)
print(largest_3)  # [9, 6, 5]

# N smallest elements
smallest_3 = heapq.nsmallest(3, arr)
print(smallest_3)  # [1, 1, 2]

# With key function
students = [('Alice', 85), ('Bob', 90), ('Charlie', 78)]
top_2 = heapq.nlargest(2, students, key=lambda x: x[1])
print(top_2)  # [('Bob', 90), ('Alice', 85)]

#### Max Heap
class MaxHeap:
    def __init__(self):
        self.heap = []
    
    def push(self, val):
        heapq.heappush(self.heap, -val)
    
    def pop(self):
        return -heapq.heappop(self.heap)
    
    def peek(self):
        return -self.heap[0] if self.heap else None
    
    def size(self):
        return len(self.heap)

# Usage
max_heap = MaxHeap()
max_heap.push(3)
max_heap.push(1)
max_heap.push(4)
print(max_heap.pop())  # 4 (largest)
print(max_heap.peek()) # 3


#### Priority Queue
class PriorityQueue:
    def __init__(self):
        self.heap = []
        self.index = 0
    
    def push(self, item, priority):
        heapq.heappush(self.heap, (priority, self.index, item))
        self.index += 1
    
    def pop(self):
        return heapq.heappop(self.heap)[2]

# Usage
pq = PriorityQueue()
pq.push("task1", 3)
pq.push("task2", 1)  # Higher priority (lower number)
pq.push("task3", 2)

print(pq.pop())  # "task2" (priority 1)
print(pq.pop())  # "task3" (priority 2)

#### K Largest/Smallest Elements
def k_largest_elements(arr, k):
    # Method 1: Using nlargest
    return heapq.nlargest(k, arr)

def k_largest_streaming(arr, k):
    # Method 2: Maintain min-heap of size k
    heap = []
    for num in arr:
        if len(heap) < k:
            heapq.heappush(heap, num)
        elif num > heap[0]:
            heapq.heapreplace(heap, num)
    return sorted(heap, reverse=True)

# Example
arr = [3, 1, 4, 1, 5, 9, 2, 6]
print(k_largest_elements(arr, 3))    # [9, 6, 5]
print(k_largest_streaming(arr, 3))   # [9, 6, 5]

#### Merge K Sorted Lists
def merge_k_sorted_lists(lists):
    heap = []
    result = []
    
    # Initialize heap with first element from each list
    for i, lst in enumerate(lists):
        if lst:
            heapq.heappush(heap, (lst[0], i, 0))  # (value, list_index, element_index)
    
    while heap:
        val, list_idx, elem_idx = heapq.heappop(heap)
        result.append(val)
        
        # Add next element from the same list
        if elem_idx + 1 < len(lists[list_idx]):
            next_val = lists[list_idx][elem_idx + 1]
            heapq.heappush(heap, (next_val, list_idx, elem_idx + 1))
    
    return result

# Example
lists = [[1, 4, 5], [1, 3, 4], [2, 6]]
print(merge_k_sorted_lists(lists))  # [1, 1, 2, 3, 4, 4, 5, 6]

#### Heap Properties
# Heap is represented as a list where:
# - Parent of node at index i is at index (i-1)//2
# - Left child of node at index i is at index 2*i+1
# - Right child of node at index i is at index 2*i+2

def heap_parent(i):
    return (i - 1) // 2

def heap_left_child(i):
    return 2 * i + 1

def heap_right_child(i):
    return 2 * i + 2

# Example heap: [1, 3, 2, 7, 5, 4, 6]
#       1
#      / \
#     3   2
#    / \ / \
#   7  5 4  6


## Real-World Applications
# 1. Dijkstra's Algorithm : Finding shortest paths
# 2. Task Scheduling : Priority-based task execution
# 3. Data Streaming : Finding top-K elements in streams
# 4. Huffman Coding : Building optimal prefix codes
# 5. A Search *: Pathfinding with heuristics
# 6. Memory Management : Heap allocation algorithms

## Best Practices
# 1. Use heapq.nlargest() and heapq.nsmallest() for small k values
# 2. Maintain heap size when dealing with streaming data
# 3. Use tuples for complex objects with priority
# 4. Remember Python uses min-heap - negate values for max-heap behavior
# 5. Consider heapreplace() instead of separate pop/push operations


# Time Complexities:
# heapify(arr)           - O(n)
# heappush(heap, item)   - O(log n)
# heappop(heap)          - O(log n)
# heappushpop(heap, item)- O(log n)
# heapreplace(heap, item)- O(log n)
# nlargest(k, arr)       - O(n log k)
# nsmallest(k, arr)      - O(n log k)
# heap[0] (peek min)     - O(1)

# Space Complexities:
# heapify(arr)           - O(1) - in-place
# heappush/heappop       - O(1) - per operation
# nlargest/nsmallest     - O(k) - for result storage
# Overall heap storage   - O(n) - where n is number of elements

### Common Heap Algorithms
# K Largest Elements:
# Time: O(n log k), Space: O(k)

# K Smallest Elements:
# Time: O(n log k), Space: O(k)

# Heap Sort:
# Time: O(n log n), Space: O(1)

# Priority Queue Operations:
# Insert: O(log n), Extract: O(log n), Peek: O(1)
# Space: O(n)

# Merge K Sorted Arrays:
# Time: O(n log k), Space: O(k)
# where n = total elements, k = number of arrays

# Dijkstra's Algorithm with Heap:
# Time: O((V + E) log V), Space: O(V)
# where V = vertices, E = edges

#### Memory Usage Notes
# Heap vs Sorted Array:
# Heap: O(log n) insert/delete, O(1) peek
# Sorted Array: O(n) insert/delete, O(1) access by index

# Heap vs Binary Search Tree:
# Heap: Better for priority queue, simpler implementation
# BST: Better for range queries, ordered traversal

# Space overhead:
# - Python list: ~8 bytes per element (64-bit)
# - Heap structure: No additional overhead beyond list
# - Custom heap classes: Additional object overhead