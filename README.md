# backdrop

A solid-colour background layer for
[VibrantGio](https://github.com/vibrantgio), a design system for native desktop
applications on macOS, Windows and Linux, written in pure Go on
[Gio](https://gioui.org). Two functions, twenty lines, no state — and the most
widely used of the organization's drawing leaves.

Filling a Gio window with a colour is three ops in a specific order: a clip to
push, a `paint.ColorOp` to set the material, and a `paint.PaintOp` to actually
put pixels down. Get the order wrong and nothing appears; leave the clip out
and the fill escapes the widget it was meant to stay in. It is the first thing
every application does and the first thing everyone writes twice.

`Widget` is that sequence as a `layout.Widget`: it clips to the incoming
constraints, fills them, and returns them. In an MVU application it is the
bottom layer of the stack, and because the layers are observables it re-emits
whenever the theme does — which is how a window follows the OS between light
and dark with no imperative code at all. `Fill` is the same paint without the
clip, for a caller that has already pushed one.

## Where it sits

Tier 0 of the stack — `mvu → spectrum → prism → pulse → cadence → markdown` —
a leaf whose only dependency is Gio. The
[organization page](https://github.com/vibrantgio) has the full tier table.

It sits alongside [gradient](https://github.com/vibrantgio/gradient) and
[circle](https://github.com/vibrantgio/circle), the other two drawing leaves,
and it is the one that is genuinely load-bearing: ten `mvu/example` programs
and four of the seven
[workbench](https://github.com/vibrantgio/workbench) applications — `todos`,
`iconbrowser`, `launcher` and `mindchat` — import it, against one consumer for
gradient and five for circle. Nothing inside the design system imports it;
components paint their own surfaces.

```sh
go get github.com/vibrantgio/backdrop
```

Every module in the organization is on gioui.org v0.10.1 and Go 1.25.1.

## Packages

One package, at the module root.

| Symbol | |
| --- | --- |
| `Widget(fill)` | A `layout.Widget` that clips to `gtx.Constraints.Max`, fills it, and returns it. |
| `Fill(ops, fill)` | The colour and paint ops straight into an `*op.Ops`, with no clip of its own — it fills whatever clip is current. |

Both take `color.Color`, not `color.NRGBA`, and convert through
`color.NRGBAModel`, so a `colornames` value goes in directly and so does a
`tokens.ColorTokens` field.

## Usage

The static case is one line. This is the whole of `mvu/example/02-backdrop`:

```go
window := mvu.NewWindow(app.Title("MVU - Backdrop"))
backdrop := rx.Of(backdrop.Widget(colornames.Grey600))
window.Render(backdrop).Wait()
```

`rx.Of` is what turns a widget into a layer: an MVU window renders
`rx.Observable[layout.Widget]` streams stacked back to front, and a background
that never changes is a stream of exactly one element.

The case worth copying is the themed one, because it is where the reactive
plumbing earns its keep. This is `backdrop.go` from
[workbench/todos](https://github.com/vibrantgio/workbench/tree/master/todos) in
full — the bottom layer of an application that follows the OS colour scheme:

```go
func BackdropLayer(th rx.Observable[theme.Theme]) rx.Observable[layout.Widget] {
	colors := rx.SwitchMap(th, func(t theme.Theme) rx.Observable[tokens.ColorTokens] {
		return t.Color
	})
	return rx.Map(colors, func(c tokens.ColorTokens) layout.Widget {
		return backdrop.Widget(c.Background)
	})
}
```

There is no invalidation call and no light/dark branch. `th` comes from
[spectrum](https://github.com/vibrantgio/spectrum), which polls the OS and
emits only on change; when it emits, `SwitchMap` resolves the new colour
tokens, `Map` builds a new `backdrop.Widget`, and the window renders it. Each
of the four workbench applications with a backdrop repeats those same few
lines.

Inside a shape you already clipped, drop to `Fill` — from `mvu/example/edit`,
filling a text field's own rectangle:

```go
defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
backdrop.Fill(gtx.Ops, backColor)
```

## For coding assistants

Read the canonical guide before writing code against this module — the module
inventory with current tags, the application skeleton, MVU and rx semantics,
typography, and the pitfalls that are not guessable:

<https://raw.githubusercontent.com/vibrantgio/.github/master/llms.txt>

[`AGENTS.md`](./AGENTS.md) in this repository has the build and test commands.

## Status

Honest about what does not work yet. Every count below is measured.

- **`Widget` always takes all the space it is offered.** It clips to and
  returns `gtx.Constraints.Max`, ignoring `Constraints.Min`, so as a flex child
  it fills the flex rather than sizing to anything, and it can never be smaller
  than what it was handed. That is correct for the bottom layer of a window and
  wrong everywhere else — a panel background inside a layout wants a clip you
  pushed and a `Fill`, not this.
- **`Fill` paints the current clip, whatever that is.** Called with no clip
  pushed it fills the entire window, silently, over everything already drawn.
  It has one call site in the whole organization — `mvu/example/edit` — against
  fourteen for `Widget`, so the trap is mostly untripped rather than mostly
  avoided.
- **The colour is an argument, not a role.** Neither function knows about
  `tokens.ColorTokens`, so every application repeats the same `SwitchMap`/`Map`
  bridge to get `c.Background` in here; four of them do, near-identically.
  There is no `backdrop.FromTokens`, and no phase of the current plan adds one.
- **Solid colour only.** No gradient, no image, no blur, no scrim. A gradient
  background is [gradient](https://github.com/vibrantgio/gradient), a different
  module with a different call shape and one consumer; a scrim over a modal is
  hand-rolled — `todos` builds its own from `color.NRGBA{A: 153}` rather than
  layering this. Phase E of the
  [org plan](https://github.com/vibrantgio/.github) builds `pulse/blur` on
  `gioui.org/gpu/headless` with a backdrop pipeline; it does not claim this
  module, and the two will need reconciling.
- **There are no tests and no golden images.** `go test ./...` reports "no test
  files". It is twenty lines and it is the first thing every application draws,
  and nothing pins it.

## License

MIT — see [LICENSE](./LICENSE).
