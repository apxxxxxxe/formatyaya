package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const dirname = "testdata"
const destDir = "test_out"

func findFiles(root string, exts []string) ([]string, error) {
	findList := []string{}

	err := filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if err := validateExt(path, exts); err != nil {
			return nil
		}

		findList = append(findList, path)
		return nil
	})
	return findList, err
}

func TestMain(*testing.T) {
	files, err := findFiles(dirname, []string{"dic", "aym"})
	if err != nil {
		panic(err)
	}

	fmt.Println("files:", files)

	for _, src := range files {
		dest := strings.Replace(src, dirname, destDir, 1)

		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			panic(err)
		}

		b, err := os.ReadFile(src)
		if err != nil {
			panic(err)
		}

		fmt.Println(src)
		if err := os.WriteFile(dest, []byte(parse(b).String()), 0644); err != nil {
			panic(err)
		}
	}
}
