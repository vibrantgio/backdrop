# backdrop

A solid-colour background layer for
[Vibrant Gio](https://github.com/vibrantgio), a design system for native desktop
applications on macOS, Windows and Linux, written in pure Go on
[Gio](https://gioui.org). Two functions and no state.

Filling a Gio window with a colour is three ops in a specific order: a clip to
push, a `paint.ColorOp` to set the colour, and a `paint.PaintOp` to actually
put pixels down. Get the order wrong and nothing appears; leave the clip out
and the fill escapes the component it was meant to stay in.

`Widget` is that sequence as a `layout.Widget`: it clips to the incoming
constraints, fills them, and returns them. In an MVU application it is the
bottom layer of the stack, and because the layers are observables it re-emits
whenever the theme does — which is how a window follows the OS between light
and dark with no imperative code at all. `Fill` is the same paint without the
clip, for a caller that has already pushed one.

## Where it sits

Tier 0 of the stack — `mvu → theme → components → effects → patterns → markdown` —
a leaf whose only dependency is Gio. The
[organization page](https://github.com/vibrantgio) has the full tier table.

It sits alongside [gradient](https://github.com/vibrantgio/gradient) and
[circle](https://github.com/vibrantgio/circle), the other two drawing leaves.

```sh
go get github.com/vibrantgio/backdrop
```

The Gio and Go versions are the ones `go.mod` names.

## Packages

One package, at the module root.

| Symbol | |
| --- | --- |
| `Widget(fill)` | A `layout.Widget` that clips to `gtx.Constraints.Max`, fills it, and returns it. |
| `Fill(ops, fill)` | The colour and paint ops straight into an `*op.Ops`, with no clip of its own — it fills whatever clip is current. |

Both take `color.Color`, not `color.NRGBA`, and convert through
`color.NRGBAModel`, so a `colornames` value goes in directly and so does a
`tokens.PlatformColors` field.

## Usage

The static case is one line in a window:

```go
window := mvu.NewWindow(app.Title("MVU - Backdrop"))
backdrop := rx.Of(backdrop.Widget(colornames.Grey600))
window.Render(backdrop).Wait()
```

`rx.Of` is what turns a component into a layer: an MVU window renders
`rx.Observable[layout.Widget]` streams stacked back to front, and a background
that never changes is a stream of exactly one element.

The case worth copying is the themed one. This is the bottom layer of a
window that follows the OS colour scheme, in full:

```go
func BackdropLayer(th rx.Observable[theme.Theme]) rx.Observable[layout.Widget] {
	colors := rx.SwitchMap(th, func(t theme.Theme) rx.Observable[tokens.PlatformColors] {
		return t.Platform
	})
	return rx.Map(colors, func(c tokens.PlatformColors) layout.Widget {
		return backdrop.Widget(c.WindowBackground)
	})
}
```

There is no invalidation call and no light/dark branch. `th` comes from
[theme](https://github.com/vibrantgio/theme), whose `Platform` observable
emits the platform's colour set, and emits it again when the appearance changes;
when it emits, `Map` builds a new `backdrop.Widget` on the window background
colour, and the window renders it.

Inside a shape you already clipped, drop to `Fill`:

```go
defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
backdrop.Fill(gtx.Ops, backColor)
```

## For coding assistants

Read the org guide before you write code against this module: the plan
root's [`AGENTS.md`](https://github.com/vibrantgio/organisation/blob/master/AGENTS.md).
The build and test commands are reached from that file.

## Status

Honest about what does not work yet.

- **`Widget` always takes all the space it is offered.** It clips to and
  returns `gtx.Constraints.Max`, ignoring `Constraints.Min`, so as a flex child
  it fills the flex rather than sizing to anything, and it can never be smaller
  than what it was handed. That is correct for the bottom layer of a window and
  wrong everywhere else — a pane background inside a layout wants a clip you
  pushed and a `Fill`, not this.
- **`Fill` paints the current clip, whatever that is.** Called with no clip
  pushed it fills the entire window, silently, over everything already drawn.
- **The colour is an argument, not a colour role.** Neither function knows
  about `tokens.PlatformColors`, so a themed caller writes the
  `SwitchMap`/`Map` bridge to get `c.WindowBackground` in here itself. There is no
  `backdrop.FromTokens`.
- **Solid colour only.** No gradient, no image, no blur, no scrim. A gradient
  background is [gradient](https://github.com/vibrantgio/gradient), a different
  module with a different call shape. `effects/blur` in
  [effects](https://github.com/vibrantgio/effects) has a backdrop pipeline of
  its own (`blur.Backdrop`) that does not use this module.
- **There are no tests and no golden images.** `go test ./...` reports "no test
  files".

## License

MIT — see [LICENSE](./LICENSE).
