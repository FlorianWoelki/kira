package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/florianwoelki/kira/file"
	"github.com/florianwoelki/kira/sandbox"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "kira",
		Usage: "Use the kira code execution engine to run your code.",
		Commands: []*cli.Command{
			{
				Name:  "execute",
				Usage: "Execute a test kira code",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "language",
						Aliases: []string{"l", "lang"},
						Value:   "python",
						Usage:   "set the language for the kira sandbox runner.",
					},
					&cli.StringFlag{
						Name:        "file",
						Aliases:     []string{"f"},
						Value:       "",
						DefaultText: "execute an example code.",
						Usage:       "set the specific file that should be executed.",
					},
					&cli.StringFlag{
						Name:        "tests",
						Aliases:     []string{"t"},
						Value:       "",
						DefaultText: "execute based on tests.",
						Usage:       "set the specific path of tests that are being executed and checked on.",
					},
				},
				Action: func(ctx *cli.Context) error {
					language := strings.ToLower(ctx.String("language"))

					var sandboxLang *sandbox.Language
					for _, l := range sandbox.Languages {
						if language == l.Name {
							sandboxLang = &l
							break
						}
					}

					if sandboxLang == nil {
						return fmt.Errorf("no language found with name %s", language)
					}

					filePath := ctx.String("file")
					code, err := file.ExtractCodeOfFile(filePath)
					if err != nil {
						return fmt.Errorf("something went wrong while reading the file %s", filePath)
					}

					testsPath := ctx.String("tests")
					sandboxTests := make([]sandbox.SandboxTest, 0)

					if testsPath != "" {
						testFiles, err := os.ReadDir(testsPath)
						if err != nil {
							return fmt.Errorf("something went wrong while reading the tests path %s", testsPath)
						}

						for _, testFile := range testFiles {
							if !strings.Contains(strings.ToLower(testFile.Name()), "test") {
								continue
							}

							testCode, err := file.ExtractCodeOfFile(testsPath + testFile.Name())
							if err != nil {
								return fmt.Errorf("something went wrong while reading the file %s", testsPath+testFile.Name())
							}

							sandboxTests = append(sandboxTests, sandbox.SandboxTest{Code: []byte(testCode), FileName: testFile.Name()})
						}
					}

					s, output, err := sandbox.Run(sandboxLang, code, sandboxTests)
					if err != nil {
						return fmt.Errorf("something went wrong while executing sandbox runner %s", err)
					}
					defer s.Clean()

					fmt.Println("\n=== BUILD OUTPUT ===")
					fmt.Printf("Error: %s, Body: %s\n\n", strconv.FormatBool(output.BuildError), output.BuildBody)
					fmt.Println("=== RUN OUTPUT ===")
					fmt.Printf("Error: %s, Body: %s\n", strconv.FormatBool(output.RunError), output.RunBody)
					fmt.Println("=== TEST OUTPUT ===")
					fmt.Printf("Error: %s, Body: \n%s\n", strconv.FormatBool(output.TestError), output.TestBody)

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
