package client

func Validate(p ProgramNode) Diagnostics {
	var out Diagnostics
	for i, s := range p.statements {
		validateStatement(s.node, "/program["+itoa(i)+"]", &out)
	}
	return out
}
func validateStatement(s statementNode, path string, out *Diagnostics) {
	switch n := s.(type) {
	case declarationStmt:
		if !validIdentifier(n.name) {
			*out = append(*out, Diagnostic{"invalid-identifier", "invalid declaration identifier", path, SeverityError})
		}
		validateExpr(n.value, path+"/value", out)
	case assignStmt:
		validateExpr(n.target, path+"/target", out)
		validateExpr(n.value, path+"/value", out)
	case exprStmt:
		validateExpr(n.value, path, out)
	case blockStmt:
		for i, x := range n.body {
			validateStatement(x.node, path+"/body["+itoa(i)+"]", out)
		}
	case ifStmt:
		validateExpr(n.test, path+"/test", out)
		for i, x := range n.yes {
			validateStatement(x.node, path+"/then["+itoa(i)+"]", out)
		}
		for i, x := range n.no {
			validateStatement(x.node, path+"/else["+itoa(i)+"]", out)
		}
	case ifNotNullStmt:
		validateExpr(n.value, path+"/value", out)
		for i, x := range n.body {
			validateStatement(x.node, path+"/body["+itoa(i)+"]", out)
		}
	case rawStatement:
		*out = append(*out, Diagnostic{"unsafe-code", "unsafe statement emitted", path, SeverityWarning})
	case returnStmt:
		validateExpr(n.value, path+"/value", out)
	}
}
func validateExpr(e expression, path string, out *Diagnostics) {
	switch n := e.(type) {
	case rawExpr:
		*out = append(*out, Diagnostic{"unsafe-code", "unsafe expression emitted", path, SeverityWarning})
	case propertyExpr:
		validateExpr(n.target, path+"/target", out)
	case indexExpr:
		validateExpr(n.target, path+"/target", out)
		validateExpr(n.index, path+"/index", out)
	case callExpr:
		validateExpr(n.callee, path+"/callee", out)
		for i, a := range n.args {
			validateExpr(a, path+"/args["+itoa(i)+"]", out)
		}
	case binaryExpr:
		validateExpr(n.left, path+"/left", out)
		validateExpr(n.right, path+"/right", out)
	case unaryExpr:
		validateExpr(n.value, path+"/value", out)
	case conditionalExpr:
		validateExpr(n.test, path+"/test", out)
		validateExpr(n.yes, path+"/yes", out)
		validateExpr(n.no, path+"/no", out)
	case arrayExpr:
		for i, a := range n.values {
			validateExpr(a, path+"/items["+itoa(i)+"]", out)
		}
	case objectExpr:
		for i, a := range n.entries {
			validateExpr(a.value, path+"/entries["+itoa(i)+"]", out)
		}
	case templateExpr:
		for i, a := range n.parts {
			validateExpr(a, path+"/parts["+itoa(i)+"]", out)
		}
	}
}
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	b := []byte{}
	for i > 0 {
		b = append([]byte{byte('0' + i%10)}, b...)
		i /= 10
	}
	return string(b)
}
