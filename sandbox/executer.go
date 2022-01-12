package sandbox

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func Run(runner *Runner, code string) {
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

	if code == "" {
		code = runner.ExampleCode
	}

	s, err := NewSandbox(runner.Name, []byte(code))
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
				h, _ := time.ParseDuration("30s")
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

	fmt.Println("\n=== BUILD OUTPUT ===")
	fmt.Printf("Error: %s, Body: %s\n\n", strconv.FormatBool(output[0].Error), output[0].Body)
	fmt.Println("=== RUN OUTPUT ===")
	fmt.Printf("Error: %s, Body: %s\n", strconv.FormatBool(output[1].Error), output[1].Body)

	stopTicking <- true
}
