package internal

import (
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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

func CreateBinaries() error {
	scriptsLogger.Println("Creating binaries...")
	err := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, "install.sh") {
			dir := filepath.Base(filepath.Dir(path))
			scriptsLogger.Printf("Downloading %s binaries...", dir)
			runInstallScript(path, dir)
		}

		return nil
	})

	if err != nil {
		return err
	}

	scriptsLogger.Println("Binaries have been downloaded successfully.")
	return nil
}

func runInstallScript(path, binary string) error {
	err := exec.Command("bash", path).Run()
	if err != nil {
		return err
	}

	scriptsLogger.Printf("%s binaries downloaded successfully.", binary)
	return nil
}
