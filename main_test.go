package main

import (
	"testing"
)

func TestConvertMarkdownToDocx(t *testing.T) {
	markdown := "# Hello, World!"
	result, err := ConvertMarkdownToDocx(markdown)
	if err != nil {
		t.Errorf("ConvertMarkdownToDocx failed with error: %v", err)
	}

	if len(result) == 0 {
		t.Errorf("ConvertMarkdownToDocx returned empty result")
	}
}