//go:build !headless
// +build !headless

package ui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// DialogueBox represents the main dialogue area.
type DialogueBox struct {
	Rect      image.Rectangle
	Frame     *NineSlice
	NameFrame *NineSlice
}

// Draw renders the dialogue box along with speaker name and text.
func (d DialogueBox) Draw(screen *ebiten.Image, face text.Face, name, txt string) {
	if d.Frame != nil {
		d.Frame.Draw(screen, d.Rect)
	} else {
		box := ebiten.NewImage(d.Rect.Dx(), d.Rect.Dy())
		box.Fill(color.RGBA{0, 0, 0, 180})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(d.Rect.Min.X), float64(d.Rect.Min.Y))
		screen.DrawImage(box, op)
	}

	y := float64(d.Rect.Min.Y + 20)

	if name != "" {
		nameHeight := 24
		nameRect := image.Rect(
			d.Rect.Min.X+20,
			d.Rect.Min.Y+10,
			d.Rect.Min.X+20+d.Rect.Dx()/3,
			d.Rect.Min.Y+10+nameHeight,
		)
		if d.NameFrame != nil {
			d.NameFrame.Draw(screen, nameRect)
		} else {
			nameBox := ebiten.NewImage(nameRect.Dx(), nameRect.Dy())
			nameBox.Fill(color.RGBA{0, 0, 0, 220})
			nOp := &ebiten.DrawImageOptions{}
			nOp.GeoM.Translate(float64(nameRect.Min.X), float64(nameRect.Min.Y))
			screen.DrawImage(nameBox, nOp)
		}

		ntOp := &text.DrawOptions{}
		ntOp.GeoM.Translate(float64(nameRect.Min.X+10), float64(nameRect.Max.Y-6))
		ntOp.ColorScale.ScaleWithColor(color.White)
		text.Draw(screen, name, face, ntOp)

		y += float64(nameRect.Dy() + 10)
	}

	tOp := &text.DrawOptions{}
	tOp.GeoM.Translate(float64(d.Rect.Min.X+20), y)
	tOp.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, txt, face, tOp)
}
