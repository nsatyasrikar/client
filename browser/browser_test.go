package browser

import (
	"testing"

	"github.com/nsatyasrikar/client"
)

func TestGetElementByIDEmitsExpectedCall(t *testing.T) {
	id := client.String("my-id")
	result := Document.GetElementByID(id)
	stmt := client.ExpressionStatement(result)
	out, err := client.TypeScript(client.Program(stmt))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "document.getElementById(\"my-id\");\n"
	if out != want {
		t.Fatalf("got %q, want %q", out, want)
	}
}

// Compile-time check that SetTextContent and AddEventListener keep their
// signatures after the move. A Value[Element] can only be constructed
// inside package client today (no exported constructor produces one
// standalone), so these two aren't independently runtime-testable from
// outside without a new export — out of scope for this move.
var (
	_ = ElementAPI.SetTextContent
	_ = ElementAPI.AddEventListener
)
