package engine

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/baking-bread/stencil/pkg/common"
	"github.com/baking-bread/stencil/pkg/funcs"
)

type Engine struct {
	CustomTemplateFuncs template.FuncMap
}

func New() Engine {
	return Engine{
		CustomTemplateFuncs: nil,
	}
}

func (e Engine) Render(templates []common.Template, values common.Values) (map[string]string, error) {
	return e.render(templates, values)
}

func Render(templates []common.Template, values common.Values) (map[string]string, error) {
	return new(Engine).Render(templates, values)
}

func (e Engine) render(templates []common.Template, values common.Values) (rendered map[string]string, err error) {

	defer func() {
		if template := recover(); template != nil {
			err = fmt.Errorf("rendering template failed: %v", template)
		}
	}()

	root := template.New("gotpl")

	e.initFuncMap(root)

	for _, template := range templates {
		if _, err := root.New(template.Name).Parse(string(template.Data)); err != nil {
			return map[string]string{}, err
		}
	}

	rendered = make(map[string]string, len(templates))
	for _, template := range templates {

		var buffer strings.Builder
		if err := root.ExecuteTemplate(&buffer, template.Name, values); err != nil {
			return map[string]string{}, err
		}

		rendered[template.Name] = buffer.String()
	}

	return rendered, nil
}

func (e Engine) initFuncMap(t *template.Template) {
	funcMap := template.FuncMap{
		"random":    funcs.Random,
		"timestamp": funcs.Timestamp,
		"upper":     funcs.Upper,
		"lower":     funcs.Lower,
		"length":    funcs.Length,
		"pick":      funcs.Pick,
	}

	t.Funcs(funcMap)
}
