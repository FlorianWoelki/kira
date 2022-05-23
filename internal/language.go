package internal

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var languageLogger *log.Logger = log.New(os.Stdout, "language: ", log.LstdFlags|log.Lshortfile)

var LoadedLanguages map[string]string

type Language struct {
	Name      string `json:"name" binding:"required"`
	Version   string `json:"version" binding:"required"`
	Extension string `json:"extension" binding:"required"`
	Timeout   int    `json:"timeout" binding:"required"`
}

func LoadLanguages() error {
	languageLogger.Println("Loading languages...")
	LoadedLanguages = make(map[string]string)

	err := filepath.Walk("./languages", func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, "metadata.json") {
			fileBytes, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			language := Language{}
			if err = json.Unmarshal(fileBytes, &language); err != nil {
				return err
			}

			LoadedLanguages[strings.ToLower(language.Name)] = string(fileBytes)
		}

		return nil
	})

	if err != nil {
		return err
	}

	languageLogger.Printf("Languages successfully loaded (amount: %d).\n", len(LoadedLanguages))
	return nil
}

func GetLanguages() ([]Language, error) {
	if LoadedLanguages == nil || len(LoadedLanguages) == 0 {
		return nil, fmt.Errorf("No loaded languages.")
	}

	result := make([]Language, 0)
	for _, languageValue := range LoadedLanguages {
		serialized := []byte(languageValue)
		language := Language{}
		err := json.Unmarshal(serialized, &language)

		if err != nil {
			return nil, err
		}

		result = append(result, language)
	}

	return result, nil
}

func GetLanguageByName(key string) (Language, error) {
	find, ok := LoadedLanguages[strings.ToLower(key)]
	if !ok {
		return Language{}, fmt.Errorf("Could not find language with key: %s", find)
	}

	language := Language{}
	if err := json.Unmarshal([]byte(find), &language); err != nil {
		return Language{}, err
	}

	return language, nil
}
