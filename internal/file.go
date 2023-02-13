package internal

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

const tempDirName = "tmp"

var fileLogger *log.Logger = log.New(os.Stdout, "file: ", log.LstdFlags|log.Lshortfile)

// CreateTempDir creates a temporary directory for the user at `/tmp/<user>/<dirName>`.
func CreateTempDir(user, dirName string) error {
	tempDirName := fmt.Sprintf("/%s/%s/%s", tempDirName, user, dirName)
	if err := exec.Command("runuser", "-u", user, "--", "mkdir", tempDirName).Run(); err != nil {
		fileLogger.Println("Could not create temp directory.")
		return err
	}

	return nil
}

// ExecutableFile returns the path to the executable file for the user.
func ExecutableFile(user, dirname, filename string) string {
	return fmt.Sprintf("/%s/%s/%s/%s", tempDirName, user, dirname, filename)
}

// CreateTempFile creates a temporary file for the user at location
// `/tmp/<user>/<dirName>/<filename><extension>`.
func CreateTempFile(user, dirName, filename, extension string) (string, error) {
	fn := fmt.Sprintf("/%s/%s/%s/%s%s", tempDirName, user, dirName, filename, extension)

	if err := exec.Command("runuser", "-u", user, "--", "touch", fn).Run(); err != nil {
		fileLogger.Println("Could not create temp file.")
		return "", err
	}

	return fn, nil
}

// WriteToFile writes the code to the file at the given path. If the file does not exist,
// it will be created. If the file exists, the code will be appended to the file.
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

// DeleteAll deletes all directories and files for the user at location `/tmp/<user>`.
func DeleteAll(user string) error {
	path := fmt.Sprintf("/%s/%s", tempDirName, user)
	err := os.RemoveAll(path)
	if err != nil {
		fileLogger.Println("Could not delete all directories.")
		return err
	}

	return nil
}

// DeleteTempDir deletes the temporary directory for the user at location
// `/tmp/<user>/<dirName>`.
func DeleteTempDir(user, dirName string) error {
	tempDirName := fmt.Sprintf("/%s/%s/%s", tempDirName, user, dirName)
	err := exec.Command("rm", "-rf", tempDirName).Run()
	if err != nil {
		fileLogger.Println("Could not delete temp directory.")
		return err
	}

	return nil
}
