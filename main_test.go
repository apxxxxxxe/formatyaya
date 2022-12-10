package main

import (
	"io/ioutil"
	"os"
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

	for _, f := range finfo {
		// if i >= 2 {
		// 	break
		// }

		if strings.HasPrefix(f.Name(), "replaced_") {
			continue
		}

		if err := os.MkdirAll("out", 0755); err != nil {
			panic(err)
		}

		if err := os.WriteFile(filepath.Join("out", "out_"+f.Name()), []byte(parse(filepath.Join(dirname, f.Name())).String()), 0644); err != nil {
			panic(err)
		}
	}
}
