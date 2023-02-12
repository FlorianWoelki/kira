package pkg

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var languageLogger *log.Logger = log.New(os.Stdout, "language: ", log.LstdFlags|log.Lshortfile)

// LoadedLanguages contains all the languages with an identifier as a key and the struct
// as a value.
var LoadedLanguages map[string]Language

type Language struct {
	// Name of the language.
	Name string `json:"name"`
	// Version that is currently used.
	Version string `json:"version"`
	// Extension of the to be executed files for the language.
	Extension string `json:"extension"`
	// Timeout of the execution, when the code will be terminated.
	Timeout int `json:"timeout"`
	// Compiled set whether the language needs to be compiled beforehand.
	Compiled bool `json:"compiled"`
}

// LoadLanguages load all the specified active languages and define all the neccessary
// information for the `Language` struct.
func LoadLanguages(activeLanguages []string) error {
	languageLogger.Println("Loading languages...")
	LoadedLanguages = make(map[string]Language)

	err := filepath.Walk("./languages", func(path string, info fs.FileInfo, err error) error {
		// Only get the information from the `metadata.json` file.
		if strings.HasSuffix(path, "metadata.json") {
			fileBytes, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			language := Language{}
			if err = json.Unmarshal(fileBytes, &language); err != nil {
				return err
			}

			// Check if the language is a compiled language.
			dir := filepath.Dir(path)
			_, err = os.Stat(fmt.Sprintf("%s/%s", dir, "compile.sh"))
			language.Compiled = err == nil

			// Check if the language is in the defined active languages.
			shouldInsert := true
			if len(activeLanguages) != 0 {
				shouldInsert = false
				for _, activeLanguage := range activeLanguages {
					if strings.EqualFold(activeLanguage, language.Name) {
						shouldInsert = true
					}
				}
			}

			if shouldInsert {
				LoadedLanguages[strings.ToLower(language.Name)] = language
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	languageLogger.Printf("Languages successfully loaded (amount: %d).\n", len(LoadedLanguages))
	return nil
}

// GetLanguages gets all loaded languages as an array.
func GetLanguages() ([]Language, error) {
	if len(LoadedLanguages) == 0 {
		return nil, fmt.Errorf("could not find any to be loaded languages")
	}

	result := make([]Language, 0)
	for _, languageValue := range LoadedLanguages {
		result = append(result, languageValue)
	}

	return result, nil
}

// GetLanguageByName gets a language by name from the loaded languages.
func GetLanguageByName(key string) (Language, error) {
	find, ok := LoadedLanguages[strings.ToLower(key)]
	if !ok {
		return Language{}, fmt.Errorf("could not find language with key: %s", key)
	}

	return find, nil
}
