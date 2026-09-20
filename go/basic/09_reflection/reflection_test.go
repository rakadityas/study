package reflection

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestJSONTags(t *testing.T) {
	u := User{ID: 1, Name: "Ada", Password: "hunter2"}

	out, err := ToJSON(u)
	if err != nil {
		t.Fatalf("ToJSON = %v", err)
	}

	if !strings.Contains(out, `"id": 1`) || !strings.Contains(out, `"name": "Ada"`) {
		t.Fatalf("output missing renamed fields:\n%s", out)
	}
	// omitempty dropped the zero-valued Email and Age
	if strings.Contains(out, "email") || strings.Contains(out, "age") {
		t.Fatalf("omitempty did not drop zero fields:\n%s", out)
	}
	// "-" kept the password out entirely
	if strings.Contains(out, "hunter2") || strings.Contains(out, "Password") {
		t.Fatalf("the password leaked into JSON:\n%s", out)
	}
}

func TestCustomMarshalling(t *testing.T) {
	r := Reading{Station: "north", Value: 21.5}

	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("Marshal = %v", err)
	}
	if got, want := string(b), `{"station":"north","value":"21.5°C"}`; got != want {
		t.Fatalf("Marshal = %s, want %s", got, want)
	}

	var back Reading
	if err := json.Unmarshal(b, &back); err != nil {
		t.Fatalf("Unmarshal = %v", err)
	}
	if back.Station != "north" || back.Value != 21.5 {
		t.Fatalf("round trip = %+v, want the original", back)
	}

	if err := json.Unmarshal([]byte(`{"value":"not-a-temp"}`), &back); err == nil {
		t.Fatal("expected an error for a malformed temperature")
	}
}

func TestDescribe(t *testing.T) {
	fields, err := Describe(User{})
	if err != nil {
		t.Fatalf("Describe = %v", err)
	}
	if len(fields) != 5 {
		t.Fatalf("got %d fields, want 5", len(fields))
	}

	byName := map[string]FieldInfo{}
	for _, f := range fields {
		byName[f.Name] = f
	}

	if byName["ID"].JSONName != "id" {
		t.Fatalf("ID JSON name = %q, want %q", byName["ID"].JSONName, "id")
	}
	if byName["Email"].JSONName != "email" {
		t.Fatalf("the omitempty option leaked into the name: %q", byName["Email"].JSONName)
	}
	if byName["Age"].Type != "int" {
		t.Fatalf("Age type = %q, want int", byName["Age"].Type)
	}

	// a pointer is dereferenced for us
	if _, err := Describe(&User{}); err != nil {
		t.Fatalf("Describe(pointer) = %v", err)
	}
	if _, err := Describe(42); err == nil {
		t.Fatal("Describe(non-struct) should fail")
	}
}

func TestValidate(t *testing.T) {
	if err := Validate(User{ID: 1, Name: "Ada", Age: 36}); err != nil {
		t.Fatalf("Validate on a valid user = %v", err)
	}

	err := Validate(User{Name: strings.Repeat("x", 60), Age: 200})
	if err == nil {
		t.Fatal("Validate should have rejected this user")
	}

	msg := err.Error()
	for _, want := range []string{"ID is required", "Name must be at most 50", "Age must be at most 150"} {
		if !strings.Contains(msg, want) {
			t.Fatalf("error %q missing %q", msg, want)
		}
	}
}

func TestSetDefaults(t *testing.T) {
	// an all-zero config picks up every default
	var c Config
	if err := SetDefaults(&c); err != nil {
		t.Fatalf("SetDefaults = %v", err)
	}
	if c.Host != "localhost" || c.Port != 8080 || !c.Debug {
		t.Fatalf("defaults not applied: %+v", c)
	}
	if c.Timeout != 0 {
		t.Fatalf("Timeout has no default tag, want 0, got %d", c.Timeout)
	}

	// values the caller already set are left alone
	explicit := Config{Host: "example.com", Port: 9000}
	if err := SetDefaults(&explicit); err != nil {
		t.Fatalf("SetDefaults = %v", err)
	}
	if explicit.Host != "example.com" || explicit.Port != 9000 {
		t.Fatalf("SetDefaults clobbered explicit values: %+v", explicit)
	}

	// a non-pointer cannot be settable
	if err := SetDefaults(Config{}); err == nil {
		t.Fatal("SetDefaults on a value should fail — reflection needs an addressable target")
	}
}
