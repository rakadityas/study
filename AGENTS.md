# Study Repo — Agent Guide

Personal study repo for algorithms, protocols, and language fundamentals.
Use this file to understand how the repo is organized and how to work within it.

---

## Repository Structure

```
.
├── dsa/
│   ├── go/          # DSA solutions in Go (single Go module)
│   └── python/      # DSA solutions in Python (standalone scripts)
├── go/              # Go protocol/networking samples (one module each)
│   ├── grpc/
│   ├── http/
│   ├── hls-dash/
│   ├── p2p/
│   ├── sse/
│   ├── webrtc/
│   ├── websocket/
│   └── go_bench/    # Go benchmarking examples
├── python/
│   └── oop_basic/   # Python OOP fundamentals
├── ruby/
│   └── ruby_basic/  # Ruby syntax basics
└── sql/             # LeetCode SQL problems (flat .sql files)
```

---

## DSA Section (`dsa/`)

### Topic Layout

Both Go and Python follow the same 12-topic numbered structure:

| # | Topic |
|---|-------|
| 01 | array-hashes |
| 02 | two-pointers |
| 03 | stack |
| 04 | binary-search |
| 05 | sliding-window |
| 06 | linked-list |
| 07 | tree |
| 08 | sort |
| 09 | backtracking |
| 10 | heap |
| 11 | graph |
| 12 | greedy (Python only) |

File naming: `NN-problem-name.py` / `NN-problem-name_test.go`

---

### Python DSA (`dsa/python/`)

- Each file is **self-contained** — solution class/function + assertions in `if __name__ == '__main__':`.
- Run a single file: `python3 dsa/python/01-array-hashes/01-contains-duplicate.py`
- No test runner, no external dependencies.

**File header convention:**
```python
# https://leetcode.com/problems/<slug>/description/
# time complexity: O(n) — brief explanation
# space complexity: O(n) — brief explanation
```

**Multi-variant files** (e.g. two approaches in one file) have per-class or per-method complexity blocks:
```python
# Approach: description of strategy
# time complexity: O(n)
# space complexity: O(1)
class SolutionOptimal:
    ...
```

---

### Go DSA (`dsa/go/`)

- **Single Go module**: `github.com/rakadityas/study/dsa/go` (rooted at `dsa/go/go.mod`, Go 1.20).
- Each file is `*_test.go` — solution function + test function live together in the same file.
- Uses standard library only (`testing` package). No external dependencies.
- Run all tests: `cd dsa/go && go test ./...`
- Run one package: `cd dsa/go && go test ./01-array-hashes/`
- Run one test: `cd dsa/go && go test ./06-linked-list/ -run TestReverseList`

**Doc comment convention** (Go doc style, above the solution function):
```go
// functionName one-line description.
// time: O(n), space: O(1) — brief explanation of why
func functionName(...) ... {
```

**Multi-variant** (struct with multiple methods, or multiple top-level functions) each get their own comment block.

**Note:** The Go DSA module has a `common/` package with shared test helpers (e.g. `Sort2DStringSlice`).

---

## Go Protocol Samples (`go/`)

Each subdirectory is an **independent Go module** with its own `go.mod`. They demonstrate protocol patterns:

| Directory | What it shows |
|-----------|---------------|
| `grpc/` | gRPC server/client |
| `http/` | HTTP server patterns |
| `websocket/` | WebSocket communication |
| `sse/` | Server-Sent Events |
| `p2p/` | Peer-to-peer networking |
| `hls-dash/` | HLS/DASH streaming |
| `webrtc/` | WebRTC signaling |
| `go_bench/` | Go benchmark patterns (5 progressive examples) |

Run a sample: `cd go/<dir> && go run ./cmd/...`

---

## SQL (`sql/`)

Flat directory of numbered LeetCode SQL files: `NN-problem-name.sql`.
No runner — files contain the query only.

---

## Commit Convention

Format: `type(scope): description` (conventional commits, lowercase)

**Types used:**
- `feat` — new solution or variant added
- `fix` / `bug` — bug fix or incorrect logic corrected
- `wip` — work in progress / structural changes

**Scopes used:**
- `dsa` — any change inside `dsa/`
- `go` — changes inside `go/` protocol samples
- `dir` — repo-level structural changes

**Examples from history:**
```
feat(dsa): account balance greedy algorithm
feat(dsa): invert binary tree
feat(dsa): kth largest element in a stream optimal
fix(dsa): improve count of substring
bug(dsa): fix bug group anagrams
feat(go): protocol sample
```

**Rules:**
- One problem or improvement per commit.
- Description names the problem or describes what changed, no period at end.
- No commit body is used — the subject line is the entire message.
- Commits go directly to `main`; no PR/branch workflow observed.

---

## Adding a New DSA Solution

### Python
1. Create `dsa/python/NN-topic/NN-problem-name.py`.
2. Add the LeetCode URL, time/space complexity comments at the top.
3. Implement the solution class/function.
4. Add assertions in `if __name__ == '__main__':`.
5. Commit: `feat(dsa): <problem name>`

### Go
1. Create `dsa/go/NN-topic/NN-problem-name_test.go` with `package <topic_package>`.
2. Add a doc comment above the solution function with `time:` and `space:` lines.
3. Implement the solution function and a `TestXxx` function.
4. Run `cd dsa/go && go test ./NN-topic/` to verify.
5. Commit: `feat(dsa): <problem name>`

---

## Complexity Comment Rules

Always annotate every solution (both languages) with:
- **Time complexity** and **space complexity**.
- A short inline explanation after `—` explaining *why* (e.g. what data structure drives the cost).
- For multi-variant files, annotate **each variant separately** with an `Approach:` line explaining the strategy.
- When correcting a wrong annotation, fix both the complexity value and the explanation.

---

## Key Differences Between Go and Python Implementations

Some problems have a better space complexity in Go than in Python:

| Problem | Python | Go | Notes |
|---------|--------|----|-------|
| Reorder List | O(n) space | O(1) space | Go uses slow/fast + in-place reverse |
| Remove Nth From End | O(n) (base variant) | O(1) | Go uses two-pointer gap directly |
| Linked List Cycle | O(n) (hashmap variant first) | O(1) | Go implements Floyd's only |
| Remove Duplicates (sorted) | O(n) (uses defaultdict) | O(1) | Go compares adjacent elements |
