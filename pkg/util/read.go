package util

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func ReadFiles(root string, ext string) (map[string]*[]string, error) {
	files := map[string]*[]string{}
	err := filepath.WalkDir(root, func(path string, file fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		file_ext := strings.ToLower(filepath.Ext(file.Name()))
		target_ext := "." + strings.ToLower(ext)
		if !file.IsDir() && file_ext == target_ext {
			fileContent, err := ReadFileToSlice(path)
			if err != nil {
				return err
			}
			files[file.Name()] = fileContent
		}
		return nil
	})
	return files, err
}

func ReadFile(path string) (*[]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func ReadFileToString(path string) (*string, error) {
	data, err := ReadFile(path)
	if err != nil {
		return nil, err
	}
	out := string(*data)
	return &out, nil
}

func ReadFileToSlice(path string) (*[]string, error) {
	file, err := os.Open(path)
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
