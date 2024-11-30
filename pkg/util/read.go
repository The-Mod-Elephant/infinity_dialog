package util

import (
	"io/fs"
	"os"
	"path/filepath"
)

func GetFiles(path string, ext string) []fs.FileInfo {
	out := []fs.FileInfo{}
	f, err := os.Open(path)
	if err != nil {
		return out
	}
	defer f.Close()
	files, err := f.Readdir(0)
	if err != nil {
		return out
	}
	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ext {
			out = append(out, f)
		}
	}
	return out
}

func ReadFile(path string) (*[]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func ReadFileToString(path string) (string, error) {
	data, err := ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(*data), nil
}
