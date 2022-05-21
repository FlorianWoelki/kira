package sandbox

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type runOutput struct {
	BuildBody        string        `json:"buildBody"`
	BuildError       bool          `json:"buildError"`
	RunBody          string        `json:"runBody"`
	RunExecutionTime time.Duration `json:"runExecutionTime"`
	RunError         bool          `json:"runError"`
	TestBody         string        `json:"testBody"`
	TestError        bool          `json:"testError"`
}

func Run(language *Language, mainCode string, files []SandboxFile, sandboxTests []SandboxFile) (*Sandbox, runOutput, error) {
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

	if mainCode == "" {
		mainCode = language.ExampleCode
	}

	s, err := NewSandbox(language.Name, []byte(mainCode), files, sandboxTests)
	if err != nil {
		return &Sandbox{}, runOutput{}, err
	}

	stopTicking := make(chan bool)
	go func() {
		timer := time.NewTicker(time.Second * 1)
		for range timer.C {
			select {
			case <-stopTicking:
				return
			default:
				h, _ := time.ParseDuration("30s")
				expireTime := s.LastTimestamp.Add(h)
				if expireTime.Before(time.Now()) {
					s.forceQuit = true
				}
			}
		}
	}()

	output, err := s.Run()
	if err != nil {
		return &Sandbox{}, runOutput{}, err
	}

	stopTicking <- true

	if s.forceQuit {
		output.RunOutput.ExecBody = "Could not execute code. Force Quit."
		output.RunOutput.Error = true
	}

	return s, runOutput{
		BuildBody:        output.SetupOutput.ExecBody,
		BuildError:       output.SetupOutput.Error,
		RunBody:          output.RunOutput.ExecBody,
		RunExecutionTime: output.RunOutput.ExecutionTime,
		RunError:         output.RunOutput.Error,
		TestBody:         output.RunOutput.TestsBody,
	}, nil
}
