# `libretranslate`

A Go Wrapper for [LibreTranslate](https://libretranslate.com/).

## Installation

```bash
go get -u github.com/piero-vic/libretranslate
```

## Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/piero-vic/libretranslate"
)

func main() {
	lt := libretranslate.NewClient("<your_api_token>")

	// Example 1: Detect Language
	detectedLanguages, err := lt.Detect("Hello, world!")
	if err != nil {
		log.Fatalf("Detection failed: %v", err)
	}

	fmt.Println("Detected languages for 'Hello, world!':")
	for _, detectedLang := range detectedLanguages {
		fmt.Printf("Language: %s, Confidence: %.2f\n", detectedLang.Language, detectedLang.Confidence)
	}

	// Example 2: Translate Text
	sourceLanguage := "en"
	targetLanguage := "fr"
	textToTranslate := "Hello, world!"

	translatedText, err := lt.Translate(textToTranslate, sourceLanguage, targetLanguage)
	if err != nil {
		log.Fatalf("Translation failed: %v", err)
	}

	fmt.Printf("Translation from '%s' to '%s':\n", sourceLanguage, targetLanguage)
	fmt.Printf("Original Text: %s\n", textToTranslate)
	fmt.Printf("Translated Text: %s\n", translatedText)
}
```

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.
