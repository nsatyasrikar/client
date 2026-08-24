package client

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func String(v string) Value[StringType] {
	b, _ := json.Marshal(v)
	return wrap[StringType](literalExpr{string(b)})
}
func Number(v float64) Value[NumberType] {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		panic("client: non-finite number")
	}
	return wrap[NumberType](literalExpr{strconv.FormatFloat(v, 'f', -1, 64)})
}
func Boolean(v bool) Value[BooleanType]  { return wrap[BooleanType](literalExpr{strconv.FormatBool(v)}) }
func Null() Value[NullType]              { return wrap[NullType](literalExpr{"null"}) }
func NullOf[T Type]() Value[Nullable[T]] { return wrap[Nullable[T]](literalExpr{"null"}) }
func Undefined() Value[UndefinedType]    { return wrap[UndefinedType](literalExpr{"undefined"}) }
func BigInt(v int64) Value[Unknown] {
	return wrap[Unknown](literalExpr{strconv.FormatInt(v, 10) + "n"})
}
func Regex(pattern, flags string) Value[Unknown] {
	return wrap[Unknown](literalExpr{"/" + strings.ReplaceAll(pattern, "/", "\\/") + "/" + flags})
}
func Template(parts ...string) Value[StringType] {
	var b strings.Builder
	b.WriteByte('`')
	for _, p := range parts {
		b.WriteString(strings.ReplaceAll(strings.ReplaceAll(p, "\\", "\\\\"), "`", "\\`"))
	}
	b.WriteByte('`')
	return wrap[StringType](literalExpr{b.String()})
}
func ArrayOf[T Type](values ...Value[T]) Value[Array[T]] {
	nodes := make([]expression, len(values))
	for i, v := range values {
		nodes[i] = v
	}
	return wrap[Array[T]](arrayExpr{nodes})
}

type arrayExpr struct{ values []expression }

func (arrayExpr) expr() {}

type objectEntry struct {
	key   string
	value expression
}
type objectExpr struct{ entries []objectEntry }

func (objectExpr) expr() {}
func Object(entries ...[2]any) Value[Unknown] {
	out := make([]objectEntry, 0, len(entries))
	for _, e := range entries {
		k, ok := e[0].(string)
		if !ok {
			panic("client: object key must be string")
		}
		v, ok := e[1].(Expression)
		if !ok {
			panic("client: object value must be expression")
		}
		out = append(out, objectEntry{k, v})
	}
	return wrap[Unknown](objectExpr{out})
}
func UnsafeExpression(code string) Value[Unknown] {
	if strings.TrimSpace(code) == "" {
		panic(fmt.Sprintf("client: empty unsafe expression"))
	}
	return wrap[Unknown](rawExpr{code})
}
