package client

func Callback(body ...Statement) Value[FunctionType] {
	return wrap[FunctionType](rawExpr{"(event) => {" + emitBlock(body) + "}"})
}
func emitBlock(body []Statement) string { p := Program(body...); s, _ := TypeScript(p); return s }
