package linked_list

import "testing"

// LRUNode is one entry in the cache's doubly linked list. The key is stored on the node so
// evicting the tail can also drop the matching map entry.
type LRUNode struct {
    key   string
    value string
    prev  *LRUNode
    next  *LRUNode
}

// LRU is a fixed-size cache that evicts the least recently used entry.
// Approach: hash map for O(1) lookup, plus a doubly linked list with sentinel head/tail nodes
// for O(1) reordering — the sentinels remove every nil check from the splice operations.
// time: O(1) per operation, space: O(size) — map entries plus one node each
type LRU struct {
    size    int
    mapData map[string]*LRUNode
    head    *LRUNode
    tail    *LRUNode
}

// ConstructorLRU wires the two sentinels together around an empty list.
// time: O(1), space: O(size)
func ConstructorLRU(size int) *LRU {
    head, tail := &LRUNode{}, &LRUNode{}
    head.next = tail
    tail.prev = head

    return &LRU{
        size:    size,
        mapData: map[string]*LRUNode{},
        head:    head,
        tail:    tail,
    }
}

// Set inserts or updates a key and marks it most recently used.
// time: O(1), space: O(1) — map lookup plus constant list splicing
func (l *LRU) Set(key, val string) {
    if node, ok := l.mapData[key]; ok {
        node.value = val
        l.rearrangeNodes(node)
        return
    }

    newNode := &LRUNode{key: key, value: val}
    l.mapData[key] = newNode
    l.addHead(newNode)
    l.cleanupLRU()
}

// Get returns the value for key and marks it most recently used, or "" when absent.
// time: O(1), space: O(1)
func (l *LRU) Get(key string) string {
    node, ok := l.mapData[key]
    if !ok { return "" }

    l.rearrangeNodes(node)
    return node.value
}

// addHead splices a node in just after the head sentinel.
// time: O(1), space: O(1)
func (l *LRU) addHead(node *LRUNode) {
    nextNode := l.head.next
    l.head.next = node
    node.prev = l.head
    node.next = nextNode
    nextNode.prev = node
}

// detach unlinks a node from its neighbours.
// time: O(1), space: O(1)
func (l *LRU) detach(node *LRUNode) {
    node.prev.next = node.next
    node.next.prev = node.prev
}

// rearrangeNodes moves an existing node to the front of the list.
// time: O(1), space: O(1)
func (l *LRU) rearrangeNodes(node *LRUNode) {
    l.detach(node)
    l.addHead(node)
}

// cleanupLRU drops the node just before the tail sentinel once the cache overflows.
// time: O(1), space: O(1)
func (l *LRU) cleanupLRU() {
    if len(l.mapData) <= l.size { return }

    lru := l.tail.prev // actual least-recently-used node
    l.detach(lru)
    delete(l.mapData, lru.key)
}

func TestLRU(t *testing.T) {
    cache := ConstructorLRU(2)

    cache.Set("a", "1")
    cache.Set("b", "2")
    if got := cache.Get("a"); got != "1" { t.Fatalf(`Get("a") = %q, want "1"`, got) }

    // "b" is now the least recently used, so adding "c" should evict it
    cache.Set("c", "3")
    if got := cache.Get("b"); got != "" { t.Fatalf(`Get("b") = %q, want "" after eviction`, got) }
    if got := cache.Get("a"); got != "1" { t.Fatalf(`Get("a") = %q, want "1"`, got) }
    if got := cache.Get("c"); got != "3" { t.Fatalf(`Get("c") = %q, want "3"`, got) }

    // updating an existing key refreshes it without growing the cache
    cache.Set("a", "9")
    if got := cache.Get("a"); got != "9" { t.Fatalf(`Get("a") = %q, want "9"`, got) }
    if len(cache.mapData) != 2 { t.Fatalf("cache size = %d, want 2", len(cache.mapData)) }

    // "c" is now least recently used
    cache.Set("d", "4")
    if got := cache.Get("c"); got != "" { t.Fatalf(`Get("c") = %q, want "" after eviction`, got) }
}
