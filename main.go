package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// ---------------- フォントを 1 回だけ生成 -----------------------------

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

// ------------------ 画面パーツごとの構造体 ----------------------------

// DialogueBox は画面下部にセリフを描画するだけの単純なパーツ
type DialogueBox struct {
	Rect image.Rectangle
}

func (d DialogueBox) draw(screen *ebiten.Image, txt string) {
	// 背景（半透明黒）
	box := ebiten.NewImage(d.Rect.Dx(), d.Rect.Dy())
	box.Fill(color.RGBA{0, 0, 0, 180})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(d.Rect.Min.X), float64(d.Rect.Min.Y))
	screen.DrawImage(box, op)

	// テキスト
	tOp := &text.DrawOptions{}
	tOp.GeoM.Translate(float64(d.Rect.Min.X+20), float64(d.Rect.Min.Y+20))
	tOp.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, txt, uiFace, tOp)
}

// StageRenderer は背景・立ち絵を描くパーツ（簡易版 ― 実画像は割愛）
type StageRenderer struct{}

func (StageRenderer) draw(screen *ebiten.Image, st *StageInfo) {
	if st == nil {
		return
	}
	// 本サンプルでは「色付き矩形」を代用品として表示
	bg := ebiten.NewImage(screen.Size())
	bg.Fill(color.RGBA{120, 180, 255, 255})
	screen.DrawImage(bg, nil)

	// 立ち絵も割愛（実装例: ファイル読み込み→キャッシュ→DrawImage）
}

// ------------------ Game 本体 -----------------------------------------

type Game struct {
	pages       []*Page
	index       int
	stage       StageRenderer
	dialogueBox DialogueBox
}

// コンストラクタ的な
func NewGame(pages []*Page) *Game {
	// 下 1/3 をテキストウィンドウにする
	w, h := 640, 480
	return &Game{
		pages: pages,
		dialogueBox: DialogueBox{
			Rect: image.Rect(0, h*2/3, w, h),
		},
	}
}

func (g *Game) Update() error {
	// マウス左クリック or スペースキー で次ページへ
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
		inpututil.IsKeyJustPressed(ebiten.KeySpace) {

		if g.index < len(g.pages)-1 {
			g.index++
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// ① 演出
	g.stage.draw(screen, g.pages[g.index].Stage)

	// ② セリフ
	if g.pages[g.index].Dialogue != nil {
		g.dialogueBox.draw(screen, g.pages[g.index].Clean)
	}
}

func (g *Game) Layout(w, h int) (int, int) { return 640, 480 }

// ------------------ main ----------------------------------------------

func main() {
	pages, err := LoadScripts("assets/scripts/demo.json")
	if err != nil {
		log.Fatal(err)
	}

	g := NewGame(pages)
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Novel Game Demo")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
