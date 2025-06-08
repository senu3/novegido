//go:build headless
// +build headless

package script

import (
	"os"
	"testing"
)

func TestParseDialogue(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"newline", "Hello\nWorld", "Hello World"},
		{"simple tag", "<b>Hello</b>", "Hello"},
		{"tag with text", "<i>Hello</i> world", "Hello world"},
		{"mixed", "<p>Hello\n<b>World</b></p>", "Hello World"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseDialogue(tt.input)
			if got != tt.want {
				t.Errorf("ParseDialogue(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestLoadScriptsChoices(t *testing.T) {
	tmp, err := os.CreateTemp("", "script*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())
	data := `[{"dialogue":{"speaker":"A","text":"hi"},"choices":[{"text":"go","page":0}]}]`
	if _, err := tmp.WriteString(data); err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	pages, err := LoadScripts(tmp.Name())
	if err != nil {
		t.Fatalf("LoadScripts error: %v", err)
	}
	if len(pages) != 1 {
		t.Fatalf("expected 1 page, got %d", len(pages))
	}
	if len(pages[0].Choices) != 1 {
		t.Fatalf("expected 1 choice, got %d", len(pages[0].Choices))
	}
	c := pages[0].Choices[0]
	if c.Text != "go" || c.Page != 0 {
		t.Fatalf("unexpected choice parsed: %+v", c)
	}
}

func TestLoadScriptsTransitions(t *testing.T) {
	tmp, err := os.CreateTemp("", "script*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())
	data := `[{"stage":{"bg":"b.png","bgFade":10,"spriteFade":5}}]`
	if _, err := tmp.WriteString(data); err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	pages, err := LoadScripts(tmp.Name())
	if err != nil {
		t.Fatalf("LoadScripts error: %v", err)
	}
	if len(pages) != 1 {
		t.Fatalf("expected 1 page, got %d", len(pages))
	}
	st := pages[0].Stage
	if st == nil || st.BGFade != 10 || st.SpriteFade != 5 {
		t.Fatalf("transitions not parsed: %+v", st)
	}
}
