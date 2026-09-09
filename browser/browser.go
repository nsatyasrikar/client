package browser

import "github.com/nsatyasrikar/client"

type documentAPI struct{}

var Document documentAPI

type elementAPI struct{}

var ElementAPI elementAPI

func (documentAPI) GetElementByID(id client.Value[client.StringType]) client.Value[client.Nullable[client.Element]] {
	doc := client.UnsafeExpression("document")
	method := client.Property[client.Unknown, client.Unknown](doc, "getElementById")
	return client.CallValue[client.Unknown, client.Nullable[client.Element]](method, id)
}

func (elementAPI) SetTextContent(el client.Value[client.Element], value client.Value[client.StringType]) client.Statement {
	return client.Assign(client.Property[client.Element, client.StringType](el, "textContent"), value)
}

func (elementAPI) AddEventListener(el client.Value[client.Element], name client.Value[client.StringType], callback client.Value[client.FunctionType]) client.Statement {
	method := client.Property[client.Element, client.Unknown](el, "addEventListener")
	call := client.CallValue[client.Unknown, client.Unknown](method, name, callback)
	return client.ExpressionStatement(call)
}
