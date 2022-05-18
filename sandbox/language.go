package sandbox

type Language struct {
	Name            string
	Ext             string
	Image           string
	BuildCmd        string
	RunCmd          string
	DefaultFileName string
	Env             []string
	TestCommand     string
	ExampleCode     string
}

var Languages = []Language{
	{
		Name:            "python",
		Ext:             ".py",
		Image:           "all-in-one-ubuntu",
		BuildCmd:        "",
		RunCmd:          "python3 app.py",
		Env:             []string{},
		DefaultFileName: "app.py",
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
		DefaultFileName: "app.go",
		TestCommand:     "go test -v {} {}",
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
		BuildCmd:        "gcc -v app.c -o app",
		RunCmd:          "./app",
		Env:             []string{},
		DefaultFileName: "app.c",
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
		BuildCmd:        "javac app.java",
		RunCmd:          "java app",
		Env:             []string{},
		DefaultFileName: "app.java",
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
		RunCmd:          "node app.js",
		Env:             []string{},
		DefaultFileName: "app.js",
		ExampleCode:     `console.log("Hello World");`,
	},
	{
		Name:            "cpp",
		Ext:             ".cpp",
		Image:           "gcc:latest",
		BuildCmd:        "gcc -v app.cpp -lstdc++ -o app",
		RunCmd:          "./app",
		Env:             []string{},
		DefaultFileName: "app.cpp",
		ExampleCode: `#include <iostream>

int main()
{
  std::cout << "Hello World" << std::endl;
  return 0;
}`,
	},
}
