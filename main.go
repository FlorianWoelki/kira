package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	code := `package main

import "fmt"
	
func main() {
	fmt.Println("Hello World")
}
`

	s, err := sandbox.NewSandbox("go", []byte(code))
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
