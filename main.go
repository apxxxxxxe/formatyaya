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
	rep   = regexp.MustCompile(`([^*])/(\r\n|\r|\n)[\t ]*`)
	repLF = regexp.MustCompile(`(\r\n|\r|\n)`)
)

func parse(filename string) *Root {
	b, err := os.ReadFile(filepath.Join(dirname, filename))
	if err != nil {
		panic(err)
	}

	src := rep.ReplaceAllString(string(b), "$1")
	src = repLF.ReplaceAllString(src, "\n")

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

func format(value interface{}, depth int, parentName string) string {
	const indent = "	"
	result := ""

	parentv := reflect.Indirect(reflect.ValueOf(value))
	parentt := parentv.Type()
	switch parentt.Kind() {
	case reflect.Ptr:
		// ポインタなら再帰呼び出しで実体を処理する
		result += format(parentv.Interface(), depth, "")

	case reflect.Slice:
		// スライスなら各要素を再帰処理
		for i := 0; i < parentv.Len(); i++ {
			e := parentv.Index(i)
			switch parentName {
			case "FuncEntities", "Sub", "FlowOneLineSub", "FlowMultiLineSub":
				result += format(e.Addr().Interface(), depth+1, "")
			case "FlowConst":
				result += format(e.Addr().Interface(), depth, "")
				if i < parentv.Len()-1 {
					result += ","
				}
			default:
				result += format(e.Addr().Interface(), depth, "")
			}
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
				s = format(fv.Interface(), depth, ft.Name)
			} else if ft.Type.Kind() == reflect.Ptr && ft.Type.Elem().Kind() == reflect.Struct {
				// 構造体のポインタ
				s = format(fv.Interface(), depth, "")
			} else {
				// 任意の型
				if ft.Type.Kind() == reflect.Int {
					fmt.Println("float:", fv)
				} else {
					s = fmt.Sprint(fv)
				}
			}

			sl := ""
			if s == "" {
			} else {
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
					sl += ", "
				case "OperOr", "OperAnd", "OperAdd", "OperCalc", "OperComp", "OperAsign", "OperMulti":
					sl += " " + s + " "
				case "FlowKey", "FlowKeyFor", "FlowKeyForEach", "FlowKeyConst", "FlowKeyExpr":
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
				case "DefSpaceBefore":
					sl += s + " "
				case "DefTabBefore":
					sl += s + "	"
				default:
					sl += s
				}
			}
			result += sl
		}
	}

	return result
}

func main() {
	finfo, err := ioutil.ReadDir(dirname)
	if err != nil {
		panic(err)
	}

	for i, f := range finfo {
		if i >= 3 {
			break
		}

		if strings.HasPrefix(f.Name(), "replaced_") {
			continue
		}

		fmt.Println(f.Name())

		actual := parse(f.Name())
		// repr.Println(actual)

		if err := os.WriteFile("out.txt", []byte(format(actual, 0, "")), 0644); err != nil {
			panic(err)
		}
	}
}
