package main

import (
	"errors"
	"strings"

	"github.com/pemistahl/lingua-go"
)

// DetectUsedLanguageInISOCode639_1 detects the used language in the provided string
// and returns the ISO 639-1 code of the detected language.
//
// Parameters:
// - str: the string to detect the language from.
//
// Returns:
// - string: the ISO 639-1 code of the detected language.
// - error: an error if the language cannot be detected from the provided string.
func DetectUsedLanguageInISOCode639_1(str string) (string, error) {
	detector := lingua.NewLanguageDetectorBuilder().
		FromAllSpokenLanguages().
		WithPreloadedLanguageModels().
		Build()

	if language, exists := detector.DetectLanguageOf(str); exists {
		return language.IsoCode639_1().String(), nil
	}

	return "", errors.New("can't detect language from provided string")
}


// GetSnowballLanguageFromISOCode639_1 returns the snowball language corresponding to the ISO 639-1 code.
//
// It takes a language code as a parameter and returns the corresponding snowball language as a string.
func GetSnowballLanguageFromISOCode639_1(language string) (newLang string) {
	switch strings.ToUpper(language) {
	case "EN":
		newLang = "english"
	case "ES":
		newLang = "spanish"
	case "FR":
		newLang = "french"
	case "RU":
		newLang = "russian"
	case "SV":
		newLang = "swedish"
	case "NO":
		newLang = "norwegian"
	case "HU":
		newLang = "hungarian"
	default:
		newLang = "english"
		return
	}
	return
}
