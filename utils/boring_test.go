package utils

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestSplitMessage(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		limit int
	}{
		{"empty", "", 2000},
		{"short", "hello world", 2000},
		{"exact limit", strings.Repeat("a", 2000), 2000},
		{"over limit no breaks", strings.Repeat("a", 4500), 2000},
		{"over limit with newlines", strings.Repeat("line of text\n", 500), 2000},
		{"multibyte runes", strings.Repeat("käyttäjä 日本語 ", 400), 2000},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			chunks := SplitMessage(tc.text, tc.limit)

			// Reassembly must equal the original.
			if got := strings.Join(chunks, ""); got != tc.text {
				t.Fatalf("reassembled text differs from original (len got %d, want %d)", len(got), len(tc.text))
			}

			for i, c := range chunks {
				if n := utf8.RuneCountInString(c); n > tc.limit {
					t.Errorf("chunk %d has %d runes, exceeds limit %d", i, n, tc.limit)
				}
				if !utf8.ValidString(c) {
					t.Errorf("chunk %d contains an invalid/broken UTF-8 rune", i)
				}
			}

			if tc.text == "" && len(chunks) != 0 {
				t.Errorf("empty input should produce no chunks, got %d", len(chunks))
			}
		})
	}
}
