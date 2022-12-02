package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alecthomas/repr"
)

var (
	rep        = regexp.MustCompile(`([^*])/(\r\n|\r|\n)[\t ]*`)
	repLF      = regexp.MustCompile(`(\r\n|\r|\n)`)
	repFor     = regexp.MustCompile(`for([^\n{;]*);([^\n{;]*);`)
	repForEach = regexp.MustCompile(`foreach([^;]*);`)

	// ;を置換する際の一時的な文字 ファイル内に存在するどの文字とも被ってはいけない
	tmpColon = string(rune(0x12))
)

func main() {
	const dirname = "files"

	finfo, err := ioutil.ReadDir(dirname)
	if err != nil {
		panic(err)
	}

	for _, f := range finfo {
		if strings.HasPrefix(f.Name(), "replaced_") {
			continue
		}

		b, err := os.ReadFile(filepath.Join(dirname, f.Name()))
		if err != nil {
			panic(err)
		}

		src := rep.ReplaceAllString(string(b), "$1")
		src = repLF.ReplaceAllString(src, "\n")

		src = repFor.ReplaceAllString(src, "for$1"+tmpColon+"$2"+tmpColon)
		src = repForEach.ReplaceAllString(src, "foreach$1"+tmpColon)

		src = strings.ReplaceAll(src, ";", "\n")
		src = strings.ReplaceAll(src, tmpColon, ";")

		if err := os.WriteFile(filepath.Join(dirname, "replaced_"+f.Name()), []byte(src), 0644); err != nil {
			panic(err)
		}

		actual, err := parser.ParseString("", string(src))
		if err != nil {
			repr.Println(actual)
			fmt.Println("NG:", f.Name())
			log.Fatal(err)
		} else {
			fmt.Println("OK:", f.Name())
		}
	}
}
