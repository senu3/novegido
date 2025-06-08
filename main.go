//go:build !headless
// +build !headless

package main

import (
	"flag"
	"fmt"
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

var (
	screenWidth  = flag.Int("width", 640, "screen width")
	screenHeight = flag.Int("height", 480, "screen height")
)

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
	Rect      image.Rectangle
	Frame     *NineSlice
	NameFrame *NineSlice
}

type DialogueEntry struct {
	Speaker string
	Text    string
}

type mp3Source struct {
	*mp3.Stream
	f *os.File
}

func (m *mp3Source) Close() error {
	return m.f.Close()
}

func (d DialogueBox) draw(screen *ebiten.Image, name, txt string) {
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
		text.Draw(screen, name, uiFace, ntOp)

		y += float64(nameRect.Dy() + 10)
	}

	tOp := &text.DrawOptions{}
	tOp.GeoM.Translate(float64(d.Rect.Min.X+20), y)
	tOp.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, txt, uiFace, tOp)
}

type Game struct {
	pages         []*Page
	index         int
	stage         *StageRenderer
	dialogueBox   DialogueBox
	audioCtx      *audio.Context
	players       map[string]*audio.Player
	sources       map[string]io.Closer
	bgm           *audio.Player
	bgmFile       string
	width         int
	height        int
	backlog       []DialogueEntry
	showBacklog   bool
	backlogOffset int
	choosing      bool
	choiceIndex   int
}

func (g *Game) nextPage() {
	if g.index >= len(g.pages)-1 {
		return
	}
	g.index++
	g.playAudio(g.pages[g.index].Audio)
	if g.pages[g.index].Dialogue != nil {
		g.backlog = append(g.backlog, DialogueEntry{
			Speaker: g.pages[g.index].Dialogue.Speaker,
			Text:    g.pages[g.index].Clean,
		})
	}
}

func (g *Game) prevPage() {
	if g.index <= 0 {
		return
	}
	g.index--
	g.playAudio(g.pages[g.index].Audio)
}

func NewGame(pages []*Page, w, h int) *Game {
	frame, err := LoadNineSlice(filepath.Join("assets", "ui", "9slice30.png"), 30)
	if err != nil {
		log.Printf("nine-slice load error: %v", err)
	}
	g := &Game{
		pages: pages,
		stage: NewStageRenderer(w, h),
		dialogueBox: DialogueBox{
			Rect:      image.Rect(0, h*2/3, w, h),
			Frame:     frame,
			NameFrame: frame,
		},
		audioCtx:    audio.NewContext(48000),
		players:     map[string]*audio.Player{},
		sources:     map[string]io.Closer{},
		width:       w,
		height:      h,
		choiceIndex: 0,
	}
	if len(pages) > 0 {
		g.playAudio(pages[0].Audio)
		if pages[0].Dialogue != nil {
			g.backlog = append(g.backlog, DialogueEntry{
				Speaker: pages[0].Dialogue.Speaker,
				Text:    pages[0].Clean,
			})
		}
	}
	return g
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		g.showBacklog = !g.showBacklog
		return nil
	}

	if g.showBacklog {
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			if g.backlogOffset < len(g.backlog)-1 {
				g.backlogOffset++
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			if g.backlogOffset > 0 {
				g.backlogOffset--
			}
		}
		return nil
	}

	if g.choosing {
		choices := g.pages[g.index].Choices
		if len(choices) == 0 {
			g.choosing = false
			return nil
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) || inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
			g.prevPage()
			g.choosing = false
			return nil
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			if g.choiceIndex > 0 {
				g.choiceIndex--
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			if g.choiceIndex < len(choices)-1 {
				g.choiceIndex++
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			dest := choices[g.choiceIndex].Page
			if dest >= 0 && dest < len(g.pages) {
				g.index = dest
				g.playAudio(g.pages[g.index].Audio)
				if g.pages[g.index].Dialogue != nil {
					g.backlog = append(g.backlog, DialogueEntry{
						Speaker: g.pages[g.index].Dialogue.Speaker,
						Text:    g.pages[g.index].Clean,
					})
				}
			}
			g.choosing = false
		}
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) || inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		g.prevPage()
	}

	trigger := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
		inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		inpututil.IsKeyJustPressed(ebiten.KeyEnter)

	if len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		trigger = true
	}

	for _, id := range ebiten.AppendGamepadIDs(nil) {
		if inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonRightBottom) {
			trigger = true
			break
		}
	}

	if trigger {

		if len(g.pages[g.index].Choices) > 0 {
			g.choosing = true
			g.choiceIndex = 0
			return nil
		}

		g.nextPage()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.stage.draw(screen, g.pages[g.index].Stage)

	if g.showBacklog {
		g.drawBacklog(screen)
		return
	}

	if g.pages[g.index].Dialogue != nil {
		dlg := g.pages[g.index].Dialogue
		g.dialogueBox.draw(screen, dlg.Speaker, g.pages[g.index].Clean)
	}

	if g.choosing {
		g.drawChoices(screen)
	}
}

func (g *Game) drawBacklog(screen *ebiten.Image) {
	box := ebiten.NewImage(g.width, g.height)
	box.Fill(color.RGBA{0, 0, 0, 220})
	screen.DrawImage(box, nil)

	lines := g.height / 24
	start := len(g.backlog) - 1 - g.backlogOffset
	y := float64(20)
	for i := 0; i < lines && start-i >= 0; i++ {
		e := g.backlog[start-i]
		tOp := &text.DrawOptions{}
		tOp.GeoM.Translate(20, y)
		tOp.ColorScale.ScaleWithColor(color.White)
		text.Draw(screen, fmt.Sprintf("%s: %s", e.Speaker, e.Text), uiFace, tOp)
		y += 24
	}
}

func (g *Game) drawChoices(screen *ebiten.Image) {
	choices := g.pages[g.index].Choices
	if len(choices) == 0 {
		return
	}
	startY := float64(g.dialogueBox.Rect.Min.Y + 60)
	for i, c := range choices {
		tOp := &text.DrawOptions{}
		tOp.GeoM.Translate(float64(g.dialogueBox.Rect.Min.X+40), startY+float64(i*24))
		col := color.RGBA{255, 255, 255, 255}
		if i == g.choiceIndex {
			col = color.RGBA{255, 255, 0, 255}
		}
		tOp.ColorScale.ScaleWithColor(col)
		text.Draw(screen, fmt.Sprintf("%d. %s", i+1, c.Text), uiFace, tOp)
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return g.width, g.height
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
	flag.Parse()

	pages, err := LoadScripts("assets/scripts/demo.json")
	if err != nil {
		log.Fatal(err)
	}

	g := NewGame(pages, *screenWidth, *screenHeight)
	ebiten.SetWindowSize(*screenWidth, *screenHeight)
	ebiten.SetWindowTitle("Novel Game Demo")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
