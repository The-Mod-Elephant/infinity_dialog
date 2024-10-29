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

func ReadFile(path string, file fs.FileInfo) (string, error) {
	data, err := os.ReadFile(filepath.Join(path, file.Name()))
	if err != nil {
		println("DIED")
		return "", err
	}
	return string(data), nil
}
