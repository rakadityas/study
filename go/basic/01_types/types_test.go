package types

import "testing"

func TestImplicitSatisfaction(t *testing.T) {
	notifiers := []Notifier{EmailNotifier{Addr: "a@b.c"}, SMSNotifier{Number: "123"}}

	want := []string{"email to a@b.c: hello", "sms to 123: hello"}
	for i, n := range notifiers {
		if got := n.Notify("hello"); got != want[i] {
			t.Fatalf("Notify() = %q, want %q", got, want[i])
		}
	}
}

func TestEmbeddingPromotesAndShadows(t *testing.T) {
	s := Service{Base: Base{ID: "7"}, Name: "billing"}

	if got, want := s.Describe(), "service:billing (base:7)"; got != want {
		t.Fatalf("s.Describe() = %q, want %q", got, want)
	}
	// the shadowed method is still reachable through the embedded field
	if got, want := s.Base.Describe(), "base:7"; got != want {
		t.Fatalf("s.Base.Describe() = %q, want %q", got, want)
	}
	// promoted fields work the same way
	if s.ID != "7" {
		t.Fatalf("s.ID = %q, want %q", s.ID, "7")
	}
}

func TestReceiverKind(t *testing.T) {
	c := Counter{}

	c.IncValue()
	if c.N != 0 {
		t.Fatalf("value receiver mutated the original: N = %d, want 0", c.N)
	}

	// Go auto-takes the address here because c is addressable
	c.IncPointer()
	if c.N != 1 {
		t.Fatalf("pointer receiver did not mutate: N = %d, want 1", c.N)
	}
}

func TestTypedNilTrap(t *testing.T) {
	if err := BadNilReturn(false); err == nil {
		t.Fatal("expected the typed-nil trap to produce a non-nil error")
	}

	if err := GoodNilReturn(false); err != nil {
		t.Fatalf("GoodNilReturn(false) = %v, want nil", err)
	}
	if err := GoodNilReturn(true); err == nil {
		t.Fatal("GoodNilReturn(true) should report an error")
	}
}

func TestClassify(t *testing.T) {
	cases := []struct {
		in   any
		want string
	}{
		{nil, "nil"},
		{42, "int:42"},
		{"hi", `string:"hi"`},
		{[]int{1, 2, 3}, "[]int len=3"},
		{EmailNotifier{Addr: "x@y.z"}, "notifier:email to x@y.z: hi"},
		{&MyError{Msg: "bad"}, "error:bad"},
		{3.5, "other:float64"},
	}

	for _, c := range cases {
		if got := Classify(c.in); got != c.want {
			t.Fatalf("Classify(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}
