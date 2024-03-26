package main

import (
	"fmt"
	"os"

	"github.com/kljensen/snowball"
	stopwords "github.com/toadharvard/stopwords-iso"
)

func main() {
	rawString := GetStringFromCLI()
	usedLanguageCode, err := DetectUsedLanguageInISOCode639_1(rawString)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Printf("Detected language: %s\n", usedLanguageCode)

	m, err := stopwords.NewStopwordsMapping()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	stringWithoutStopwords := m.ClearStringByLang(rawString, usedLanguageCode)
	fmt.Printf("String without stopwords: %s\n", stringWithoutStopwords)

	snowballLang := GetSnowballLanguageFromISOCode639_1(usedLanguageCode)
	stemmedString, err := snowball.Stem(stringWithoutStopwords, snowballLang, false)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Printf("Stemmed string: %s\n", stemmedString)
	fmt.Printf("Unique words: %s\n", StringToWords(stemmedString))
}
