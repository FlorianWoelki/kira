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

func ExecutableFile(user, dirname, filename string) string {
	return fmt.Sprintf("/tmp/%s/%s/%s", user, dirname, filename)
}

func CreateTempFile(user, dirName, filename, extension string) (string, error) {
	fn := fmt.Sprintf("/tmp/%s/%s/%s%s", user, dirName, filename, extension)

	if err := exec.Command("runuser", "-u", user, "--", "touch", fn).Run(); err != nil {
		fileLogger.Println("Could not create temp file.")
		return "", err
	}

	return fn, nil
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

func DeleteAll(user string) error {
	path := fmt.Sprintf("/tmp/%s", user)
	err := os.RemoveAll(path)
	if err != nil {
		fileLogger.Println("Could not delete all directories.")
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
