package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var uiFace text.Face

func init() {
	f, err := os.Open("assets/fonts/DotGothic16-Regular.ttf")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	src, err := text.NewGoTextFaceSource(f)
	if err != nil {
		log.Fatal(err)
	}
	uiFace = &text.GoTextFace{Source: src, Size: 22}
}

type DialogueBox struct {
	Rect image.Rectangle
}

func (d DialogueBox) draw(screen *ebiten.Image, txt string) {

	box := ebiten.NewImage(d.Rect.Dx(), d.Rect.Dy())
	box.Fill(color.RGBA{0, 0, 0, 180})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(d.Rect.Min.X), float64(d.Rect.Min.Y))
	screen.DrawImage(box, op)

	tOp := &text.DrawOptions{}
	tOp.GeoM.Translate(float64(d.Rect.Min.X+20), float64(d.Rect.Min.Y+20))
	tOp.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, txt, uiFace, tOp)
}

type Game struct {
	pages       []*Page
	index       int
	stage       *StageRenderer
	dialogueBox DialogueBox
}

func NewGame(pages []*Page) *Game {
	const w, h = 640, 480
	return &Game{
		pages:       pages,
		stage:       NewStageRenderer(w, h),
		dialogueBox: DialogueBox{Rect: image.Rect(0, h*2/3, w, h)},
	}
}

func (g *Game) Update() error {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
		inpututil.IsKeyJustPressed(ebiten.KeySpace) {

		if g.index < len(g.pages)-1 {
			g.index++
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.stage.draw(screen, g.pages[g.index].Stage)

	if g.pages[g.index].Dialogue != nil {
		g.dialogueBox.draw(screen, g.pages[g.index].Clean)
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return g.stage.screenW, g.stage.screenH
}

func main() {
	pages, err := LoadScripts("assets/scripts/demo.json")
	if err != nil {
		log.Fatal(err)
	}

	g := NewGame(pages)
	ebiten.SetWindowSize(1920, 1080)
	ebiten.SetWindowTitle("Novel Game Demo")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
