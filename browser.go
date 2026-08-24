package client

type documentAPI struct{}

var Document documentAPI

type elementAPI struct{}

var ElementAPI elementAPI

type documentValue struct{}

func (documentAPI) GetElementByID(id Value[StringType]) Value[Nullable[Element]] {
	return wrap[Nullable[Element]](callExpr{propertyExpr{literalExpr{"document"}, "getElementById", false}, []expression{id.node}})
}
func (elementAPI) SetTextContent(el Value[Element], value Value[StringType]) Statement {
	return Assign(Property[Element, StringType](el, "textContent"), value)
}
func (elementAPI) AddEventListener(el Value[Element], name Value[StringType], callback Value[FunctionType]) Statement {
	return ExpressionStatement(wrap[Unknown](callExpr{propertyExpr{el.node, "addEventListener", false}, []expression{name.node, callback.node}}))
}
