package service

import (
	"strings"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

func ConvertText(text string) string {
	f := func(r rune) bool {
		return r != '.' && r != '-' && r != ' '
	}

	if strings.ContainsFunc(text, f) {
		return morse.ToMorse(text)
	} else {
		return morse.ToText(text)
	}

}
