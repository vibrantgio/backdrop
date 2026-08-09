# AGENTS.md — backdrop

A solid-colour background, in two functions: `Fill` paints the current clip
into an op list, and `Widget` returns a `layout.Widget` that clips to the
incoming constraints and fills them.

**Layer.** Tier 0 of ADR-001's table — a leaf whose only dependency outside
the organization is Gio. Its root module imports nothing else in the
organization. No other repository's root module imports it; outside the
tier table it is imported by the demo module `mvu/example` and the
workbench applications `iconbrowser`, `launcher`, `mindchat` and `todos`.
Both directions are measured rather than typed — `scripts/check-layers.sh
--edges` reports the graph and `scripts/sync-agents.sh` renders these
sentences from it — so correcting them here changes nothing.

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
