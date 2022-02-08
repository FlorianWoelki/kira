package sandbox

type Language struct {
	Name            string
	Ext             string
	Image           string
	BuildCmd        string
	RunCmd          string
	DefaultFileName string
	Env             []string
	MaxCPUs         int64
	MaxMemory       int64
	TestCommand     string
	ExampleCode     string
}

var Languages = []Language{
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
		TestCommand:     "python3 -m unittest example_test",
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
		TestCommand:     "go test -v ./...",
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
	{
		Name:            "typescript",
		Ext:             ".ts",
		Image:           "florianwoelki/kira-typescript",
		BuildCmd:        "tsc code.ts",
		RunCmd:          "node code.js",
		Env:             []string{},
		DefaultFileName: "code.ts",
		MaxCPUs:         2,
		MaxMemory:       128,
		ExampleCode:     `console.log("Hello World");`,
	},
	{
		Name:            "julia",
		Ext:             ".jl",
		Image:           "julia:1.7.1-alpine",
		BuildCmd:        "",
		RunCmd:          "julia code.jl",
		Env:             []string{},
		DefaultFileName: "code.jl",
		MaxCPUs:         2,
		MaxMemory:       128,
		ExampleCode:     `print("Hello World")`,
	},
	{
		Name:            "cpp",
		Ext:             ".cpp",
		Image:           "gcc:latest",
		BuildCmd:        "gcc -v code.cpp -lstdc++ -o code",
		RunCmd:          "./code",
		Env:             []string{},
		DefaultFileName: "code.cpp",
		MaxCPUs:         2,
		MaxMemory:       128,
		ExampleCode: `#include <iostream>

int main()
{
  std::cout << "Hello World" << std::endl;
  return 0;
}`,
	},
	{
		Name:            "elixir",
		Ext:             ".exs",
		Image:           "elixir:1.13.1-alpine",
		BuildCmd:        "",
		RunCmd:          "elixir code.exs",
		Env:             []string{},
		DefaultFileName: "code.exs",
		MaxCPUs:         2,
		MaxMemory:       128,
		ExampleCode:     `IO.puts "Hello World"`,
	},
	{
		Name:            "swift",
		Ext:             ".swift",
		Image:           "swift:5.5.2",
		BuildCmd:        "",
		RunCmd:          "swift -module-cache-path . code.swift",
		Env:             []string{},
		DefaultFileName: "code.swift",
		MaxCPUs:         2,
		MaxMemory:       128,
		ExampleCode:     `print("Hello World")`,
	},
}
