package file

import (
	"bufio"
	"os"
)

func Reader(filePath string) (*bufio.Scanner, func(), error) {
	readFile, err := os.Open(filePath)
	close := func() {
		readFile.Close()
	}

	if err != nil {
		return nil, close, err
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	return fileScanner, close, nil
}
