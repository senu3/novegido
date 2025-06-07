package main

import (
	"image"
	"image/color"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
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

type mp3Source struct {
	*mp3.Stream
	f *os.File
}

func (m *mp3Source) Close() error {
	return m.f.Close()
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
	audioCtx    *audio.Context
	players     map[string]*audio.Player
	sources     map[string]io.Closer
	bgm         *audio.Player
	bgmFile     string
}

func NewGame(pages []*Page) *Game {
	const w, h = 640, 480
	g := &Game{
		pages:       pages,
		stage:       NewStageRenderer(w, h),
		dialogueBox: DialogueBox{Rect: image.Rect(0, h*2/3, w, h)},
		audioCtx:    audio.NewContext(48000),
		players:     map[string]*audio.Player{},
		sources:     map[string]io.Closer{},
	}
	if len(pages) > 0 {
		g.playAudio(pages[0].Audio)
	}
	return g
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
		inpututil.IsKeyJustPressed(ebiten.KeySpace) {

		if g.index < len(g.pages)-1 {
			g.index++
			g.playAudio(g.pages[g.index].Audio)
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

func (g *Game) playAudio(info *AudioInfo) {
	if info == nil || info.File == "" {
		return
	}

	if p, ok := g.players[info.File]; ok {
		_ = p.Rewind()
		if info.Loop {
			if g.bgm != nil && g.bgm != p {
				g.bgm.Pause()
				g.bgmFile = info.File
			}
			g.bgm = p
		}
		p.Play()
		return
	}

	f, err := os.Open(filepath.Join("assets", info.File))
	if err != nil {
		log.Printf("audio load error: %v", err)
		return
	}
	stream, err := mp3.DecodeWithoutResampling(f)
	if err != nil {
		log.Printf("decode error: %v", err)
		_ = f.Close()
		return
	}
	src := &mp3Source{Stream: stream, f: f}
	var reader io.ReadSeeker = src
	if info.Loop {
		reader = audio.NewInfiniteLoop(reader, stream.Length())
	}
	p, err := g.audioCtx.NewPlayer(reader)
	if err != nil {
		log.Printf("player error: %v", err)
		_ = src.Close()
		return
	}
	g.players[info.File] = p
	g.sources[info.File] = src
	if info.Loop {
		if g.bgm != nil && g.bgm != p {
			g.discardPlayer(g.bgmFile)
		}
		g.bgm = p
		g.bgmFile = info.File
	}
	p.Play()
}

func (g *Game) discardPlayer(file string) {
	if file == "" {
		return
	}
	if p, ok := g.players[file]; ok {
		p.Close()
		delete(g.players, file)
	}
	if s, ok := g.sources[file]; ok {
		s.Close()
		delete(g.sources, file)
	}
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
