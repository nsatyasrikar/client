package client

type declarationStmt struct {
	kind, name string
	value      expression
}

func (declarationStmt) stmt() {}

type assignStmt struct{ target, value expression }

func (assignStmt) stmt() {}
func Const[T Type](name string, value Value[T]) (Value[T], Statement) {
	return wrap[T](literalExpr{name}), Statement{node: declarationStmt{"const", name, value.node}}
}
func Let[T Type](name string, value Value[T]) (Value[T], Statement) {
	return wrap[T](literalExpr{name}), Statement{node: declarationStmt{"let", name, value.node}}
}
func Assign[T Type](target, value Value[T]) Statement {
	return Statement{node: assignStmt{target.node, value.node}}
}
func UnsafeStatement(code string) Statement { return Statement{node: rawStatement{code}} }

type rawStatement struct{ code string }

func (rawStatement) stmt() {}
