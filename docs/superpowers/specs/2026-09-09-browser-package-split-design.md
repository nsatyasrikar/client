# Design: split generated browser bindings into their own package

## Problem

Repo root has 1049 flat `.go` files: 15 hand-written AST/compiler files and
1034 generated `browser_*.go` binding stubs (one per HTML/DOM attribute,
written by `cmd/clientgen`). Root is unnavigable.

## Constraint discovered during brainstorming

`expression`, `Statement`, `Type` are Go's classic sealed-sum-type pattern:
interfaces with unexported methods (`expr()`, `stmt()`, `clientType()`).
`emit.go` and `validate.go` type-switch directly on unexported concrete
structs (`callExpr`, `assignStmt`, `binaryExpr`, ...). Go has no
same-package-multiple-directories mechanism, so a real submodule needs a
real package boundary — for the 15 core files that means exporting ~20
internal node types, a public-API redesign disproportionate to a
navigation cleanup. Scope is narrowed accordingly (user-approved):
only the generated-bindings layer moves.

## Change

1. **New package `browser/`** (`package browser`) holds:
   - all 1034 `browser_*.go` generated stub files
   - the hand-written DOM API wrappers currently in root `browser.go`
     (`documentAPI`, `elementAPI`, `Document`, `ElementAPI`)

2. **Sealing fix** — `type.go` gains:
   ```go
   type TypeMarker struct{}
   func (TypeMarker) clientType() {}
   ```
   Generated stubs become `type ActionType struct{ client.TypeMarker }`
   (method promotion satisfies `client.Type` without a per-file method).
   `generator.go`'s template updates to emit `package browser`,
   `import "github.com/nsatyasrikar/client"`, and the embed form.

3. **browser.go rewrite** — `GetElementByID`/`SetTextContent`/
   `AddEventListener` currently construct unexported AST nodes
   (`callExpr`, `propertyExpr`, `literalExpr`) directly. Rewrite using only
   already-exported `client.Property`, `client.CallValue`, `client.Assign`,
   `client.ExpressionStatement`, `client.UnsafeExpression("document")`.
   Same emitted output, zero new exports needed on the `client` side beyond
   `TypeMarker`.

4. **Import fixups** — anything referencing `client.Document`,
   `client.ElementAPI`, or any `client.XType` switches to `browser.`.
   No external consumers yet (module has no downstream importers), so
   blast radius is internal: `cmd/clientgen`, any `_test.go` files, and
   root package files that reference browser bindings.

5. **Knowledge graph** — after the move commit, and per future commit
   touching `client`/`browser`, run `graphify update .`. Install
   graphify's post-commit git hook so this is automatic rather than
   manual.

## Out of scope

- Splitting `ast`/`compiler`/`emit`/`validate` into separate packages
  (would require exporting internal sum-type nodes — deferred, not
  requested).
- Any behavior change to emitted TypeScript/JS output.

## Testing

- `go build ./...` and `go vet ./...` after the move.
- Existing tests (if any) pass unchanged in behavior — this is a pure
  move + mechanical rewrite, no semantic change to generated code.
- `cmd/clientgen -check` (freshness check flag already in generator)
  confirms regenerated files match committed files.
