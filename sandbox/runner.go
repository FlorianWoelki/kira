package sandbox

type Runner struct {
	Name            string
	Ext             string
	Image           string
	BuildCmd        string
	RunCmd          string
	DefaultFileName string
	Env             []string
	MaxCPUs         int64
	MaxMemory       int64
	ExampleCode     string
}

var Runners = []Runner{
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
		ExampleCode:     `print("Hello World")`,
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
		ExampleCode: `package main

import "fmt"

func main() {
	fmt.Println("Hello World")
}`,
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
		ExampleCode: `#include <stdio.h>

int main()
{
	printf("Hello World");
	return 0;
}`,
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
		ExampleCode: `class code {
	public static void main(String[] args) {
		System.out.println("Hello World");
	}
}`,
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
		ExampleCode:     `console.log("Hello World");`,
	},
}
