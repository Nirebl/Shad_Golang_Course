package language

import (
	"encoding/json"
	"strings"
)

type Language struct {
	Name       string
	Extensions []string
}

var langMap = make(map[string]Language)

func init() {
	var languages []Language
	_ = json.Unmarshal(languagesJSON, &languages)

	for _, lang := range languages {
		langMap[strings.ToLower(lang.Name)] = lang
	}
}

func GetLanguageExtensions(lang string) []string {
	return langMap[strings.ToLower(lang)].Extensions
}
