package internal

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

var scriptsLogger *log.Logger = log.New(os.Stdout, "scripts: ", log.LstdFlags|log.Lshortfile)

func CreateRunners() error {
	scriptsLogger.Println("Creating runners...")

	if err := exec.Command("/bin/bash", "/kira/scripts/create-runners.sh").Run(); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			scriptsLogger.Println("Runners already exist.")
		} else {
			return err
		}
	}

	scriptsLogger.Println("Runners created successfully.")
	return nil
}

func CreateUsers() error {
	scriptsLogger.Println("Creating users...")

	if err := exec.Command("/bin/bash", "/kira/scripts/create-users.sh").Run(); err != nil {
		return err
	}

	scriptsLogger.Println("Users created successfully.")
	return nil
}
