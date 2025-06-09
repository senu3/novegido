//go:build !headless
// +build !headless

package ui

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// UI holds assets used for rendering user interface elements.
type UI struct {
	Face text.Face
}

// New loads the default font and returns a UI object.
func New() (*UI, error) {
	f, err := os.Open("assets/fonts/DotGothic16-Regular.ttf")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	src, err := text.NewGoTextFaceSource(f)
	if err != nil {
		return nil, err
	}
	return &UI{Face: &text.GoTextFace{Source: src, Size: 22}}, nil
}
