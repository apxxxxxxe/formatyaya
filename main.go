package main

import (
	"log"
	"strings"

	"github.com/alecthomas/repr"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Root struct {
	Functions []*Func `"\n"* (@@ "\n"*)*`
}

type Func struct {
	FunctionName string   `@Ident`
	LF           struct{} `"\n"?`
	Expr         []*Expr  `"{" "\n"* (@@ "\n"*)* "}"`
}

type Expr struct {
	Left  *Terminal `@@`
	Op    string    `( @Oper`
	Right *Terminal `  @@    )?`
}

type Terminal struct {
	String *String `( @@`
	Ident  string  `| @Ident)`

	Call *Call `("(" @@ ")")?`
}

type Call struct {
	Terminal []*Expr `(@@ ","?)*`
}

type String struct {
	Value string `"\"" @Char "\""`
}

var (
	def = lexer.MustStateful(lexer.Rules{
		"Root": { // Rootは特別な名前　これが基点のRule
			// TODO: defineとか
			{Name: "Ident", Pattern: `\w+`, Action: nil},
			{Name: `Function`, Pattern: `\{`, Action: lexer.Push("FuncRule")},
			lexer.Include("Common"),
		},
		"FuncRule": {
			{Name: "Ident", Pattern: `\w+`, Action: nil},
			{Name: "String", Pattern: `"`, Action: lexer.Push("StringRule")},
			{Name: `Oper`, Pattern: `(-|\+|/|\*|%|=|==|!=|>=|<=|>|<)`, Action: nil},
			{Name: "Call", Pattern: `\(`, Action: lexer.Push("CallRule")},
			{Name: "FuncEnd", Pattern: `\}`, Action: lexer.Pop()},
			lexer.Include("Common"),
		},
		"StringRule": {
			{Name: "Char", Pattern: `[^"]+`, Action: nil},
			{Name: "StringEnd", Pattern: `"`, Action: lexer.Pop()},
		},
		"CallRule": {
			lexer.Include("FuncRule"),
			{Name: "Delim", Pattern: `,`, Action: nil},
			{Name: "CallEnd", Pattern: `\)`, Action: lexer.Pop()},
		},
		"Common": {
			{Name: "LF", Pattern: `\n`, Action: nil},
			{Name: `Whitespace`, Pattern: `\s+`, Action: nil},
		},
	})
	parser = participle.MustBuild[Root](
		participle.Lexer(def),
		participle.Elide("Whitespace"),
	)
)

func main() {
	src := `
  Func
  {
    a = b(c,"d")
    e = f(g,func("h"))
  }
  `

	actual, err := parser.Parse("", strings.NewReader(src))
	repr.Println(actual)
	if err != nil {
		log.Fatal(err)
	}
}
