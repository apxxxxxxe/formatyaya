package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/alecthomas/repr"
)

const version = "0.1.0"

var (
	repSlash         = regexp.MustCompile(`([^*])/\n[\t ]*`)
	repLF            = regexp.MustCompile(`(\r\n|\r)`)
	repTrailingSpace = regexp.MustCompile(`[\t ]+\n`)
)

func parse(filename string) *Root {
	b, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	// 改行コードの統一
	src := repLF.ReplaceAllString(string(b), "\n")
	// 行末の空白文字を削除
	src = repTrailingSpace.ReplaceAllString(src, "\n")
	// 行末のスラッシュ(次行との連結)を処理
	src = repSlash.ReplaceAllString(src, "$1")

	actual, err := parser.ParseString("", string(src))
	if err != nil {
		log.Println(repr.String(actual))
		log.Println(err)
		return nil
	}
	return actual
}

func validateExt(file string, exts []string) error {
	finfo, err := os.Stat(file)
	if err != nil {
		return err
	} else if finfo.IsDir() {
		return errors.New(file + ": is directory")
	}

	for _, e := range exts {
		if strings.HasSuffix(file, e) {
			return nil
		}
	}
	return errors.New(file + ": invalid extension")
}

func main() {
	var (
		inPlace    bool
		useSpace   bool
		spaceCount int
	)
	flag.BoolVar(&inPlace, "i", false, "Inplace: overwrite original file")
	flag.BoolVar(&useSpace, "s", false, "useSpace: indent as space")
	flag.IntVar(&spaceCount, "c", 2, "spaceCount: the count of space when flag -s turned on")

	flag.Parse()
	file := flag.Arg(0)
	exts := []string{"dic", "aym"}

	if err := validateExt(file, exts); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if useSpace {
		Indent = strings.Repeat(" ", spaceCount)
	} else {
		Indent = "	"
	}

	parsed := parse(file)
	if parsed == nil {
		os.Exit(1)
	}

	if inPlace {
		if err := os.WriteFile(file, []byte(parse(file).String()), 0644); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Println(parsed)
	}

	os.Exit(0)
}
