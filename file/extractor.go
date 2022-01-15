package file

import "io/ioutil"

func ExtractCodeOfFile(filePath string) (string, error) {
	if filePath == "" {
		return "", nil
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	code := string(content)
	return code, nil
}
