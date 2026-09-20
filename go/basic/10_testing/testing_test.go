package testingdemo

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// Table-driven tests with subtests
// ---------------------------------------------------------------------------

// The default shape for a pure function. Each case gets a name, so `go test -run
// 'TestParseTags/dedup'` runs exactly one, and a failure report names the case rather than
// a line number.
func TestParseTags(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []string
	}{
		{"simple", "go,rust", []string{"go", "rust"}},
		{"trims whitespace", " go , rust ", []string{"go", "rust"}},
		{"lowercases", "Go,RUST", []string{"go", "rust"}},
		{"dedups", "go,Go,GO", []string{"go"}},
		{"drops blanks", "go,,rust,", []string{"go", "rust"}},
		{"empty input", "", []string{}},
		{"only separators", ",,,", []string{}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := ParseTags(c.in)

			if len(got) != len(c.want) {
				t.Fatalf("ParseTags(%q) = %v, want %v", c.in, got, c.want)
			}
			for i := range got {
				if got[i] != c.want[i] {
					t.Fatalf("ParseTags(%q) = %v, want %v", c.in, got, c.want)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Fakes and helpers
// ---------------------------------------------------------------------------

// fakeStore is an in-memory Store. A hand-written fake like this is usually better than a
// mocking framework: it is obvious, it compiles, and it fails loudly when the interface
// changes.
type fakeStore struct {
	data    map[string]string
	putErr  error
	putCall int
}

func newFakeStore() *fakeStore { return &fakeStore{data: map[string]string{}} }

func (f *fakeStore) Get(key string) (string, error) {
	v, ok := f.data[key]
	if !ok {
		return "", ErrNoRecord
	}
	return v, nil
}

func (f *fakeStore) Put(key, value string) error {
	f.putCall++
	if f.putErr != nil {
		return f.putErr
	}
	f.data[key] = value
	return nil
}

// newService is a fixture. t.Helper() makes failures point at the caller's line rather than
// here, and t.Cleanup registers teardown next to setup so the two cannot drift apart.
func newService(t *testing.T) (*TagService, *fakeStore) {
	t.Helper()

	store := newFakeStore()
	t.Cleanup(func() { store.data = nil })

	return NewTagService(store), store
}

func TestTagServiceSaveAndLoad(t *testing.T) {
	svc, store := newService(t)

	got, err := svc.Save("post-1", "Go, rust ,go")
	if err != nil {
		t.Fatalf("Save = %v, want nil", err)
	}
	if strings.Join(got, ",") != "go,rust" {
		t.Fatalf("Save = %v, want [go rust]", got)
	}
	if store.putCall != 1 {
		t.Fatalf("Put called %d times, want 1", store.putCall)
	}

	back, err := svc.Load("post-1")
	if err != nil {
		t.Fatalf("Load = %v, want nil", err)
	}
	if strings.Join(back, ",") != "go,rust" {
		t.Fatalf("Load = %v, want [go rust]", back)
	}
}

func TestTagServiceErrors(t *testing.T) {
	t.Run("no usable tags", func(t *testing.T) {
		svc, store := newService(t)

		if _, err := svc.Save("post-1", " , , "); err == nil {
			t.Fatal("Save should reject an input with no usable tags")
		}
		if store.putCall != 0 {
			t.Fatal("Save should not touch the store when there is nothing to save")
		}
	})

	t.Run("store failure is wrapped", func(t *testing.T) {
		svc, store := newService(t)
		boom := errors.New("disk full")
		store.putErr = boom

		_, err := svc.Save("post-1", "go")
		if !errors.Is(err, boom) {
			t.Fatalf("Save err = %v, want it to wrap %v", err, boom)
		}
	})

	t.Run("missing key", func(t *testing.T) {
		svc, _ := newService(t)

		if _, err := svc.Load("nope"); !errors.Is(err, ErrNoRecord) {
			t.Fatalf("Load err = %v, want ErrNoRecord", err)
		}
	})
}

// ---------------------------------------------------------------------------
// Parallel tests
// ---------------------------------------------------------------------------

// t.Parallel() pauses a subtest until its siblings are also ready, then runs them together.
// It catches shared-state bugs and shortens slow suites — but only when each case owns its
// own fixture, as these do.
func TestParseTagsParallel(t *testing.T) {
	cases := []struct{ in, want string }{
		{"a,b", "a,b"},
		{"B,A,b", "b,a"},
		{" c ", "c"},
	}

	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			t.Parallel()
			if got := strings.Join(ParseTags(c.in), ","); got != c.want {
				t.Fatalf("ParseTags(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Examples: documentation that is compiled and verified
// ---------------------------------------------------------------------------

// An Example function appears in godoc AND runs as a test — the Output comment is compared
// against stdout, so documentation that drifts from behaviour fails the build.
func ExampleParseTags() {
	fmt.Println(ParseTags("Go, rust , GO"))
	// Output: [go rust]
}

func ExampleNormalize() {
	fmt.Println(Normalize([]string{"rust", "go", "zig"}))
	// Output: [go rust zig]
}

// ---------------------------------------------------------------------------
// Benchmarks
// ---------------------------------------------------------------------------

// Run with: go test -bench=. -benchmem ./10_testing
//
// b.N is chosen by the framework. Anything that is setup rather than the thing being
// measured goes before b.ResetTimer.
func BenchmarkParseTags(b *testing.B) {
	input := strings.Repeat("Go, Rust, Zig, ", 20)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ParseTags(input)
	}
}

func BenchmarkParseTagsBySize(b *testing.B) {
	for _, size := range []int{1, 10, 100} {
		input := strings.Repeat("tag,", size)

		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.ReportAllocs() // per-benchmark equivalent of -benchmem
			for i := 0; i < b.N; i++ {
				ParseTags(input)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Fuzzing
// ---------------------------------------------------------------------------

// Run the corpus only:      go test ./10_testing
// Actually fuzz for 10s:    go test -fuzz=FuzzParseTags -fuzztime=10s ./10_testing
//
// A fuzz target asserts *invariants* rather than exact outputs, since the input is generated.
func FuzzParseTags(f *testing.F) {
	for _, seed := range []string{"", ",", "go", "Go,GO", " a , b ", strings.Repeat("x,", 50)} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, in string) {
		tags := ParseTags(in)

		seen := map[string]bool{}
		for _, tag := range tags {
			if tag == "" {
				t.Fatalf("ParseTags(%q) produced an empty tag", in)
			}
			if tag != strings.ToLower(tag) {
				t.Fatalf("ParseTags(%q) produced an uppercase tag %q", in, tag)
			}
			if tag != strings.TrimSpace(tag) {
				t.Fatalf("ParseTags(%q) produced an untrimmed tag %q", in, tag)
			}
			if seen[tag] {
				t.Fatalf("ParseTags(%q) produced duplicate %q", in, tag)
			}
			seen[tag] = true
		}

		// idempotent: re-parsing the joined output changes nothing
		if again := ParseTags(strings.Join(tags, ",")); len(again) != len(tags) {
			t.Fatalf("ParseTags is not idempotent for %q: %v then %v", in, tags, again)
		}
	})
}
