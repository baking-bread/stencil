package engine

import (
	"text/template/parse"

	"github.com/baking-bread/stencil/pkg/common"
)

type Analyzer struct {
	fields map[string]struct{}
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

func (a *Analyzer) AnalyzeTemplate(template common.Template) ([]string, error) {

	tree, err := parse.Parse(template.Name, string(template.Data), "{{", "}}", nil)
	if err != nil {
		return nil, err
	}

	for _, t := range tree {
		a.walkNode(t.Root)
	}

	result := make([]string, 0, len(a.fields))
	for field := range a.fields {
		result = append(result, field)
	}

	return result, nil
}

func (a *Analyzer) walkNode(node parse.Node) {
	if node == nil {
		return
	}

	switch n := node.(type) {
	case *parse.ActionNode:
		a.walkNode(n.Pipe)

	case *parse.IfNode:
		a.walkNode(n.Pipe)
		a.walkNode(n.List)
		a.walkNode(n.ElseList)

	case *parse.RangeNode:
		a.walkNode(n.Pipe)
		a.walkNode(n.List)
		a.walkNode(n.ElseList)

	case *parse.ListNode:
		if n == nil {
			return
		}
		for _, node := range n.Nodes {
			a.walkNode(node)
		}

	case *parse.PipeNode:
		for _, cmd := range n.Cmds {
			a.walkNode(cmd)
		}

	case *parse.CommandNode:
		for _, arg := range n.Args {
			a.walkNode(arg)
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
			a.fields[fieldPath] = struct{}{}
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
					a.fields[fieldPath] = struct{}{}
				}
			}
		}

	case *parse.TemplateNode:
		a.walkNode(n.Pipe)
	}
}
