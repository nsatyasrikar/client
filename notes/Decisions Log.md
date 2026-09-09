---
title: Decisions Log
tags:
  - codebase
  - decisions
---

# Decisions Log

Part of [[Codebase Overview]]. Non-obvious choices, recorded so nobody
re-derives them from scratch later.

## Why browser split but not core

Repo root used to hold 1049 flat `.go` files: 15 hand-written AST/compiler
files plus 1034 generated DOM binding stubs. Only the generated-bindings
layer moved to `browser/`.

Splitting the 15 core files (`ast`/`compiler`/`emit`/`validate`) into real
Go packages was considered and **rejected** — `expression`/`Statement`/`Type`
are sealed by unexported methods and `emit.go`/`validate.go` type-switch on
unexported concrete structs (see [[Architecture#AST shape]]). A real package
boundary would require exporting ~20 internal node types: a public-API
redesign, not a navigation cleanup. Scope was narrowed to the generated
layer only.

## The `TypeMarker` sealing fix

To let generated types live in `browser/` while `client.Type` stays sealed,
`client.TypeMarker` was added as an embeddable exported struct with the
unexported `clientType()` method — the standard Go pattern for a
cross-package-embeddable sealed interface. See [[Browser Package#TypeMarker]].

## The `"undefined"` bug

While rewriting `browser.go`'s DOM wrappers to use only `client`'s exported
API (previously they reached into unexported AST fields directly, only
possible because they lived inside `package client`), a test surfaced a
pre-existing, previously undetected bug: `Call`, `CallValue`, `Join`,
`ArrayOf`, `Object`, and `ExpressionStatement` all stored a `Value[T]`
argument directly instead of unwrapping its AST node, so **any** argument
passed through them rendered as the literal string `"undefined"` in emitted
output.

It had never been caught because there were no tests in the repo before
this work, and the only real caller (`browser.go`'s original hand-built
`AddEventListener`) happened to reach into unexported fields directly,
bypassing the broken path.

Fixed with one shared helper (`Value[T].unwrapNode()` + `unwrapExpr()` in
`type.go`), applied at all six call sites, landed as its own commit
separate from the package move (never mix a structural change with a
behavioral bug fix in one commit).

## Why `types_generated.go` is one file, not 1034

The 1034 generated stub types are pure markers with zero methods, never
hand-edited. Real Go code generators (protoc-gen-go, stringer, mockgen) emit
one aggregated file per generation run, not one file per declaration — the
one-file-per-type layout here was the unusual choice, not the consolidated
one. Collapsed to `browser/types_generated.go`.

## Why the code generator was removed

`cmd/clientgen` + `generator.go` (`GenerateBindings`) read an `index.json`
schema file that **never existed in this repo** — confirmed via git history
and a full filesystem search. The generator had been dead, unrunnable code
since the repo's first commit. Removed; `types_generated.go`'s header no
longer claims a live generator owns it — it's now explicitly static and
hand-maintained.

## graphify

A knowledge graph of this codebase lives at `graphify-out/` (gitignored —
[[Codebase Overview]]), rebuilt automatically by a post-commit git hook
(`graphify hook install`). AST-only extraction, no LLM/API cost for a
code-only corpus like this one. Query it with `graphify query "<question>"`
before grepping raw source for architecture questions.
