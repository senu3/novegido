package main

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"
)

// --- JSON に合わせたデータ構造 -----------------------------------------

type SpriteInfo struct {
	ID   string `json:"id"`
	File string `json:"file"`
	Pos  string `json:"pos"` // "left" / "center" / "right" など
}

type StageInfo struct {
	BG      string       `json:"bg"`
	Sprites []SpriteInfo `json:"sprites"`
}

type DialogueInfo struct {
	Speaker string `json:"speaker"`
	Text    string `json:"text"`
}

type Page struct {
	Stage    *StageInfo    `json:"stage,omitempty"`
	Dialogue *DialogueInfo `json:"dialogue,omitempty"`
	Clean    string        `json:"-"` // 描画用にキャッシュ
}

// --- ローダ／パーサ ------------------------------------------------------

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

// ParseDialogue は BBCode やタグを除去するなど、想定フォーマットに応じて実装してください
func ParseDialogue(src string) string {
	// ここでは単純に\n → 空白、 <>タグを除去するだけのダミー処理
	out := strings.ReplaceAll(src, "\n", " ")
	out = regexp.MustCompile(`<[^>]+>`).ReplaceAllString(out, "")
	return strings.TrimSpace(out)
}
