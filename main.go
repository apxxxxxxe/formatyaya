package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

const dirname = "files"

var (
	rep        = regexp.MustCompile(`([^*])/(\r\n|\r|\n)[\t ]*`)
	repLF      = regexp.MustCompile(`(\r\n|\r|\n)`)
	repFor     = regexp.MustCompile(`for([^\n{;]*);([^\n{;]*);`)
	repForEach = regexp.MustCompile(`foreach([^;]*);`)

	// ;を置換する際の一時的な文字 ファイル内に存在するどの文字とも被ってはいけない
	tmpColon = string(rune(0x12))
)

func parse(filename string) *Root {
	b, err := os.ReadFile(filepath.Join(dirname, filename))
	if err != nil {
		panic(err)
	}

	src := rep.ReplaceAllString(string(b), "$1")
	src = repLF.ReplaceAllString(src, "\n")

	src = repFor.ReplaceAllString(src, "for$1"+tmpColon+"$2"+tmpColon)
	src = repForEach.ReplaceAllString(src, "foreach$1"+tmpColon)

	src = strings.ReplaceAll(src, ";", "\n")
	src = strings.ReplaceAll(src, tmpColon, ";")

	if err := os.WriteFile(filepath.Join(dirname, "replaced_"+filename), []byte(src), 0644); err != nil {
		panic(err)
	}

	actual, err := parser.ParseString("", string(src))
	if err != nil {
		fmt.Println("NG:", filename)
		log.Fatal(err)
	} else {
		fmt.Println("OK:", filename)
	}

	return actual
}

func format(value interface{}, depth int) string {
	const indent = "	"
	result := ""

	parentv := reflect.Indirect(reflect.ValueOf(value))
	parentt := parentv.Type()
	switch parentt.Kind() {
	case reflect.Ptr:
		// ポインタなら再帰呼び出しで実体を処理する
		result += format(parentv.Interface(), depth)

	case reflect.Slice:
		// スライスなら各要素を再帰処理
		for i := 0; i < parentv.Len(); i++ {
			e := parentv.Index(i)
			d := depth
			if e.Type().String() == "*main.FuncEntity" {
				d++
			}
			result += format(e.Addr().Interface(), d)
		}

	case reflect.Struct:
		for i := 0; i < parentt.NumField(); i++ {
			ft := parentt.Field(i)
			fv := parentv.FieldByName(ft.Name)

			s := ""
			if ft.Type.Kind() == reflect.Ptr && fv.IsNil() {
				// nilポインタなら処理しない
				continue
			} else if ft.Type.Kind() == reflect.Slice {
				// スライス
				s = format(fv.Interface(), depth)
			} else if ft.Type.Kind() == reflect.Ptr && ft.Type.Elem().Kind() == reflect.Struct {
				// 構造体のポインタ
				s = format(fv.Interface(), depth)
			} else {
				// 任意の型
				if ft.Type.Kind() == reflect.Int {
					fmt.Println("float:", fv)
				} else {
					s = fmt.Sprint(fv)
				}
			}

			if s != "" {
				sl := ""
				switch parentt.Name() {
				case "FuncEntity":
					//多重にインデントを付けるのを防ぐ 応急処置的？アルゴリズムに不備があるかも
					if !strings.HasPrefix(s, indent) {
						sl += strings.Repeat(indent, depth)
					}
				}

				switch ft.Name {
				case "RootEntity":
					sl += s + "\n"
				case "FuncArgs":
					sl += "(" + s + ")"
				case "ArrayArgs":
					sl += "[" + s + "]"
				case "ForInitAsign", "ForEndExpr":
					sl += s + "; "
				case "OperEnum":
					sl += s + " "
				case "OperOr", "OperAnd", "OperAdd", "OperCalc", "OperComp", "OperAsign", "OperMulti":
					sl += " " + s + " "
				case "KeyFlow", "KeyFor", "KeyForEach", "KeyConst", "KeyExpr":
					sl += s + " "
				case "Definition":
					sl += s + "\n"
				case "FuncEntitiesBegin", "FlowMultiLineSubBegin":
					sl += " " + s + "\n"
				case "SubBegin":
					sl += s + "\n"
				case "FlowOneLineSub":
					sl += "\n" + strings.Repeat(indent, depth) + s
				case "FlowMultiLineSubEnd":
					sl += strings.Repeat(indent, depth) + s + "\n"
				case "Separator", "CommentMultiLine", "Value", "SubEnd":
					sl += s + "\n"
				case "SingleQuote":
					sl += "'" + s + "'"
				case "DoubleQuote":
					sl += "\"" + s + "\""
				case "HearDocumentSingle":
					sl += "<<'" + s + "'>>"
				case "HearDocumentDouble":
					sl += "<<\"" + s + "\">>"
				default:
					sl += s
				}

				result += sl
			}
		}
	}

	return result
}

func main() {

	finfo, err := ioutil.ReadDir(dirname)
	if err != nil {
		panic(err)
	}

	for _, f := range finfo {
		if strings.HasPrefix(f.Name(), "replaced_") {
			continue
		}

		fmt.Println(f.Name(), "\n")

		actual := parse(f.Name())
		// repr.Println(actual)

		if err := os.WriteFile("out.txt", []byte(format(actual, 0)), 0644); err != nil {
			panic(err)
		}

		break
	}
}
