package main

import (
	"log"
	"os"
	"regexp"

	"github.com/alecthomas/repr"
)

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
		log.Println("NG:", filename)
		log.Fatal(err)
	} else {
		log.Println("OK:", filename)
	}

	return actual
}

func main() {
}
