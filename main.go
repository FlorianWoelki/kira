package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

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
				},
				Action: func(ctx *cli.Context) error {
					language := ctx.String("language")

					var runner *sandbox.Runner
					for _, r := range sandbox.Runners {
						if language == r.Name {
							runner = &r
							break
						}
					}

					if runner == nil {
						return fmt.Errorf("no language found with name %s", language)
					}

					filePath := ctx.String("file")
					code, err := extractCodeOfFile(filePath)
					if err != nil {
						return fmt.Errorf("something went wrong while reading the file %s", filePath)
					}

					s, output, err := sandbox.Run(runner, code)
					if err != nil {
						return fmt.Errorf("something went wrong while executing sandbox runner %s", err)
					}
					defer s.Clean()

					fmt.Println("\n=== BUILD OUTPUT ===")
					fmt.Printf("Error: %s, Body: %s\n\n", strconv.FormatBool(output.BuildError), output.BuildBody)
					fmt.Println("=== RUN OUTPUT ===")
					fmt.Printf("Error: %s, Body: %s\n", strconv.FormatBool(output.RunError), output.RunBody)
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

func extractCodeOfFile(filePath string) (string, error) {
	if filePath == "" {
		return "", nil
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	code := string(content)
	return code, nil
}
