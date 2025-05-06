package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"

	engine "novegido/internal"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Game ----------------------------------------------------------------------

var (
	helloFace text.Face
)

type Game struct {
	player   Player
	fontface text.Face
}

type Player struct {
	Img  *ebiten.Image
	X, Y float64
}

func init() {

	f, err := os.Open("asset/fonts/DotGothic16-Regular.ttf")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	src, err := text.NewGoTextFaceSource(f)
	if err != nil {
		return
	}
	helloFace = &text.GoTextFace{Source: src, Size: 24}
}

func (g *Game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.X += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.player.Y += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.Y -= 2
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{120, 180, 255, 255})

	opts := ebiten.DrawImageOptions{}

	opts.GeoM.Translate(g.player.X, g.player.Y)

	screen.DrawImage(
		g.player.Img.SubImage(
			image.Rect(0, 0, 35, 35),
		).(*ebiten.Image),
		&opts,
	)

	opts.GeoM.Reset()

	tOp := &text.DrawOptions{}
	tOp.GeoM.Translate(10, 30) // 画面左上に配置
	tOp.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, "Hello World", g.fontface, tOp)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}

func Load() {
	// デモ用の読み込み
	pages, err := engine.LoadPages("asset/scripts/demo.json")
	if err != nil {
		panic(err)
	}

	// 台本を順に処理するデモ
	for i, page := range pages {
		fmt.Printf("ページ %d:\n", i+1)

		// 演出情報の表示
		if page.Stage != nil {
			fmt.Printf("- 背景: %s\n", page.Stage.BG)
			for _, s := range page.Stage.Sprites {
				fmt.Printf("- 立ち絵 [%s]: %s (%s)\n", s.ID, s.File, s.Pos)
			}
		}

		// セリフとTIPS用語を表示
		if page.Dialogue != nil {
			cleanText := engine.ParseDialogue(page.Dialogue.Text)
			fmt.Printf("%s「%s」\n", page.Dialogue.Speaker, cleanText)
		}
	}
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	PlayerImg, _, err := ebitenutil.NewImageFromFile("asset/sprite/tc_cat_calico[A].png")
	if err != nil {
		log.Fatal(err)
	}

	game := Game{
		player: Player{
			Img: PlayerImg,
			X:   100,
			Y:   100,
		},
		fontface: helloFace,
	}
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
