package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

const dirname = "testfiles"

func TestMain(*testing.T) {
	finfo, err := ioutil.ReadDir(dirname)
	if err != nil {
		panic(err)
	}

	for i, f := range finfo {
		if i >= 2 {
			break
		}

		if strings.HasPrefix(f.Name(), "replaced_") {
			continue
		}

		fmt.Println(f.Name())

		fmt.Println(parse(filepath.Join(dirname, f.Name())))

		// if err := os.WriteFile("out.txt", []byte(), 0644); err != nil {
		// 	panic(err)
		// }
	}
}
