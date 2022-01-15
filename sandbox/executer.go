package sandbox

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type runOutput struct {
	BuildBody  string `json:"buildBody"`
	BuildError bool   `json:"buildError"`
	RunBody    string `json:"runBody"`
	RunError   bool   `json:"runError"`
}

func Run(language *Language, code string) (*Sandbox, runOutput, error) {
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
		code = language.ExampleCode
	}

	s, err := NewSandbox(language.Name, []byte(code))
	if err != nil {
		return &Sandbox{}, runOutput{}, err
	}

	stopTicking := make(chan bool)
	go func() {
		timer := time.NewTicker(time.Second * 1)
		for t := range timer.C {
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
		return &Sandbox{}, runOutput{}, err
	}

	stopTicking <- true
	return s, runOutput{
		BuildBody:  output[0].Body,
		BuildError: output[0].Error,
		RunBody:    output[1].Body,
		RunError:   output[1].Error,
	}, nil
}
