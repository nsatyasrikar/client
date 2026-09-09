package client

import "testing"

func emitOrFatal(t *testing.T, p ProgramNode) string {
	t.Helper()
	out, err := TypeScript(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return out
}

func TestCallValueUnwrapsArgs(t *testing.T) {
	fn := UnsafeExpression("foo")
	arg := String("bar")
	result := CallValue[Unknown, Unknown](fn, arg)
	out := emitOrFatal(t, Program(ExpressionStatement(result)))
	want := "foo(\"bar\");\n"
	if out != want {
		t.Fatalf("got %q, want %q", out, want)
	}
}

func TestJoinUnwrapsParts(t *testing.T) {
	result := Join(String("a"), String("b"))
	out := emitOrFatal(t, Program(ExpressionStatement(result)))
	want := "\"a\" + \"b\";\n"
	if out != want {
		t.Fatalf("got %q, want %q", out, want)
	}
}

func TestArrayOfUnwrapsValues(t *testing.T) {
	result := ArrayOf(String("a"), String("b"))
	out := emitOrFatal(t, Program(ExpressionStatement(result)))
	want := "[\"a\", \"b\"];\n"
	if out != want {
		t.Fatalf("got %q, want %q", out, want)
	}
}

func TestObjectUnwrapsValues(t *testing.T) {
	result := Object([2]any{"k", String("v")})
	out := emitOrFatal(t, Program(ExpressionStatement(result)))
	want := "{k: \"v\"};\n"
	if out != want {
		t.Fatalf("got %q, want %q", out, want)
	}
}
