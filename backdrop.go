package backdrop

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

func Fill(ops *op.Ops, fill color.Color) {
	paint.ColorOp{Color: color.NRGBAModel.Convert(fill).(color.NRGBA)}.Add(ops)
	paint.PaintOp{}.Add(ops)
}

func Widget(fill color.Color) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		defer clip.Outline{Path: clip.Rect{Max: gtx.Constraints.Max}.Path()}.Op().Push(gtx.Ops).Pop()
		Fill(gtx.Ops, fill)
		return layout.Dimensions{Size: gtx.Constraints.Max}
	}
}
