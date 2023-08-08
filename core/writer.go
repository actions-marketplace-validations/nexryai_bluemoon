package core

import (
	"io/ioutil"
	"os"
)

func WriteToFile(content string, path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		ExitOnError(err, "Failed to create file.")
		defer file.Close()
	}

	err := ioutil.WriteFile(path, []byte(content), os.ModePerm)
	ExitOnError(err, "Failed to write rules.")

}
