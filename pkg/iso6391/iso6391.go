package iso6391

import (
	"errors"
	"strings"

	iso "github.com/emvi/iso-639-1"
)

type ISOCode6391 struct {
	Code string
	Name string
}

func NewLanguage(code string) (ISOCode6391, error) {
	if !iso.ValidCode(code) {
		return ISOCode6391{}, errors.New("invalid language code")
	}
	name := strings.ToLower(iso.Name(code))
	return ISOCode6391{Code: code, Name: name}, nil
}

func (l *ISOCode6391) String() string {
	return l.Code
}
