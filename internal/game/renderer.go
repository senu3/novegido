//go:build !headless
// +build !headless

package game

import (
	"image/color"
	"log"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"novegido/internal/script"
)

// StageRenderer handles rendering of backgrounds and sprites with simple fades.
type StageRenderer struct {
	bgCache     map[string]*ebiten.Image
	spriteCache map[string]*ebiten.Image

	currBG        string
	prevBG        string
	bgFadeFrames  int
	bgFadeCounter int

	currSprites       []script.SpriteInfo
	prevSprites       []script.SpriteInfo
	spriteFadeFrames  int
	spriteFadeCounter int

	black *ebiten.Image

	screenW, screenH int
}

// NewStageRenderer creates a renderer for a stage of the given screen size.
func NewStageRenderer(w, h int) *StageRenderer {
	black := ebiten.NewImage(w, h)
	black.Fill(color.Black)
	return &StageRenderer{
		bgCache:     map[string]*ebiten.Image{},
		spriteCache: map[string]*ebiten.Image{},
		black:       black,
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

func (r *StageRenderer) draw(dst *ebiten.Image, st *script.StageInfo) {
	if st != nil {
		if st.BG != "" && st.BG != r.currBG {
			if st.BGFade > 0 {
				r.prevBG = r.currBG
				r.currBG = st.BG
				r.bgFadeFrames = st.BGFade
				r.bgFadeCounter = 0
			} else {
				r.prevBG = ""
				r.currBG = st.BG
				r.bgFadeFrames = 0
			}
		}

		if !spritesEqual(st.Sprites, r.currSprites) {
			if st.SpriteFade > 0 {
				r.prevSprites = r.currSprites
				r.currSprites = append([]script.SpriteInfo(nil), st.Sprites...)
				r.spriteFadeFrames = st.SpriteFade
				r.spriteFadeCounter = 0
			} else {
				r.prevSprites = nil
				r.currSprites = append([]script.SpriteInfo(nil), st.Sprites...)
				r.spriteFadeFrames = 0
			}
		}
	}

	r.drawBackground(dst)
	r.drawSprites(dst)
}

func (r *StageRenderer) drawBackground(dst *ebiten.Image) {
	if r.bgFadeFrames == 0 {
		if r.currBG != "" {
			bg := r.load(r.bgCache, "bg", r.currBG)
			op := &ebiten.DrawImageOptions{}
			bw, bh := bg.Size()
			op.GeoM.Scale(float64(r.screenW)/float64(bw), float64(r.screenH)/float64(bh))
			dst.DrawImage(bg, op)
		} else {
			dst.DrawImage(r.black, nil)
		}
		return
	}

	ratio := float64(r.bgFadeCounter) / float64(r.bgFadeFrames)

	if r.prevBG != "" {
		bg := r.load(r.bgCache, "bg", r.prevBG)
		op := &ebiten.DrawImageOptions{}
		bw, bh := bg.Size()
		op.GeoM.Scale(float64(r.screenW)/float64(bw), float64(r.screenH)/float64(bh))
		op.ColorScale.ScaleAlpha(float32(1 - ratio))
		dst.DrawImage(bg, op)
	} else {
		op := &ebiten.DrawImageOptions{}
		op.ColorScale.ScaleAlpha(float32(1 - ratio))
		dst.DrawImage(r.black, op)
	}

	if r.currBG != "" {
		bg := r.load(r.bgCache, "bg", r.currBG)
		op := &ebiten.DrawImageOptions{}
		bw, bh := bg.Size()
		op.GeoM.Scale(float64(r.screenW)/float64(bw), float64(r.screenH)/float64(bh))
		op.ColorScale.ScaleAlpha(float32(ratio))
		dst.DrawImage(bg, op)
	} else {
		op := &ebiten.DrawImageOptions{}
		op.ColorScale.ScaleAlpha(float32(ratio))
		dst.DrawImage(r.black, op)
	}

	if r.bgFadeCounter < r.bgFadeFrames {
		r.bgFadeCounter++
	}
}

func (r *StageRenderer) drawSprites(dst *ebiten.Image) {
	if r.spriteFadeFrames == 0 {
		r.drawSpriteSet(dst, r.currSprites, 1)
		return
	}

	ratio := float64(r.spriteFadeCounter) / float64(r.spriteFadeFrames)
	r.drawSpriteSet(dst, r.prevSprites, 1-ratio)
	r.drawSpriteSet(dst, r.currSprites, ratio)

	if r.spriteFadeCounter < r.spriteFadeFrames {
		r.spriteFadeCounter++
	}
}

func (r *StageRenderer) drawSpriteSet(dst *ebiten.Image, sprites []script.SpriteInfo, alpha float64) {
	if alpha <= 0 {
		return
	}
	for _, s := range sprites {
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
		if alpha < 1 {
			op.ColorScale.ScaleAlpha(float32(alpha))
		}
		dst.DrawImage(sp, op)
	}
}

func spritesEqual(a, b []script.SpriteInfo) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].File != b[i].File || a[i].Pos != b[i].Pos {
			return false
		}
	}
	return true
}
