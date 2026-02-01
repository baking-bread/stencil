package engine

import (
	"text/template/parse"

	"github.com/baking-bread/stencil/pkg/common"
)

func (e Engine) Analyze(templates []common.Template, values common.Values) (map[string]bool, error) {
	return e.analyze(templates, values)
}

func Analyze(templates []common.Template, values common.Values) (map[string]bool, error) {
	return new(Engine).Analyze(templates, values)
}

func (e Engine) analyze(templates []common.Template, values common.Values) (map[string]bool, error) {

	fields := make(map[string]bool)

	for _, template := range templates {
		e.analyzeTemplate(template, &fields)
	}

	for field := range fields {
		check := common.Pick(values, field)
		if check != nil {
			fields[field] = true
		}
	}

	return fields, nil
}

func (e Engine) analyzeTemplate(template common.Template, fields *map[string]bool) (map[string]bool, error) {
	result := make(map[string]bool)

	tree, err := parse.Parse(template.Name, string(template.Data), "{{", "}}", nil)
	if err != nil {
		return nil, err
	}

	for _, t := range tree {
		e.walkNode(t.Root, fields)
	}

	return result, nil
}

func (e Engine) walkNode(node parse.Node, fields *map[string]bool) {
	if node == nil {
		return
	}

	_fields := *fields

	switch n := node.(type) {
	case *parse.ActionNode:
		e.walkNode(n.Pipe, fields)

	case *parse.IfNode:
		e.walkNode(n.Pipe, fields)
		e.walkNode(n.List, fields)
		e.walkNode(n.ElseList, fields)

	case *parse.RangeNode:
		e.walkNode(n.Pipe, fields)
		e.walkNode(n.List, fields)
		e.walkNode(n.ElseList, fields)

	case *parse.ListNode:
		if n == nil {
			return
		}
		for _, node := range n.Nodes {
			e.walkNode(node, fields)
		}

	case *parse.PipeNode:
		for _, cmd := range n.Cmds {
			e.walkNode(cmd, fields)
		}

	case *parse.CommandNode:
		for _, arg := range n.Args {
			e.walkNode(arg, fields)
		}

	case *parse.FieldNode:
		fieldPath := ""
		for i, ident := range n.Ident {
			if i > 0 {
				fieldPath += "."
			}
			fieldPath += ident
		}
		if fieldPath != "" {
			_fields[fieldPath] = false
		}

	case *parse.VariableNode:
		for i, ident := range n.Ident {
			if i > 0 {
				fieldPath := ""
				for j := 1; j <= i; j++ {
					if j > 1 {
						fieldPath += "."
					}
					fieldPath += ident
				}
				if fieldPath != "" {
					_fields[fieldPath] = false
				}
			}
		}

	case *parse.TemplateNode:
		e.walkNode(n.Pipe, fields)
	}
}
