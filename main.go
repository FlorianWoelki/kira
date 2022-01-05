package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/florianwoelki/kira/sandbox"
)

func main() {
	c := make(chan os.Signal)

	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			default:
				fmt.Println("signal received", s)
			}
		}
	}()

	code := "print(\"Hello World\")"

	s, err := sandbox.NewSandbox("python", []byte(code))
	if err != nil {
		panic(err)
	}

	output, err := s.Run()
	if err != nil {
		panic(err)
	}

	for _, op := range output {
		fmt.Println(op.Body)
	}

	s.Clean()
}
