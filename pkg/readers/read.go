package readers

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
)

func GetFiles(path, ext string) []fs.FileInfo {
	out := []fs.FileInfo{}
	f, err := os.Open(filepath.Clean(path))
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
	data, err := os.ReadFile(filepath.Clean(path))
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

func ReadFileToSlice(path string) (*[]string, error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return &lines, scanner.Err()
}
