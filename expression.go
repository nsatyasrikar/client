package client

import "strings"

type propertyExpr struct {
	target   expression
	name     string
	optional bool
}

func (propertyExpr) expr() {}

type indexExpr struct{ target, index expression }

func (indexExpr) expr() {}

type callExpr struct {
	callee expression
	args   []expression
}

func (callExpr) expr() {}

type unaryExpr struct {
	op    string
	value expression
}

func (unaryExpr) expr() {}

type binaryExpr struct {
	op          string
	left, right expression
}

func (binaryExpr) expr() {}

type conditionalExpr struct{ test, yes, no expression }

func (conditionalExpr) expr() {}

type templateExpr struct{ parts []expression }

func (templateExpr) expr() {}

func Property[T, R Type](target Value[T], name string) Value[R] {
	return wrap[R](propertyExpr{target.node, name, false})
}
func OptionalProperty[T, R Type](target Value[Nullable[T]], name string) Value[Nullable[R]] {
	return wrap[Nullable[R]](propertyExpr{target.node, name, true})
}
func Index[T, R Type](target Value[T], index Value[NumberType]) Value[R] {
	return wrap[R](indexExpr{target.node, index.node})
}
func Call[R Type](callee Value[Unknown], args ...Expression) Value[R] {
	nodes := make([]expression, len(args))
	for i, a := range args {
		nodes[i] = unwrapExpr(a)
	}
	return wrap[R](callExpr{callee.node, nodes})
}
func CallValue[T, R Type](callee Value[T], args ...Expression) Value[R] {
	nodes := make([]expression, len(args))
	for i, a := range args {
		nodes[i] = unwrapExpr(a)
	}
	return wrap[R](callExpr{callee.node, nodes})
}
func Unary[T Type](op string, value Value[T]) Value[T] { return wrap[T](unaryExpr{op, value.node}) }
func Binary[T, R Type](op string, left Value[T], right Value[R]) Value[BooleanType] {
	return wrap[BooleanType](binaryExpr{op, left.node, right.node})
}
func Conditional[T Type](test Value[BooleanType], yes, no Value[T]) Value[T] {
	return wrap[T](conditionalExpr{test.node, yes.node, no.node})
}
func Equal[T Type](left, right Value[T]) Value[BooleanType] { return Binary[T, T]("===", left, right) }
func NotEqual[T Type](left, right Value[T]) Value[BooleanType] {
	return Binary[T, T]("!==", left, right)
}
func And(left, right Value[BooleanType]) Value[BooleanType] {
	return Binary[BooleanType, BooleanType]("&&", left, right)
}
func Or(left, right Value[BooleanType]) Value[BooleanType] {
	return Binary[BooleanType, BooleanType]("||", left, right)
}
func Coalesce[T Type](left Value[Nullable[T]], right Value[T]) Value[T] {
	return wrap[T](binaryExpr{"??", left.node, right.node})
}
func Join(parts ...Expression) Value[StringType] {
	nodes := make([]expression, len(parts))
	for i, p := range parts {
		nodes[i] = unwrapExpr(p)
	}
	return wrap[StringType](templateExpr{nodes})
}
func validIdentifier(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if !(r == '_' || r == '$' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || i > 0 && r >= '0' && r <= '9') {
			return false
		}
	}
	return !strings.Contains("class function const let var return if else for while", s)
}
