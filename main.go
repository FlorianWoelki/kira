package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/florianwoelki/kira/sandbox"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "kira"
	app.Usage = "Use the kira code execution engine to run your code."
	app.Commands = []cli.Command{
		{
			Name:  "execute",
			Usage: "Execute a test kira code",
			Action: func(ctx *cli.Context) error {
				execute()
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func execute() {
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

	code := `var i = 42; console.log("Hello World " + i);`

	s, err := sandbox.NewSandbox("javascript", []byte(code))
	if err != nil {
		panic(err)
	}
	defer s.Clean()

	stopTicking := make(chan bool)
	go func() {
		for t := range time.Tick(time.Second * 1) {
			select {
			case <-stopTicking:
				return
			default:
				fmt.Println("ticking", t)
				h, _ := time.ParseDuration("5s")
				expireTime := s.LastTimestamp.Add(h)
				if expireTime.Before(time.Now()) {
					s.Clean()
				}
			}
		}
	}()

	output, err := s.Run()
	if err != nil {
		panic(err)
	}

	for _, op := range output {
		fmt.Println(op)
	}

	stopTicking <- true
}
