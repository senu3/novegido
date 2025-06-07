package main

import "testing"

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
