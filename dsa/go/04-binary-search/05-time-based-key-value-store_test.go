package binary_search

import "testing"

type TimeMap struct{ store map[string][]entry }
type entry struct{ t int; v string }

func ConstructorTimeMap() TimeMap { return TimeMap{store: map[string][]entry{}} }
func (tm *TimeMap) Set(key, value string, timestamp int) { tm.store[key] = append(tm.store[key], entry{t: timestamp, v: value}) }
func (tm *TimeMap) Get(key string, timestamp int) string {
    arr := tm.store[key]
    l, r := 0, len(arr)-1
    ans := ""
    for l <= r {
        m := l + (r-l)/2
        if arr[m].t <= timestamp { ans = arr[m].v; l = m + 1 } else { r = m - 1 }
    }
    return ans
}

func TestTimeMap(t *testing.T) {
    tm := ConstructorTimeMap()
    tm.Set("foo", "bar", 1)
    if tm.Get("foo", 1) != "bar" { t.Fatalf("expected bar") }
    if tm.Get("foo", 3) != "bar" { t.Fatalf("expected bar at 3") }
    tm.Set("foo", "bar2", 4)
    if tm.Get("foo", 4) != "bar2" { t.Fatalf("expected bar2") }
    if tm.Get("foo", 5) != "bar2" { t.Fatalf("expected bar2 at 5") }
}

