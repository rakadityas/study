// Package testingdemo covers Go's testing toolkit beyond `func TestX(t *testing.T)`:
// table-driven tests, subtests, fixtures and cleanup, interface-based fakes, benchmarks,
// examples and fuzzing.
//
// Go deliberately ships no assertion library. The idiom is a plain if plus t.Fatalf with a
// message naming the input, what you got and what you wanted — because that message is all
// you have when the test fails in CI six weeks from now.
package testingdemo

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// ---------------------------------------------------------------------------
// The code under test
// ---------------------------------------------------------------------------

// ParseTags splits a comma-separated tag list, trimming, lowercasing, dropping blanks and
// deduplicating. Small, pure and full of edge cases — an ideal table-test subject.
func ParseTags(s string) []string {
	seen := map[string]bool{}
	out := []string{}

	for _, raw := range strings.Split(s, ",") {
		tag := strings.ToLower(strings.TrimSpace(raw))
		if tag == "" || seen[tag] {
			continue
		}
		seen[tag] = true
		out = append(out, tag)
	}

	return out
}

// Normalize sorts a tag list so output is stable regardless of input order.
func Normalize(tags []string) []string {
	out := append([]string{}, tags...)
	sort.Strings(out)
	return out
}

// ---------------------------------------------------------------------------
// A dependency, behind an interface, so tests can substitute a fake
// ---------------------------------------------------------------------------

// ErrNoRecord is returned when a lookup misses.
var ErrNoRecord = errors.New("no record")

// Store is the narrow interface the service depends on. Defining it here — at the consumer,
// with exactly the two methods this package needs — is what makes it trivially fakeable.
type Store interface {
	Get(key string) (string, error)
	Put(key, value string) error
}

// TagService is the unit under test. It takes its dependency as an interface, so the test
// supplies a fake and never needs a real database.
type TagService struct{ store Store }

func NewTagService(s Store) *TagService { return &TagService{store: s} }

// Save parses, normalises and stores a tag list.
func (t *TagService) Save(key, raw string) ([]string, error) {
	tags := Normalize(ParseTags(raw))
	if len(tags) == 0 {
		return nil, fmt.Errorf("save %q: no usable tags", key)
	}

	if err := t.store.Put(key, strings.Join(tags, ",")); err != nil {
		return nil, fmt.Errorf("save %q: %w", key, err)
	}
	return tags, nil
}

// Load reads a tag list back.
func (t *TagService) Load(key string) ([]string, error) {
	raw, err := t.store.Get(key)
	if err != nil {
		return nil, fmt.Errorf("load %q: %w", key, err)
	}
	return ParseTags(raw), nil
}
