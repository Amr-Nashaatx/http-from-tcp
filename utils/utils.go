package utils

import (
	"iter"
	"strings"
	"text/scanner"
)

func GetTokensFromText(text string) iter.Seq[string] {
	var s scanner.Scanner
	return func(yield func(string) bool) {
		s.Init(strings.NewReader(text))
		for token := s.Scan(); token != scanner.EOF; token = s.Scan() {
			if !yield(s.TokenText()) {
				return
			}
		}
	}
}
