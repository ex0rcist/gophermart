package entities

import (
	"strings"
)

type Secret string

func (s *Secret) Set(src string) error {
	*s = Secret(src)
	return nil
}

func (s Secret) Type() string {
	return "string"
}

func (s Secret) String() string {
	if len(s) <= 2 {
		return string(s)
	}

	masked := strings.Repeat("*", len(s)-2)
	return string(s[0]) + masked + string(s[len(s)-1])
}
