package util

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

func UnmarshalStrict(data []byte, out any) error {
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true) // Disallow unknown fields
	return decoder.Decode(out)
}
