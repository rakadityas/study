// Package types covers the parts of Go's type system that surprise people coming from
// class-based languages: implicit interface satisfaction, embedding instead of inheritance,
// value vs pointer receivers, and the typed-nil trap.
package types

import "fmt"

// ---------------------------------------------------------------------------
// Implicit interface satisfaction
// ---------------------------------------------------------------------------

// Notifier is satisfied by any type with this method. There is no "implements" keyword —
// the compiler checks structurally at the point of assignment. Interfaces therefore belong
// to the consumer, not the producer: define them where they are used, keep them small.
type Notifier interface {
	Notify(msg string) string
}

// EmailNotifier never mentions Notifier, yet satisfies it.
type EmailNotifier struct{ Addr string }

func (e EmailNotifier) Notify(msg string) string { return "email to " + e.Addr + ": " + msg }

// SMSNotifier likewise.
type SMSNotifier struct{ Number string }

func (s SMSNotifier) Notify(msg string) string { return "sms to " + s.Number + ": " + msg }

// compile-time assertions: these cost nothing at runtime but fail the build if the
// contract is broken. Prefer them over discovering the mismatch at a call site.
var (
	_ Notifier = EmailNotifier{}
	_ Notifier = SMSNotifier{}
)

// ---------------------------------------------------------------------------
// Embedding: composition that promotes methods, not inheritance
// ---------------------------------------------------------------------------

// Base supplies an ID and a Describe method.
type Base struct{ ID string }

func (b Base) Describe() string { return "base:" + b.ID }

// Service embeds Base, so Service values get Describe for free. This is delegation with
// syntactic sugar — there is no vtable and no virtual dispatch. Service.Describe below
// shadows the promoted method; the embedded one is still reachable as s.Base.Describe().
type Service struct {
	Base
	Name string
}

func (s Service) Describe() string { return "service:" + s.Name + " (" + s.Base.Describe() + ")" }

// ---------------------------------------------------------------------------
// Value vs pointer receivers
// ---------------------------------------------------------------------------

// Counter demonstrates why the receiver kind matters. IncValue operates on a copy and is
// therefore a no-op for the caller; IncPointer mutates the original.
type Counter struct{ N int }

func (c Counter) IncValue()    { c.N++ }
func (c *Counter) IncPointer() { c.N++ }

// The method set rule that trips people up: *Counter has both methods, Counter has only the
// value one. So a Counter value does NOT satisfy an interface requiring IncPointer, while
// &Counter{} does. Rule of thumb: if any method needs a pointer receiver, give them all one.

// ---------------------------------------------------------------------------
// The typed-nil trap
// ---------------------------------------------------------------------------

// An interface value is a (type, value) pair. It is nil only when BOTH halves are nil.
// Assigning a nil *concrete pointer* to an interface produces a non-nil interface holding
// a nil pointer — the single most common source of "but I returned nil!" bugs.

type MyError struct{ Msg string }

func (e *MyError) Error() string { return e.Msg }

// BadNilReturn looks like it returns nil on success, but the declared return type is error
// while the variable is *MyError. The implicit conversion packs a nil pointer into a
// non-nil interface.
func BadNilReturn(fail bool) error {
	var err *MyError // nil *MyError
	if fail {
		err = &MyError{Msg: "boom"}
	}
	return err // always non-nil as an error!
}

// GoodNilReturn keeps the concrete type out of the return path until there is a real value.
func GoodNilReturn(fail bool) error {
	if fail {
		return &MyError{Msg: "boom"}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Type assertions and type switches
// ---------------------------------------------------------------------------

// Classify recovers the dynamic type behind an interface value. The comma-ok form never
// panics; the single-value form does, so reach for comma-ok unless a failure is a bug.
func Classify(v any) string {
	switch x := v.(type) {
	case nil:
		return "nil"
	case int:
		return fmt.Sprintf("int:%d", x)
	case string:
		return fmt.Sprintf("string:%q", x)
	case []int:
		return fmt.Sprintf("[]int len=%d", len(x))
	case Notifier:
		// interfaces work as cases too — the first matching case wins, so order matters
		return "notifier:" + x.Notify("hi")
	case error:
		return "error:" + x.Error()
	default:
		return fmt.Sprintf("other:%T", x)
	}
}
