package client

type Type interface{ clientType() }
type Unknown struct{}

func (Unknown) clientType() {}

type StringType struct{}

func (StringType) clientType() {}

type NumberType struct{}

func (NumberType) clientType() {}

type BooleanType struct{}

func (BooleanType) clientType() {}

type NullType struct{}

func (NullType) clientType() {}

type UndefinedType struct{}

func (UndefinedType) clientType() {}

type FunctionType struct{}

func (FunctionType) clientType() {}

type Element struct{}

func (Element) clientType() {}

type Node struct{}

func (Node) clientType() {}

type Event struct{}

func (Event) clientType() {}

type Nullable[T Type] struct{}

func (Nullable[T]) clientType() {}

type Array[T Type] struct{}

func (Array[T]) clientType() {}

type Promise[T Type] struct{}

func (Promise[T]) clientType() {}

// TypeMarker lets types outside this package satisfy Type by embedding it.
type TypeMarker struct{}

func (TypeMarker) clientType() {}

type expression interface{ expr() }
type statementNode interface{ stmt() }
type Expression interface{ expression }

type Value[T Type] struct{ node expression }

func (Value[T]) expr() {}
func (v Value[T]) unwrapNode() expression { return v.node }

// valueNode is implemented by every Value[T]; unwrapExpr uses it to reach
// the underlying AST node when an expression arrives boxed in a Value[T],
// which happens whenever a caller passes one through an Expression-typed
// parameter (Call, CallValue, Join, ArrayOf, Object, ExpressionStatement).
type valueNode interface{ unwrapNode() expression }

func unwrapExpr(x expression) expression {
	if v, ok := x.(valueNode); ok {
		return v.unwrapNode()
	}
	return x
}

type Statement struct{ node statementNode }

func (Statement) stmt() {}

type ProgramNode struct{ statements []Statement }

func Program(statements ...Statement) ProgramNode {
	return ProgramNode{statements: append([]Statement(nil), statements...)}
}

type literalExpr struct{ code string }

func (literalExpr) expr() {}

type rawExpr struct{ code string }

func (rawExpr) expr() {}

type exprStmt struct{ value expression }

func (exprStmt) stmt() {}

type blockStmt struct{ body []Statement }

func (blockStmt) stmt() {}

func wrap[T Type](n expression) Value[T]             { return Value[T]{node: n} }
func ExpressionStatement(value Expression) Statement {
	return Statement{node: exprStmt{value: unwrapExpr(value)}}
}
func Block(statements ...Statement) Statement {
	return Statement{node: blockStmt{body: append([]Statement(nil), statements...)}}
}
