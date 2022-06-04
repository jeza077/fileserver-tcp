package main

import (
	"errors"
	"os"
)

type File struct {
	Name    string
	Size    int64
	Content []byte
}

const MAX_FILE_SIZE = 100000000

func LoadFile(path string) (string, int64, []byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", 0, nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return "", 0, nil, err
	}

	if fileInfo.Size() > MAX_FILE_SIZE {
		return "", 0, nil, errors.New("tamaño maximo por archivo es de 100mb")
	}

	buffer := make([]byte, fileInfo.Size())
	if _, err := file.Read(buffer); err != nil {
		return "", 0, nil, err
	}
	file.Close()

	return fileInfo.Name(), fileInfo.Size(), buffer, nil
}
