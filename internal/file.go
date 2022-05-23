package internal

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

var fileLogger *log.Logger = log.New(os.Stdout, "file: ", log.LstdFlags|log.Lshortfile)

func CreateTempDir(user, dirName string) error {
	tempDirName := fmt.Sprintf("/tmp/%s/%s", user, dirName)
	if err := exec.Command("runuser", "-u", user, "--", "mkdir", tempDirName).Run(); err != nil {
		fileLogger.Println("Could not create temp directory.")
		return err
	}

	return nil
}

func CreateTempFile(user, dirName, extension string) (string, error) {
	filename := fmt.Sprintf("/tmp/%s/%s/code%s", user, dirName, extension)

	if err := exec.Command("runuser", "-u", user, "--", "touch", filename).Run(); err != nil {
		fileLogger.Println("Could not create temp file.")
		return "", err
	}

	return filename, nil
}

func WriteToFile(filename, code string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		fileLogger.Printf("Could not write to file %s.", filename)
		return err
	}

	defer file.Close()

	if _, err := file.WriteString(code); err != nil {
		fileLogger.Println("Could not write line to file.")
		return err
	}

	return nil
}

func DeleteTempDir(user, dirName string) error {
	tempDirName := fmt.Sprintf("/tmp/%s/%s", user, dirName)
	err := exec.Command("rm", "-rf", tempDirName).Run()
	if err != nil {
		fileLogger.Println("Could not delete temp directory.")
		return err
	}

	return nil
}
