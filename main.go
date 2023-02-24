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

	"github.com/apxxxxxxe/formatyaya/ast"
)

const version = "0.1.2"

var (
	repLF              = regexp.MustCompile(`(\r\n|\r)`)
	repTrailingSpace   = regexp.MustCompile(`[\t ]+\n`)
	repEdgeLF          = regexp.MustCompile(`(^\n+|\n+$)`)
	repDoubleBlankLine = regexp.MustCompile(`(?m)(^[\t ]*\n){2}`)
)

const (
	hasCRLF = iota
	hasCR
	hasLF
)

func identifyCRLF(s string) int {
	if crlf := strings.Contains(s, "\r\n"); crlf {
		return hasCRLF
	} else if cr := strings.Contains(s, "\r"); cr {
		return hasCR
	} else if lf := strings.Contains(s, "\n"); lf {
		return hasLF
	} else {
		return hasCRLF
	}
}

func parse(b []byte) *ast.Root {
	// 改行コードの統一
	src := repLF.ReplaceAllString(string(b), "\n")
	// 行末の空白文字を削除
	src = repTrailingSpace.ReplaceAllString(src, "\n")

	actual, err := ast.Parser.ParseString("", string(src))
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
		ast.Indent = strings.Repeat(" ", spaceCount)
	} else {
		ast.Indent = "	"
	}

	b, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	lfCode := identifyCRLF(string(b))

	parsed := parse(b)
	if parsed == nil {
		os.Exit(3)
	}

	parsedString := repEdgeLF.ReplaceAllString(parsed.String(), "")
	parsedString = repDoubleBlankLine.ReplaceAllString(parsedString, "\n")

	lfDict := map[int]string{
		hasCRLF: "\r\n",
		hasCR:   "\r",
		hasLF:   "\n",
	}
	parsedString = strings.ReplaceAll(parsedString, "\n", lfDict[lfCode])

	if inPlace {
		if err := os.WriteFile(file, []byte(parsedString), 0644); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Println(parsedString)
	}

	os.Exit(0)
}
