---
title: Browser Package
tags:
  - codebase
  - browser
---

# Browser Package

Part of [[Codebase Overview]]. Lives at `browser/` (package `browser`).

## Contents

- `browser.go` — hand-written DOM API wrappers: `Document.GetElementByID`,
  `ElementAPI.SetTextContent`, `ElementAPI.AddEventListener`. Built entirely
  from `client`'s exported API (`Property`, `CallValue`, `Assign`,
  `ExpressionStatement`, `UnsafeExpression`) — no access to `client`'s
  unexported AST internals.
- `types_generated.go` — **1034** static marker types (`ActionType`,
  `AlarmType`, ... one per browser/DOM interface), each just:
  ```go
  type ActionType struct{ client.TypeMarker }
  ```

## TypeMarker

`client.Type` is sealed by an unexported method (see [[Architecture#AST shape]]),
which normally means only `package client` can implement it. `TypeMarker` is
the escape hatch:

```go
// in package client
type TypeMarker struct{}
func (TypeMarker) clientType() {}
```

Any type outside `package client` can embed `TypeMarker` and inherit the
promoted `clientType()` method — satisfying `Type` without needing the
unexported method itself. This is what makes the 1034 generated types in
this package legal.

> [!info] Only a few of the 1034 types have builder methods
> The generated types exist so the Go type system *recognizes* those DOM
> interfaces as valid `client.Type`s (usable in `Value[T]`, `Property[T,R]`,
> etc.) — not because each one has a hand-written function like
> `GetElementByID`. Adding real DOM interaction for a new type means writing
> the wrapper function by hand in `browser.go`.

## History

This package didn't always exist — see [[Decisions Log]] for how it got
here (originally 1034 separate files at repo root, moved and consolidated
over several passes) and why the generator that used to produce
`types_generated.go` was removed.
