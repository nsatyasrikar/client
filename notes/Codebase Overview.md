---
title: Codebase Overview
tags:
  - codebase
  - index
aliases:
  - client
  - Overview
---

# client — Codebase Overview

Index note for the `client` Go module. Start here.

> [!abstract] What this is
> A typed Go AST for building JavaScript/TypeScript, compiled via
> [esbuild](https://github.com/evanw/esbuild). Types are enforced at the Go
> compiler level through `Value[T Type]` generics — see [[Architecture#Type system]].
> User-facing docs and runnable examples live in the repo [README](../README.md).

## Map of notes

- [[Architecture]] — AST/compiler design, sealed-interface pattern, package boundary
- [[Browser Package]] — the `browser/` DOM bindings package
- [[Decisions Log]] — non-obvious choices and why, for future-you

## Package layout at a glance

| Package | Path | Role |
|---|---|---|
| `client` | repo root | AST, type system, validator, TypeScript emitter, esbuild compiler |
| `browser` | `browser/` | DOM/browser API bindings — hand-written wrappers + 1034 static marker types |

## Quick facts

- Module: `github.com/nsatyasrikar/client`
- No API key or network access needed to build/test
- `module.go` and `class.go` are currently empty stubs — module/class emission isn't implemented
- Knowledge graph lives at `graphify-out/` (gitignored, regenerates via `graphify update .` — see [[Decisions Log#graphify]])

## Related code

- [README.md](../README.md) — the canonical, example-driven documentation
- [go.mod](../go.mod)
