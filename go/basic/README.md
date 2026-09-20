# basic — Go advanced crash course

A runnable tour of the Go language and runtime, aimed at someone who can already write Go
and wants the parts that decide whether production code is correct: the type system's sharp
edges, concurrency, context, memory behaviour, and the testing toolkit.

Every topic is a package whose tests *are* the lesson — they assert the behaviour the
comments describe, so nothing here can quietly rot. Read the `.go` file, then read the
`_test.go` file to see the claims proven.

## Run

```bash
cd go/basic

go test ./...              # everything
go test -race ./...        # the concurrency chapters are written to pass under -race
go test -v ./02_concurrency
go test -bench=. -benchmem ./07_memory
```

## Topics

| # | Package | What it covers |
|---|---------|----------------|
| 01 | [01_types](./01_types) | implicit interface satisfaction, embedding vs inheritance, value/pointer receivers and method sets, the typed-nil trap, type switches |
| 02 | [02_concurrency](./02_concurrency) | goroutines, buffered vs unbuffered channels, directional types, pipelines, fan-in/fan-out, `select`, closed-channel semantics |
| 03 | [03_context](./03_context) | cancellation, deadlines, request-scoped values, collision-free context keys, first-result-wins |
| 04 | [04_sync](./04_sync) | `Mutex`/`RWMutex`, `WaitGroup`, `Once`, typed atomics and compare-and-swap, when `sync.Map` is actually the right call |
| 05 | [05_errors](./05_errors) | sentinels, custom error types, `%w` wrapping, `errors.Is`/`As`/`Join`, inspection without string matching |
| 06 | [06_generics](./06_generics) | type parameters, union and `~` constraints, `cmp.Ordered`, generic containers, inference — and when not to bother |
| 07 | [07_memory](./07_memory) | slice headers and aliasing, the append-clobber bug, three-index slices, growth and preallocation, `strings.Builder`, map addressability, escape analysis |
| 08 | [08_patterns](./08_patterns) | worker pools, semaphores for bounded parallelism, an errgroup from scratch, token-bucket rate limiting, graceful shutdown |
| 09 | [09_reflection](./09_reflection) | struct tags, custom JSON marshalling, reflecting over fields, a tag-driven validator, setting values through `reflect` |
| 10 | [10_testing](./10_testing) | table tests and subtests, fixtures with `t.Helper`/`t.Cleanup`, hand-written fakes, parallel tests, examples, benchmarks, fuzzing |

## Things worth running yourself

The escape-analysis and allocation chapters are more convincing when you watch them:

```bash
# where does each value live, and why?
go build -gcflags='-m' ./07_memory

# preallocation vs naive append; strings.Builder vs +=
go test -bench=. -benchmem ./07_memory

# generate inputs until an invariant breaks
go test -fuzz=FuzzParseTags -fuzztime=30s ./10_testing

# the race detector finds bugs tests alone cannot
go test -race ./...
```

Sample output from `07_memory`, which is the whole argument for preallocating in one table:

```text
BenchmarkBuildNaive-8       4779 ns/op    25207 B/op    11 allocs/op
BenchmarkBuildPrealloc-8    1009 ns/op        0 B/op     0 allocs/op
BenchmarkConcat/naive-8    27511 ns/op   174175 B/op   199 allocs/op
BenchmarkConcat/builder-8   1635 ns/op     1792 B/op     1 allocs/op
```

## Notes

- Standard library only; no external dependencies.
- `go.mod` targets Go 1.22. The generics chapter needs 1.18+, `errors.Join` needs 1.20,
  and the typed atomics in `04_sync` need 1.19.
- Unlike `dsa/go`, this module is gofmt-clean — run `gofmt -w .` before committing.
