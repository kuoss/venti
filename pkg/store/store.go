package store

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
)

func loadYaml(r io.Reader, y interface{}) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("cannot ReadFile: %w", err)
	}
	if err := yaml.Unmarshal(b, y); err != nil {
		return fmt.Errorf("cannot Unmarshal: %w", err)
	}
	return nil
}
