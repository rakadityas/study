// Package reflection covers struct tags, encoding/json, and the reflect package.
//
// Reflection trades compile-time safety and speed for flexibility. It is the right tool for
// serialisation, validation and ORM-style mapping — places where the code genuinely cannot
// know the type. Everywhere else, generics or an interface will be faster and clearer.
package reflection

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ---------------------------------------------------------------------------
// Struct tags and JSON
// ---------------------------------------------------------------------------

// User shows the tag options that matter in practice:
//   - a renamed field, because Go exports with capitals and JSON usually does not
//   - omitempty, which drops zero values from the output
//   - "-", which excludes a field entirely (the way to keep secrets out of logs)
//   - a validate tag, read by the custom validator below
type User struct {
	ID       int    `json:"id" validate:"required"`
	Name     string `json:"name" validate:"required,max=50"`
	Email    string `json:"email,omitempty"`
	Age      int    `json:"age,omitempty" validate:"min=0,max=150"`
	Password string `json:"-"`
}

// Only exported fields are visible to reflection and therefore to encoding/json. An
// unexported field is silently skipped — a common cause of "my JSON is empty".

// Temperature demonstrates custom marshalling. Implementing MarshalJSON/UnmarshalJSON lets a
// type control its own wire format without the caller knowing.
type Temperature float64

func (t Temperature) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(fmt.Sprintf("%.1f°C", float64(t)))), nil
}

func (t *Temperature) UnmarshalJSON(data []byte) error {
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return fmt.Errorf("temperature: %w", err)
	}

	f, err := strconv.ParseFloat(strings.TrimSuffix(s, "°C"), 64)
	if err != nil {
		return fmt.Errorf("temperature %q: %w", s, err)
	}

	*t = Temperature(f)
	return nil
}

// Reading is a struct embedding a custom-marshalled type.
type Reading struct {
	Station string      `json:"station"`
	Value   Temperature `json:"value"`
}

// ---------------------------------------------------------------------------
// Reading tags with reflect
// ---------------------------------------------------------------------------

// FieldInfo is one field's reflected metadata.
type FieldInfo struct {
	Name     string
	Type     string
	JSONName string
	Exported bool
}

// Describe walks a struct's fields. reflect.TypeOf gives the static type; Elem() steps
// through a pointer. Doing this once at startup and caching the result is how real
// libraries keep reflection off the hot path.
func Describe(v any) ([]FieldInfo, error) {
	t := reflect.TypeOf(v)
	if t == nil {
		return nil, fmt.Errorf("describe: nil value")
	}
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("describe: want a struct, got %s", t.Kind())
	}

	out := make([]FieldInfo, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		jsonName := f.Name
		if tag, ok := f.Tag.Lookup("json"); ok {
			// the tag is "name,opt1,opt2" — take the name and drop the options
			if name := strings.Split(tag, ",")[0]; name != "" {
				jsonName = name
			}
		}

		out = append(out, FieldInfo{
			Name:     f.Name,
			Type:     f.Type.String(),
			JSONName: jsonName,
			Exported: f.IsExported(),
		})
	}
	return out, nil
}

// ---------------------------------------------------------------------------
// A tag-driven validator
// ---------------------------------------------------------------------------

// Validate enforces the `validate` tags. This is a miniature of what go-playground/validator
// does, and shows the standard shape: reflect over fields, parse the tag, check the value.
func Validate(v any) error {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return fmt.Errorf("validate: nil pointer")
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("validate: want a struct, got %s", rv.Kind())
	}

	rt := rv.Type()
	var problems []string

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if !field.IsExported() {
			continue // reflection cannot read unexported values anyway
		}

		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}

		value := rv.Field(i)
		for _, rule := range strings.Split(tag, ",") {
			if err := applyRule(field.Name, value, rule); err != nil {
				problems = append(problems, err.Error())
			}
		}
	}

	if len(problems) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(problems, "; "))
	}
	return nil
}

// applyRule checks one rule against one field value.
func applyRule(name string, value reflect.Value, rule string) error {
	key, arg, hasArg := strings.Cut(rule, "=")

	switch key {
	case "required":
		if value.IsZero() {
			return fmt.Errorf("%s is required", name)
		}

	case "min", "max":
		if !hasArg {
			return fmt.Errorf("%s: rule %q needs an argument", name, key)
		}
		bound, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			return fmt.Errorf("%s: bad bound %q", name, arg)
		}

		// min/max means length for strings and magnitude for numbers
		var actual int64
		switch value.Kind() {
		case reflect.String:
			actual = int64(len(value.String()))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			actual = value.Int()
		default:
			return nil // rule does not apply to this kind
		}

		if key == "min" && actual < bound {
			return fmt.Errorf("%s must be at least %d", name, bound)
		}
		if key == "max" && actual > bound {
			return fmt.Errorf("%s must be at most %d", name, bound)
		}
	}

	return nil
}

// ---------------------------------------------------------------------------
// Setting values through reflection
// ---------------------------------------------------------------------------

// SetDefaults fills zero-valued fields from their `default` tag. Note the two conditions
// reflect.Value.Set needs: the value must come from a pointer (so it is addressable) and
// the field must be exported (so it is settable). CanSet reports both at once.
func SetDefaults(ptr any) error {
	rv := reflect.ValueOf(ptr)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("setDefaults: want a non-nil pointer, got %T", ptr)
	}
	rv = rv.Elem()

	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		def, ok := field.Tag.Lookup("default")
		if !ok {
			continue
		}

		value := rv.Field(i)
		if !value.CanSet() || !value.IsZero() {
			continue // never clobber a value the caller set
		}

		switch value.Kind() {
		case reflect.String:
			value.SetString(def)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(def, 10, 64)
			if err != nil {
				return fmt.Errorf("setDefaults %s: %w", field.Name, err)
			}
			value.SetInt(n)
		case reflect.Bool:
			b, err := strconv.ParseBool(def)
			if err != nil {
				return fmt.Errorf("setDefaults %s: %w", field.Name, err)
			}
			value.SetBool(b)
		}
	}

	return nil
}

// Config exercises SetDefaults.
type Config struct {
	Host    string `default:"localhost"`
	Port    int    `default:"8080"`
	Debug   bool   `default:"true"`
	Timeout int    // no tag: left alone
}

// ---------------------------------------------------------------------------
// Round-tripping
// ---------------------------------------------------------------------------

// ToJSON marshals with indentation.
func ToJSON(v any) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal: %w", err)
	}
	return string(b), nil
}
