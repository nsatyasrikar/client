package client

import (
	"fmt"
	"strings"
)

type emitter struct {
	b      strings.Builder
	indent int
	temp   int
}

func TypeScript(p ProgramNode) (string, error) {
	d := Validate(p)
	for _, x := range d {
		if x.Severity == SeverityError {
			return "", d
		}
	}
	e := &emitter{}
	for _, s := range p.statements {
		e.statement(s.node)
	}
	return e.b.String(), nil
}
func (e *emitter) line(s string) {
	e.b.WriteString(strings.Repeat("  ", e.indent))
	e.b.WriteString(s)
	e.b.WriteByte('\n')
}
func (e *emitter) statement(s statementNode) {
	switch n := s.(type) {
	case declarationStmt:
		e.line(n.kind + " " + n.name + " = " + e.expr(n.value) + ";")
	case assignStmt:
		e.line(e.expr(n.target) + " = " + e.expr(n.value) + ";")
	case exprStmt:
		e.line(e.expr(n.value) + ";")
	case rawStatement:
		e.line(n.code)
	case returnStmt:
		e.line("return " + e.expr(n.value) + ";")
	case blockStmt:
		e.line("{")
		e.indent++
		for _, x := range n.body {
			e.statement(x.node)
		}
		e.indent--
		e.line("}")
	case ifStmt:
		e.line("if (" + e.expr(n.test) + ") {")
		e.indent++
		for _, x := range n.yes {
			e.statement(x.node)
		}
		e.indent--
		if len(n.no) == 0 {
			e.line("}")
		} else {
			e.line("} else {")
			e.indent++
			for _, x := range n.no {
				e.statement(x.node)
			}
			e.indent--
			e.line("}")
		}
	case ifNotNullStmt:
		name := "value" + itoa(e.temp)
		e.temp++
		e.line("const " + name + " = " + e.expr(n.value) + ";")
		e.line("if (" + name + " !== null) {")
		e.indent++
		for _, x := range n.body {
			e.statement(replaceNarrow(x.node, name))
		}
		e.indent--
		e.line("}")
	}
}
func replaceNarrow(s statementNode, name string) statementNode {
	switch n := s.(type) {
	case exprStmt:
		return exprStmt{replaceExpr(n.value, name)}
	case assignStmt:
		return assignStmt{replaceExpr(n.target, name), replaceExpr(n.value, name)}
	case returnStmt:
		return returnStmt{replaceExpr(n.value, name)}
	}
	return s
}
func replaceExpr(x expression, name string) expression {
	if l, ok := x.(literalExpr); ok && l.code == "__narrowed" {
		return literalExpr{name}
	}
	return x
}
func (e *emitter) expr(x expression) string {
	switch n := x.(type) {
	case literalExpr:
		return n.code
	case rawExpr:
		return n.code
	case propertyExpr:
		sep := "."
		if n.optional {
			sep = "?."
		}
		return e.expr(n.target) + sep + n.name
	case indexExpr:
		return e.expr(n.target) + "[" + e.expr(n.index) + "]"
	case callExpr:
		a := make([]string, len(n.args))
		for i, x := range n.args {
			a[i] = e.expr(x)
		}
		return e.expr(n.callee) + "(" + strings.Join(a, ", ") + ")"
	case unaryExpr:
		return n.op + e.expr(n.value)
	case binaryExpr:
		return e.expr(n.left) + " " + n.op + " " + e.expr(n.right)
	case conditionalExpr:
		return e.expr(n.test) + " ? " + e.expr(n.yes) + " : " + e.expr(n.no)
	case arrayExpr:
		a := make([]string, len(n.values))
		for i, x := range n.values {
			a[i] = e.expr(x)
		}
		return "[" + strings.Join(a, ", ") + "]"
	case objectExpr:
		a := make([]string, len(n.entries))
		for i, x := range n.entries {
			a[i] = fmt.Sprintf("%s: %s", x.key, e.expr(x.value))
		}
		return "{" + strings.Join(a, ", ") + "}"
	case templateExpr:
		a := make([]string, len(n.parts))
		for i, x := range n.parts {
			a[i] = e.expr(x)
		}
		return strings.Join(a, " + ")
	}
	return "undefined"
}
