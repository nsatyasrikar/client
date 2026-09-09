# client

A typed Go AST for building JavaScript/TypeScript expressions and statements,
compiled to a JS bundle via [esbuild](https://github.com/evanw/esbuild).

Every value in the AST carries its JS type as a Go generic parameter
(`Value[StringType]`, `Value[Element]`, ...), so an expression built with the
wrong type is a Go compile error, not a runtime one.

## Contents

- [Package layout](#package-layout)
- [Quick start](#quick-start)
- [Examples](#examples)
  - [Values and declarations](#values-and-declarations)
  - [Control flow](#control-flow)
  - [Conditional expressions](#conditional-expressions)
  - [Arrays and objects](#arrays-and-objects)
  - [Browser DOM access](#browser-dom-access)
  - [Compiling to a JS bundle](#compiling-to-a-js-bundle)
  - [Validation diagnostics](#validation-diagnostics)
- [Building and testing](#building-and-testing)
- [Limitations](#limitations)

## Package layout

- **`client`** (repo root) — the AST, type system, statement/expression
  builders, validator, TypeScript emitter, and the esbuild-backed compiler.
- **`browser/`** — DOM/browser API bindings: hand-written wrappers
  (`Document.GetElementByID`, `ElementAPI.SetTextContent`,
  `ElementAPI.AddEventListener`) plus 1034 static marker types
  (`browser/types_generated.go`) satisfying `client.Type` for browser
  interfaces (`ActionType`, `AlarmType`, ...) via an embedded
  `client.TypeMarker`.

## Quick start

```go
import "github.com/nsatyasrikar/client"

name, decl := client.Let("name", client.String("world"))
greet := client.Join(client.String("hello, "), name)
call := client.Call[client.Unknown](client.UnsafeExpression("console.log"), greet)
stmt := client.ExpressionStatement(call)

out, err := client.TypeScript(client.Program(decl, stmt))
// out == "let name = \"world\";\nconsole.log(\"hello, \" + name);\n"
```

Every example below is a real program: build a `ProgramNode` with
`client.Program(...statements)`, then either:

- `client.TypeScript(p)` — emit readable TypeScript/JS source, or
- `client.Compile(p)` / `client.CompileWithOptions(p, opts)` — run that
  source through esbuild for bundling, minification, and ESM/IIFE output.

## Examples

All output comments below were captured by actually running the snippet as a
Go test — none are hand-typed guesses.

### Values and declarations

`Let`/`Const` declare a variable and hand back a `Value[T]` referencing it:

```go
name, decl := client.Let("name", client.String("world"))
out, _ := client.TypeScript(client.Program(decl))
// out == "let name = \"world\";\n"
```

### Control flow

```go
x, decl := client.Let("x", client.Number(5))
cond := client.Binary[client.NumberType, client.NumberType](">", x, client.Number(3))
stmt := client.If(cond, client.ExpressionStatement(client.UnsafeExpression(`console.log("big")`)))

out, _ := client.TypeScript(client.Program(decl, stmt))
// out == "let x = 5;\nif (x > 3) {\n  console.log(\"big\");\n}\n"
```

`IfElse` takes explicit `yes`/`no` statement slices; `IfNotNull` narrows a
`Value[Nullable[T]]` to `Value[T]` inside its callback.

### Conditional expressions

```go
x, decl := client.Let("x", client.Number(10))
cond := client.Binary[client.NumberType, client.NumberType](">", x, client.Number(5))
result := client.Conditional(cond, client.String("big"), client.String("small"))
_, decl2 := client.Let("label", result)

out, _ := client.TypeScript(client.Program(decl, decl2))
// out == "let x = 10;\nlet label = x > 5 ? \"big\" : \"small\";\n"
```

### Arrays and objects

```go
arr := client.ArrayOf(client.Number(1), client.Number(2), client.Number(3))
_, decl := client.Let("nums", arr)

obj := client.Object([2]any{"name", client.String("client")}, [2]any{"stable", client.Boolean(true)})
_, decl2 := client.Let("info", obj)

out, _ := client.TypeScript(client.Program(decl, decl2))
// out == "let nums = [1, 2, 3];\nlet info = {name: \"client\", stable: true};\n"
```

### Browser DOM access

```go
import "github.com/nsatyasrikar/client/browser"

id := client.String("my-id")
call := client.ExpressionStatement(browser.Document.GetElementByID(id))

out, _ := client.TypeScript(client.Program(call))
// out == "document.getElementById(\"my-id\");\n"
```

### Compiling to a JS bundle

```go
_, decl := client.Let("greeting", client.String("hi"))
out, err := client.Compile(client.Program(decl))
// out == "let e=\"hi\";\n"  (esbuild's default minify renames the identifier)
```

`Compile` uses `CompileOptions{Format: FormatESM, Target: "es2020", Minify: true}`;
call `CompileWithOptions` directly to control format (`FormatESM`/`FormatIIFE`),
target, minification, and source maps.

### Validation diagnostics

`Validate` walks a program and reports problems without emitting anything;
`TypeScript`/`Compile` call it internally and fail only on `SeverityError`
(a `SeverityWarning`, like the one below, still emits):

```go
stmt := client.UnsafeStatement("debugger;")
diags := client.Validate(client.Program(stmt))
// diags == Diagnostics{{Code: "unsafe-code", Message: "unsafe statement emitted", ...}}
```

`UnsafeStatement`/`UnsafeExpression` exist as an escape hatch for code this
AST can't express yet; both are flagged in diagnostics precisely so that
escape hatch stays visible in review.

## Building and testing

```bash
go build ./...
go vet ./...
go test ./...
```

No API key, network access, or external data source is required.

## Limitations

- `module.go` and `class.go` are currently empty — module/class emission
  isn't implemented yet.
- `browser/types_generated.go` was originally produced by a `clientgen` tool
  from a JSON API schema. That generator and its schema data source were
  never part of this repo and have since been removed; the file is now
  static and hand-maintained. Adding a new browser API type means adding a
  `type FooType struct{ client.TypeMarker }` line by hand.
- Only a handful of DOM APIs have hand-written wrappers in `browser/browser.go`
  (`GetElementByID`, `SetTextContent`, `AddEventListener`) — the 1034
  generated types exist so the type system recognizes those interfaces as
  `client.Type`s, not because they all have builder functions.
