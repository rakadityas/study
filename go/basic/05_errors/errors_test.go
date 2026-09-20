package errorsdemo

import (
	"errors"
	"strings"
	"testing"
)

func TestWrappingPreservesSentinel(t *testing.T) {
	s := NewStore(map[string]string{"1": "ada"})

	if got, err := FetchUser(s, "1"); err != nil || got != "ada" {
		t.Fatalf("FetchUser(1) = (%q, %v), want (\"ada\", nil)", got, err)
	}

	_, err := FetchUser(s, "404")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want it to wrap ErrNotFound", err)
	}

	// the message carries every layer of context
	if msg := err.Error(); !strings.Contains(msg, "users.404") || !strings.Contains(msg, "not found") {
		t.Fatalf("err message %q lost context", msg)
	}

	// errors.As recovers the wrapper type so its fields are readable
	var qe *QueryError
	if !errors.As(err, &qe) || qe.Query != "users.404" {
		t.Fatalf("errors.As did not recover *QueryError from %v", err)
	}
}

func TestValidationErrorAs(t *testing.T) {
	err := ValidateAge(200)

	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("errors.As did not find *ValidationError in %v", err)
	}
	if ve.Field != "age" || ve.Value != 200 {
		t.Fatalf("ValidationError = %+v, want field age / value 200", ve)
	}

	if err := ValidateAge(30); err != nil {
		t.Fatalf("ValidateAge(30) = %v, want nil", err)
	}
}

func TestValidateAllJoins(t *testing.T) {
	err := ValidateAll("", 200)
	if err == nil {
		t.Fatal("ValidateAll should report both failures")
	}

	msg := err.Error()
	if !strings.Contains(msg, "name") || !strings.Contains(msg, "age") {
		t.Fatalf("joined error %q should mention both fields", msg)
	}

	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("errors.As should reach into a joined error, got %v", err)
	}

	if err := ValidateAll("ada", 36); err != nil {
		t.Fatalf("ValidateAll on valid input = %v, want nil", err)
	}
}

func TestClassify(t *testing.T) {
	s := NewStore(map[string]string{})
	_, notFound := FetchUser(s, "x")

	cases := []struct {
		err  error
		want string
	}{
		{nil, "ok"},
		{notFound, "not-found"},
		{ErrUnauthorized, "unauthorized"},
		{ValidateAge(-1), "validation:age"},
		{errors.New("something else"), "unknown"},
	}

	for _, c := range cases {
		if got := Classify(c.err); got != c.want {
			t.Fatalf("Classify(%v) = %q, want %q", c.err, got, c.want)
		}
	}
}
