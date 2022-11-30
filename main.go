package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/alecthomas/repr"
)

var (
	rep = regexp.MustCompile(`([^*])/(\r\n|\r|\n)[\t ]*`)
)

func main() {
	var actual *Root
	files := []string{"yaya_bootend.dic", "yaya_aitalk.dic"}

	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			panic(err)
		}

		src := rep.ReplaceAllString(string(b), "$1")
		if err := os.WriteFile("replaced_"+f, []byte(src), 0644); err != nil {
			panic(err)
		}

		actual, err = parser.ParseString("", string(src))
		repr.Println(actual)
		if err != nil {
			fmt.Println("NG:", f)
			log.Fatal(err)
		} else {
			fmt.Println("OK:", f)
		}
	}
}
