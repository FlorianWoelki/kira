package sandbox

type runner struct {
	Name            string
	Ext             string
	Image           string
	BuildCmd        string
	RunCmd          string
	DefaultFileName string
	Env             []string
	MaxCPUs         int64
	MaxMemory       int64
}

var runners = []runner{
	{
		Name:            "python",
		Ext:             ".py",
		Image:           "python:3.9.1-alpine",
		BuildCmd:        "",
		RunCmd:          "python3 code.py",
		Env:             []string{},
		DefaultFileName: "code.py",
		MaxCPUs:         2,
		MaxMemory:       128,
	},
	{
		Name:            "go",
		Ext:             ".go",
		Image:           "golang:1.17-alpine",
		BuildCmd:        "rm -rf go.mod && go mod init kira && go build -v .",
		RunCmd:          "./kira",
		Env:             []string{"GOPROXY=https://goproxy.io,direct"},
		DefaultFileName: "code.go",
		MaxCPUs:         2,
		MaxMemory:       128,
	},
	{
		Name:            "c",
		Ext:             ".c",
		Image:           "gcc:latest",
		BuildCmd:        "gcc -v code.c -o code",
		RunCmd:          "./code",
		Env:             []string{},
		DefaultFileName: "code.c",
		MaxCPUs:         2,
		MaxMemory:       128,
	},
	{
		Name:            "java",
		Ext:             ".java",
		Image:           "openjdk:8u232-jdk",
		BuildCmd:        "javac code.java",
		RunCmd:          "java code",
		Env:             []string{},
		DefaultFileName: "code.java",
		MaxCPUs:         2,
		MaxMemory:       128,
	},
	{
		Name:            "javascript",
		Ext:             ".js",
		Image:           "node:lts-alpine",
		BuildCmd:        "",
		RunCmd:          "node code.js",
		Env:             []string{},
		DefaultFileName: "code.js",
		MaxCPUs:         2,
		MaxMemory:       128,
	},
}
