package util

import (
	"testing"
)

type TestStruct struct {
	Name string `yaml:"name"`
	Age  int    `yaml:"age"`
}

func TestUnmarshalStrict(t *testing.T) {
	t.Run("Valid YAML", func(t *testing.T) {
		data := []byte(`
name: John Doe
age: 30
`)
		var ts TestStruct
		err := UnmarshalStrict(data, &ts)
		if err != nil {
			t.Errorf("UnmarshalStrict failed: %v", err)
		}
		if ts.Name != "John Doe" || ts.Age != 30 {
			t.Errorf("UnmarshalStrict failed: expected name: John Doe, age: 30, got: name: %s, age: %d", ts.Name, ts.Age)
		}
	})

	t.Run("Invalid YAML - Unknown Field", func(t *testing.T) {
		data := []byte(`
name: John Doe
age: 30
city: New York
`)
		var ts TestStruct
		err := UnmarshalStrict(data, &ts)
		if err == nil {
			t.Errorf("UnmarshalStrict should have failed due to unknown field 'city'")
		}
	})

	t.Run("Invalid YAML - Invalid Type", func(t *testing.T) {
		data := []byte(`
name: John Doe
age: "thirty"
`)
		var ts TestStruct
		err := UnmarshalStrict(data, &ts)
		if err == nil {
			t.Errorf("UnmarshalStrict should have failed due to invalid type for 'age'")
		}
	})
}
