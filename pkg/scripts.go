package pkg

import (
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var scriptsLogger *log.Logger = log.New(os.Stdout, "scripts: ", log.LstdFlags|log.Lshortfile)

// CreateRunners execute the `/scripts/create-runners.sh` script which creates the runners
// group.
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

// CreateUsers tries to run the `/scripts/create-users.sh` script which adds the linux
// users to the machine.
func CreateUsers() error {
	scriptsLogger.Println("Creating users...")

	if err := exec.Command("/bin/bash", "/kira/scripts/create-users.sh").Run(); err != nil {
		return err
	}

	scriptsLogger.Println("Users created successfully.")
	return nil
}

// CreateBinaries gets all the `install.sh` scripts from the activated languages and tries
// to install all their binaries.
func CreateBinaries() error {
	scriptsLogger.Println("Creating binaries...")
	err := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		dir := filepath.Base(filepath.Dir(path))
		if _, ok := LoadedLanguages[dir]; ok && strings.HasSuffix(path, "install.sh") {
			scriptsLogger.Printf("Downloading %s binaries...", dir)
			err := exec.Command("/bin/bash", path).Run()
			if err != nil {
				return err
			}

			scriptsLogger.Printf("%s binaries downloaded successfully.", dir)
		}

		return nil
	})

	if err != nil {
		return err
	}

	scriptsLogger.Println("Binaries have been downloaded successfully.")
	return nil
}
