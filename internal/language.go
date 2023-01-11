package internal

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/florianwoelki/kira/internal/pool"
)

var languageLogger *log.Logger = log.New(os.Stdout, "language: ", log.LstdFlags|log.Lshortfile)

var LoadedLanguages map[string]Language

type Language struct {
	Name      string `json:"name" binding:"required"`
	Version   string `json:"version" binding:"required"`
	Extension string `json:"extension" binding:"required"`
	Timeout   int    `json:"timeout" binding:"required"`
	TestInfo  struct {
		Regex           string `json:"regex,omitempty"`
		FailedTestRegex string `json:"failedTestRegex,omitempty"`
		AssertionRegex  string `json:"assertionRegex,omitempty"`
		PassedString    string `json:"passedString,omitempty"`
	} `json:"testInfo,omitempty"`
	Compiled bool
}

func LoadLanguages(activeLanguages []string) error {
	languageLogger.Println("Loading languages...")
	LoadedLanguages = make(map[string]Language)

	err := filepath.Walk("./languages", func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, "metadata.json") {
			fileBytes, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			language := Language{}
			if err = json.Unmarshal(fileBytes, &language); err != nil {
				return err
			}

			dir := filepath.Dir(path)
			_, err = os.Stat(fmt.Sprintf("%s/%s", dir, "compile.sh"))
			language.Compiled = err == nil

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

func GetLanguageByName(key string) (Language, error) {
	find, ok := LoadedLanguages[strings.ToLower(key)]
	if !ok {
		return Language{}, fmt.Errorf("could not find language with key: %s", key)
	}

	return find, nil
}

func PrettifyTestOutput(output string, language Language) []*pool.TestResult {
	regex := regexp.MustCompile(language.TestInfo.Regex)
	failedTestRegex := regexp.MustCompile(language.TestInfo.FailedTestRegex)
	assertionRegex := regexp.MustCompile(language.TestInfo.AssertionRegex)
	lines := strings.Split(output, "\n")

	formattedOutput := map[string]*pool.TestResult{}

	lastFailedTest := ""
	for _, line := range lines {
		regexMatch := regex.FindStringSubmatch(line)
		failedTestMatch := failedTestRegex.FindStringSubmatch(line)

		if len(regexMatch) > 0 {
			name := regexMatch[1]
			status := regexMatch[2]
			formattedOutput[name] = &pool.TestResult{
				Name:   name,
				Passed: status == language.TestInfo.PassedString,
			}
		} else if len(failedTestMatch) > 0 {
			lastFailedTest = failedTestMatch[1]
		}

		if len(lastFailedTest) > 0 {
			assertionMatch := assertionRegex.FindStringSubmatch(line)
			if len(assertionMatch) > 0 {
				formattedOutput[lastFailedTest].Received = assertionMatch[1]
				formattedOutput[lastFailedTest].Actual = assertionMatch[2]
				lastFailedTest = ""
			}
		}
	}

	result := []*pool.TestResult{}
	for _, data := range formattedOutput {
		result = append(result, data)
	}

	return result
}
