package engine

import (
	"testing"

	"github.com/baking-bread/stencil/pkg/common"
)

func TestAnalyzeSimpleTemplateWithSingleValue(t *testing.T) {

	var expected = map[string]bool{
		"test": true,
	}

	var templates = []common.Template{
		{
			Name: "Test",
			Data: []byte("Hello {{ .test }}!"),
		},
	}

	var values, _ = common.ReadValues(
		[]byte("test: \"World\""),
	)

	var output, _ = Analyze(templates, values)

	for name := range expected {
		if output[name] != expected[name] {
			t.Fail()
		}
	}
}

func TestAnalyzeSimpleTemplateWithNestedValue(t *testing.T) {

	var expected = map[string]bool{
		"test.name": true,
	}

	var templates = []common.Template{
		{
			Name: "Test",
			Data: []byte("Hello {{ .test.name }}!"),
		},
	}

	var values, _ = common.ReadValues(
		[]byte("test:\n  name: \"World\""),
	)

	var output, _ = Analyze(templates, values)

	for name := range expected {
		if output[name] != expected[name] {
			t.Fail()
		}
	}
}
