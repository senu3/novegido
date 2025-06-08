package script

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"
)

// SpriteInfo describes a character sprite on screen.
type SpriteInfo struct {
	ID   string `json:"id"`
	File string `json:"file"`
	Pos  string `json:"pos"`
}

// StageInfo describes background and sprite placement along with transitions.
type StageInfo struct {
	BG         string       `json:"bg"`
	Sprites    []SpriteInfo `json:"sprites"`
	BGFade     int          `json:"bgFade,omitempty"`
	SpriteFade int          `json:"spriteFade,omitempty"`
}

// DialogueInfo holds spoken text and speaker name.
type DialogueInfo struct {
	Speaker string `json:"speaker"`
	Text    string `json:"text"`
}

// AudioInfo describes a sound file that should be played.
type AudioInfo struct {
	File string `json:"file"`
	Loop bool   `json:"loop"`
}

// ChoiceInfo represents a selectable option leading to another page.
type ChoiceInfo struct {
	Text string `json:"text"`
	Page int    `json:"page"`
}

// Page is a single entry of a script.
type Page struct {
	Stage    *StageInfo    `json:"stage,omitempty"`
	Dialogue *DialogueInfo `json:"dialogue,omitempty"`
	Audio    *AudioInfo    `json:"audio,omitempty"`
	Choices  []ChoiceInfo  `json:"choices,omitempty"`
	Clean    string        `json:"-"`
}

// LoadScripts reads a JSON script file and returns parsed pages.
func LoadScripts(path string) ([]*Page, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var pages []*Page
	if err := json.NewDecoder(f).Decode(&pages); err != nil {
		return nil, err
	}

	for _, p := range pages {
		if p.Dialogue != nil {
			p.Clean = ParseDialogue(p.Dialogue.Text)
		}
	}
	return pages, nil
}

// ParseDialogue removes any markup such as HTML tags and normalises whitespace.
func ParseDialogue(src string) string {
	out := strings.ReplaceAll(src, "\n", " ")
	out = regexp.MustCompile(`<[^>]+>`).ReplaceAllString(out, "")
	return strings.TrimSpace(out)
}
