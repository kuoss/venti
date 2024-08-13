package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name string `yaml:"name"`
	Age  int    `yaml:"age"`
}

func TestUnmarshalStrict_ValidYAML(t *testing.T) {
	data := []byte(`
name: John Doe
age: 30
`)
	var ts TestStruct
	err := UnmarshalStrict(data, &ts)
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", ts.Name)
	assert.Equal(t, 30, ts.Age)
}

func TestUnmarshalStrict_InvalidYAML_UnknownField(t *testing.T) {
	data := []byte(`
name: John Doe
age: 30
city: New York
`)
	var ts TestStruct
	err := UnmarshalStrict(data, &ts)
	assert.EqualError(t, err, "yaml: unmarshal errors:\n  line 4: field city not found in type util.TestStruct")
}

func TestUnmarshalStrict_InvalidYAML_InvalidType(t *testing.T) {

	t.Run("Invalid YAML - Unknown Field", func(t *testing.T) {
		data := []byte(`
name: John Doe
age: 30
city: New York
`)
		var ts TestStruct
		err := UnmarshalStrict(data, &ts)
		assert.EqualError(t, err, "yaml: unmarshal errors:\n  line 4: field city not found in type util.TestStruct")
	})

	t.Run("Invalid YAML - Invalid Type", func(t *testing.T) {
		data := []byte(`
name: John Doe
age: "thirty"
`)
		var ts TestStruct
		err := UnmarshalStrict(data, &ts)
		assert.EqualError(t, err, "yaml: unmarshal errors:\n  line 3: cannot unmarshal !!str `thirty` into int")
	})
}
