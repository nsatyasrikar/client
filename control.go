package client

type ifStmt struct {
	test    expression
	yes, no []Statement
}

func (ifStmt) stmt() {}

type ifNotNullStmt struct {
	value expression
	body  []Statement
}

func (ifNotNullStmt) stmt() {}

type returnStmt struct{ value expression }

func (returnStmt) stmt() {}

func If(test Value[BooleanType], body ...Statement) Statement {
	return Statement{node: ifStmt{test.node, append([]Statement(nil), body...), nil}}
}
func IfElse(test Value[BooleanType], yes []Statement, no []Statement) Statement {
	return Statement{node: ifStmt{test.node, append([]Statement(nil), yes...), append([]Statement(nil), no...)}}
}
func IfNotNull[T Type](value Value[Nullable[T]], callback func(Value[T]) Statement) Statement {
	narrowed := wrap[T](literalExpr{"__narrowed"})
	return Statement{node: ifNotNullStmt{value.node, []Statement{callback(narrowed)}}}
}
func Return[T Type](value Value[T]) Statement { return Statement{node: returnStmt{value.node}} }
