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
func ExpressionStatement(value Expression) Statement { return Statement{node: exprStmt{value: value}} }
func Block(statements ...Statement) Statement {
	return Statement{node: blockStmt{body: append([]Statement(nil), statements...)}}
}
