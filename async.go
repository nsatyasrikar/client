package client

func Await[T Type](value Value[Promise[T]]) Value[T] {
	return wrap[T](rawExpr{"await " + (&emitter{}).expr(value.node)})
}
