package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// === 構造体定義 ===

type Page struct {
	Stage    *Stage    `json:"stage,omitempty"`
	Dialogue *Dialogue `json:"dialogue,omitempty"`
}

type Stage struct {
	BG      string   `json:"bg,omitempty"`
	Sprites []Sprite `json:"sprites,omitempty"`
	Effects []string `json:"effects,omitempty"`
}

type Sprite struct {
	ID   string `json:"id"`
	File string `json:"file"`
	Pos  string `json:"pos"`
}

type Dialogue struct {
	Speaker string `json:"speaker,omitempty"`
	Text    string `json:"text,omitempty"`
}

type DictionaryEntry struct {
	Title  string `json:"title"`
	Short  string `json:"short"`
	Detail string `json:"detail"`
}

type Dictionary map[string]DictionaryEntry

// === データ読み込み関数 ===

func LoadPages(path string) ([]Page, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var pages []Page
	err = json.Unmarshal(b, &pages)
	return pages, err
}

func LoadDictionary(path string) (Dictionary, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var dict Dictionary
	err = json.Unmarshal(b, &dict)
	return dict, err
}

// === TIPS用語マークアップのパース ===

var tipRegex = regexp.MustCompile(`\\[\\[(.+?)\\|(.+?)\\]\\]`)

type TipLink struct {
	ID   string
	Text string
	Pos  int
}

func ParseDialogue(raw string) (cleanText string, tips []TipLink) {
	cleanText = tipRegex.ReplaceAllStringFunc(raw, func(m string) string {
		matches := tipRegex.FindStringSubmatch(m)
		tips = append(tips, TipLink{
			ID:   matches[1],
			Text: matches[2],
			Pos:  strings.Index(raw, m),
		})
		return matches[2]
	})
	return
}

// === メイン（動作デモ用） ===

func main() {
	// デモ用の読み込み
	pages, err := LoadPages("asset/scripts/demo.json")
	if err != nil {
		panic(err)
	}
	dict, err := LoadDictionary("asset/dict/dictionary.json")
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
			cleanText, tips := ParseDialogue(page.Dialogue.Text)
			fmt.Printf("%s「%s」\n", page.Dialogue.Speaker, cleanText)
			for _, tip := range tips {
				title, short, ok := dict.GetTip(tip.ID)
				if ok {
					fmt.Printf("  [TIPS] %s: %s\n", title, short)
				} else {
					fmt.Printf("  [TIPS] 用語 %s は辞書に未登録です。\n", tip.ID)
				}
			}
		}
		fmt.Println("-----------")
	}
}

// === 辞書からTIPS取得 ===

func (dict Dictionary) GetTip(id string) (title, short string, ok bool) {
	entry, ok := dict[id]
	return entry.Title, entry.Short, ok
}
