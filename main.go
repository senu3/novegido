//go:build !headless
// +build !headless

package main

import (
	"flag"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"novegido/internal/game"
	"novegido/internal/script"
	"novegido/internal/ui"
)

var (
	screenWidth  = flag.Int("width", 640, "screen width")
	screenHeight = flag.Int("height", 480, "screen height")
)

func main() {
	flag.Parse()

	pages, err := script.LoadScripts("assets/scripts/demo.json")
	if err != nil {
		log.Fatal(err)
	}

	uiObj, err := ui.New()
	if err != nil {
		log.Fatal(err)
	}
	g := game.NewGame(uiObj, pages, *screenWidth, *screenHeight)
	ebiten.SetWindowSize(*screenWidth, *screenHeight)
	ebiten.SetWindowTitle("Novel Game Demo")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
