package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const dirname = "ghost"
const destDir = "out"

func findFilesWithWalkDir(root, ext string) ([]string, error) {
	findList := []string{}

	err := filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ext) {
			return nil
		}

		findList = append(findList, path)
		return nil
	})
	return findList, err
}

func TestMain(*testing.T) {
	files, err := findFilesWithWalkDir(dirname, "dic")
	if err != nil {
		panic(err)
	}

	fmt.Println(files)

	for _, src := range files {
		dest := strings.Replace(src, dirname, destDir, 1)

		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			panic(err)
		}

		if err := os.WriteFile(dest, []byte(parse(src).String()), 0644); err != nil {
			panic(err)
		}
	}
}
