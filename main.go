package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/alecthomas/repr"
)

var (
	rep   = regexp.MustCompile(`([^*])/(\r\n|\r|\n)[\t ]*`)
	repLF = regexp.MustCompile(`(\r\n|\r|\n)`)
)

func main() {
	const dirname = "files"

	finfo, err := ioutil.ReadDir(dirname)
	if err != nil {
		panic(err)
	}

	for _, f := range finfo {
		b, err := os.ReadFile(filepath.Join(dirname, f.Name()))
		if err != nil {
			panic(err)
		}

		src := rep.ReplaceAllString(string(b), "$1")
		src = repLF.ReplaceAllString(src, "\n")

		if err := os.WriteFile(filepath.Join(dirname, "replaced_"+f.Name()), []byte(src), 0644); err != nil {
			panic(err)
		}

		actual, err := parser.ParseString("", string(src))
		repr.Println(actual)
		if err != nil {
			fmt.Println("NG:", f.Name())
			log.Fatal(err)
		} else {
			fmt.Println("OK:", f.Name())
		}
	}
}
