package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/repr"
)

func main() {
	var fp *os.File
	var err error
	var actual *Root
	files := []string{"./yaya_bootend.dic", "./yaya_aitalk.dic"}

	for _, f := range files {
		fp, err = os.Open(f)
		if err != nil {
			log.Fatal(err)
		}
		defer fp.Close()

		actual, err = parser.Parse("", fp)
		repr.Println(actual)
		if err != nil {
			fmt.Println("NG:", fp.Name())
			log.Fatal(err)
		} else {
			fmt.Println("OK:", fp.Name())
		}
	}
}
