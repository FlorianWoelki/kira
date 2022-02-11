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
						Name:        "main",
						Aliases:     []string{"m"},
						Value:       "",
						DefaultText: "default main file from example code.",
						Usage:       "set the specific main file for executing first.",
					},
					&cli.StringFlag{
						Name:        "dir",
						Aliases:     []string{"d"},
						Value:       "",
						DefaultText: "copies all the example code.",
						Usage:       "set the specific directory that should be copied and executed from.",
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

					mainFilePath := ctx.String("main")
					mainCode, err := file.ExtractCodeOfFile(mainFilePath)
					if err != nil {
						return fmt.Errorf("something went wrong while reading the main file path %s", mainFilePath)
					}

					dirPath := ctx.String("dir")
					sandboxTests := make([]sandbox.SandboxFile, 0)
					files := make([]sandbox.SandboxFile, 0)

					if dirPath != "" {
						dirFiles, err := os.ReadDir(dirPath)
						if err != nil {
							return fmt.Errorf("something went wrong while reading the dir path %s", dirPath)
						}

						for _, dirFile := range dirFiles {
							if dirPath+dirFile.Name() == mainFilePath {
								continue
							}

							fileCode, err := file.ExtractCodeOfFile(dirPath + dirFile.Name())
							if err != nil {
								return fmt.Errorf("something went wrong while reading the file %s", dirPath+dirFile.Name())
							}

							if strings.Contains(strings.ToLower(dirFile.Name()), "test") {
								sandboxTests = append(sandboxTests, sandbox.SandboxFile{Code: []byte(fileCode), FileName: dirFile.Name()})
							} else {
								files = append(files, sandbox.SandboxFile{Code: []byte(fileCode), FileName: dirFile.Name()})
							}
						}
					}

					s, output, err := sandbox.Run(sandboxLang, mainCode, files, sandboxTests)
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

					/*testsPath := ctx.String("tests")
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

					return nil*/
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
