package common

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

type Frontmatter struct {
	Fields map[string]FieldMeta `yaml:"fields"`
}

type FieldMeta struct {
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
}

var frontmatterDelimiter = []byte("---\n")

func ParseFrontmatter(data []byte) (Frontmatter, []byte, error) {
	if !bytes.HasPrefix(data, frontmatterDelimiter) {
		return Frontmatter{}, data, nil
	}

	rest := data[len(frontmatterDelimiter):]
	end := bytes.Index(rest, frontmatterDelimiter)
	if end == -1 {
		return Frontmatter{}, data, nil
	}

	yamlData := rest[:end]
	body := rest[end+len(frontmatterDelimiter):]

	var fm Frontmatter
	if err := yaml.Unmarshal(yamlData, &fm); err != nil {
		return Frontmatter{}, nil, fmt.Errorf("parsing frontmatter: %w", err)
	}

	return fm, body, nil
}
