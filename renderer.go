package main

import (
	"image/color"
	"log"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// StageRenderer は背景と立ち絵を描く & 画像をキャッシュする
type StageRenderer struct {
	bgCache     map[string]*ebiten.Image // assets/bg/ 以下
	spriteCache map[string]*ebiten.Image // assets/sprites/ 以下
	lastBG      string                   // BG が省略されたページ用
	screenW, screenH int
}

// ctor
func NewStageRenderer(w, h int) *StageRenderer {
	return &StageRenderer{
		bgCache:     map[string]*ebiten.Image{},
		spriteCache: map[string]*ebiten.Image{},
		screenW:     w,
		screenH:     h,
	}
}

// 内部: 画像をロードしてキャッシュ
func (r *StageRenderer) load(cache map[string]*ebiten.Image, dir, file string) *ebiten.Image {
	if img, ok := cache[file]; ok {
		return img
	}
	full := filepath.Join("assets", dir, file)
	img, _, err := ebitenutil.NewImageFromFile(full)
	if err != nil {
		// 読めなければ 1×1 マゼンタで代用
		log.Printf("image load error: %v", err)
		img = ebiten.NewImage(1, 1)
		img.Fill(color.RGBA{255, 0, 255, 255})
	}
	cache[file] = img
	return img
}

// 描画本体
func (r *StageRenderer) draw(dst *ebiten.Image, st *StageInfo) {
	// -------- ① 背景 --------
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

	// -------- ② 立ち絵 --------
	if st == nil {
		return
	}
	for _, s := range st.Sprites {
		sp := r.load(r.spriteCache, "sprites", s.File)

		// X 座標を pos で決定
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
		// 下端をそろえて描く
		y := float64(r.screenH) - float64(sp.Bounds().Dy())

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, y)
		dst.DrawImage(sp, op)
	}
}