package main

import (
	"image/color"
	"log"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type StageRenderer struct {
	bgCache          map[string]*ebiten.Image
	spriteCache      map[string]*ebiten.Image
	lastBG           string
	screenW, screenH int
}

func NewStageRenderer(w, h int) *StageRenderer {
	return &StageRenderer{
		bgCache:     map[string]*ebiten.Image{},
		spriteCache: map[string]*ebiten.Image{},
		screenW:     w,
		screenH:     h,
	}
}

func (r *StageRenderer) load(cache map[string]*ebiten.Image, dir, file string) *ebiten.Image {
	if img, ok := cache[file]; ok {
		return img
	}
	full := filepath.Join("assets", dir, file)
	img, _, err := ebitenutil.NewImageFromFile(full)
	if err != nil {

		log.Printf("image load error: %v", err)
		img = ebiten.NewImage(1, 1)
		img.Fill(color.RGBA{255, 0, 255, 255})
	}
	cache[file] = img
	return img
}

func (r *StageRenderer) draw(dst *ebiten.Image, st *StageInfo) {

	if st != nil && st.BG != "" {
		r.lastBG = st.BG
	}
	if r.lastBG != "" {
		bg := r.load(r.bgCache, "bg", r.lastBG)
		op := &ebiten.DrawImageOptions{}
		bw, bh := bg.Size()
		op.GeoM.Scale(float64(r.screenW)/float64(bw), float64(r.screenH)/float64(bh))
		dst.DrawImage(bg, op)
	} else {
		dst.Fill(color.Black)
	}

	if st == nil {
		return
	}
	for _, s := range st.Sprites {
		sp := r.load(r.spriteCache, "sprites", s.File)

		var x float64
		switch s.Pos {
		case "left":
			x = float64(r.screenW)*0.2 - float64(sp.Bounds().Dx())/2
		case "center":
			x = float64(r.screenW)*0.5 - float64(sp.Bounds().Dx())/2
		case "right":
			x = float64(r.screenW)*0.8 - float64(sp.Bounds().Dx())/2
		default:
			x = 0
		}

		y := float64(r.screenH) - float64(sp.Bounds().Dy())

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, y)
		dst.DrawImage(sp, op)
	}
}
