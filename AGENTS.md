# AGENTS.md — backdrop

A solid-colour background, in two functions: `Fill` paints the current clip
into an op list, and `Widget` returns a `layout.Widget` that clips to the
incoming constraints and fills them.

**Layer.** Tier 0 of ADR-001's table — a leaf whose only dependency is Gio.
`mvu/example` and four of the workbench applications use it; nothing in the
design system does.

**Read the canonical guide before you write code against this module.** It is
the organization's one agent guide — the module inventory with current tags,
the application skeleton, the MVU loop and rx semantics, typography, and the
pitfalls that are not guessable. It lives exactly once, in `vibrantgio/.github`,
and this file links it rather than copying it:

    https://raw.githubusercontent.com/vibrantgio/.github/master/llms.txt

**Module.** `github.com/vibrantgio/backdrop`, one module at the repository
root.

**Build and test.** From the repository root:

    go build ./... && go test ./...
