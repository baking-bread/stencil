package engine

import (
	"regexp"
	"testing"

	"github.com/baking-bread/stencil/pkg/common"
)

func TestRenderSimpleTemplateWithoutValues(t *testing.T) {

	var expected = map[string]string{
		"Test": "Hello World!",
	}

	var templates = []common.Template{
		{
			Name: "Test",
			Data: []byte("Hello World!"),
		},
	}

	var values = common.Values{}

	var output, _ = Render(templates, values)

	for name := range expected {
		if output[name] != expected[name] {
			t.Fail()
		}
	}
}

func TestRenderSimpleTemplateWithSingleValue(t *testing.T) {

	var expected = map[string]string{
		"Test": "Hello World!",
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

	var output, _ = Render(templates, values)

	for name := range expected {
		if output[name] != expected[name] {
			t.Fail()
		}
	}
}

func TestRenderRandomStringFromArrayValues(t *testing.T) {

	var expected = map[string]string{
		"Test": "Hello (World|Steve)!",
	}

	var templates = []common.Template{
		{
			Name: "Test",
			Data: []byte("Hello {{ pick .test (random (length .test)) }}!"),
		},
	}

	var values, _ = common.ReadValues(
		[]byte("test: [\"World\", \"Steve\"]"),
	)

	var output, err = Render(templates, values)

	if err != nil {
		t.Errorf("%e", err)
	}

	for name := range expected {
		if match, _ := regexp.Match(expected[name], []byte(output[name])); !match {
			t.Errorf("Expected %s; Got %s", expected[name], output[name])
		}
	}
}

func TestRenderTemplateWithFrontmatter(t *testing.T) {

	var templates, err = common.LoadText("Test", "---\nfields:\n  name:\n    type: string\n---\nHello {{ .name }}!")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var values, _ = common.ReadValues(
		[]byte("name: \"World\""),
	)

	output, err := Render(templates, values)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Hello World!"
	if output["Test"] != expected {
		t.Errorf("expected %q, got %q", expected, output["Test"])
	}
}

func TestRenderTimestamp(t *testing.T) {

	var expected = map[string]string{
		"Test": "Hello .*!",
	}

	var templates = []common.Template{
		{
			Name: "Test",
			Data: []byte("Hello {{ timestamp \"2006-01-02\" }}!"),
		},
	}

	var values = common.Values{}

	var output, err = Render(templates, values)

	if err != nil {
		t.Errorf("%e", err)
	}

	for name := range expected {

		if match, _ := regexp.Match(expected[name], []byte(output[name])); !match {
			t.Errorf("Expected %s; Got %s", expected[name], output[name])
		}
	}
}
