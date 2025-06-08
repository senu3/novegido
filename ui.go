package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// NineSlice represents an image that can be drawn using the nine-slice
// technique. The Corner field specifies the size of each corner region
// in pixels.
type NineSlice struct {
	Image  *ebiten.Image
	Corner int
}

// LoadNineSlice loads an image from the given path and returns a NineSlice.
// If the image fails to load, a placeholder is returned and the error is
// reported so the caller can handle it gracefully.
func LoadNineSlice(path string, corner int) (*NineSlice, error) {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Printf("nine-slice load error: %v", err)
		// create a visible placeholder so the game can continue running
		ph := ebiten.NewImage(corner*2, corner*2)
		ph.Fill(color.RGBA{255, 0, 255, 255})
		return &NineSlice{Image: ph, Corner: corner}, err
	}
	return &NineSlice{Image: img, Corner: corner}, nil
}

// Draw renders the nine-slice image into the destination rectangle on dst.
// The source image is split into nine regions using the Corner size. The
// edges and center are scaled to fill the specified rectangle.
func (ns *NineSlice) Draw(dst *ebiten.Image, rect image.Rectangle) {
	if ns == nil || ns.Image == nil {
		return
	}
	cw := ns.Corner
	sw, sh := ns.Image.Size()
	if rect.Dx() < 2*cw || rect.Dy() < 2*cw {
		// rectangle too small, just draw the whole image scaled
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(rect.Dx())/float64(sw), float64(rect.Dy())/float64(sh))
		op.GeoM.Translate(float64(rect.Min.X), float64(rect.Min.Y))
		dst.DrawImage(ns.Image, op)
		return
	}

	// helper to draw a subimage with scaling
	drawPart := func(src image.Rectangle, x, y, w, h int) {
		sub := ns.Image.SubImage(src).(*ebiten.Image)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(w)/float64(src.Dx()), float64(h)/float64(src.Dy()))
		op.GeoM.Translate(float64(x), float64(y))
		dst.DrawImage(sub, op)
	}

	left, top := rect.Min.X, rect.Min.Y
	right, bottom := rect.Max.X, rect.Max.Y
	// corners
	drawPart(image.Rect(0, 0, cw, cw), left, top, cw, cw)
	drawPart(image.Rect(sw-cw, 0, sw, cw), right-cw, top, cw, cw)
	drawPart(image.Rect(0, sh-cw, cw, sh), left, bottom-cw, cw, cw)
	drawPart(image.Rect(sw-cw, sh-cw, sw, sh), right-cw, bottom-cw, cw, cw)

	// edges
	w := rect.Dx() - 2*cw
	h := rect.Dy() - 2*cw
	drawPart(image.Rect(cw, 0, sw-cw, cw), left+cw, top, w, cw)           // top
	drawPart(image.Rect(cw, sh-cw, sw-cw, sh), left+cw, bottom-cw, w, cw) // bottom
	drawPart(image.Rect(0, cw, cw, sh-cw), left, top+cw, cw, h)           // left
	drawPart(image.Rect(sw-cw, cw, sw, sh-cw), right-cw, top+cw, cw, h)   // right

	// center
	drawPart(image.Rect(cw, cw, sw-cw, sh-cw), left+cw, top+cw, w, h)
}
