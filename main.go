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
	Begin        struct{} `"\n"? "{" "\n"*`

	FuncEntity []*FuncEntity `(@@ "\n"*)*`

	End struct{} `"}"`
}

type FuncEntity struct {
	Main *Expr `( @@`
	Sub  *Expr `| "{" "\n"* (@@ "\n"*)* "}")`
}

// 関数内に1行で取りうる式
type Expr struct {
	Asign     *Asign `( @@ `
	Separator string `| @Separator)`
}

// 代入式
type Asign struct {
	Left *Primary `@@`

	OperAsign      string      `( @OperAsign`
	Right          *Comparison `  @@        `
	OperAsignUnary string      `| @OperAsignUnary)?`
}

type MuitiPulation struct {
	Left *Addition `@@`

	Op    string         `( @("*"|"/"|"%")`
	Right *MuitiPulation `  @@)*`
}

type Addition struct {
	Left *Comparison `@@`

	Op    string    `( @("+"|"-")`
	Right *Addition `  @@)*`
}

type Comparison struct {
	Left *Logic `@@`

	OperComp string      `( @OperComp`
	Right    *Comparison `  @@)*`
}

type Logic struct {
	Left *Primary `@@`

	AND struct{} `(( "&&"`
	OR  struct{} ` | "||")`

	Right *Logic `@@)*`
}

// 単項: 左辺と右辺の両方になりうる式
type Primary struct {
	String   *String   `( @@`
	FuncCall *FuncCall `| @@`
	Number   *Number   `| @@`
	Ident    string    `| @Ident)`
}

type Flow struct {
	Key string `@Flow`
}

type FuncCall struct {
	FuncName string     `@Ident`
	Args     []*Primary `"(" (@@ ","?)* ")"`
}

type String struct {
	Value string `"\"" @Char "\""`
}

type Number struct {
	Value float64 `@Number`
}

var (
	def = lexer.MustStateful(lexer.Rules{
		"Root": { // Rootは特別な名前　これが基点のRule
			// TODO: defineとか
			{Name: `Function`, Pattern: `\{`, Action: lexer.Push("FuncRule")},
			lexer.Include("Common"),
		},
		"FuncRule": {
			{Name: `OperAsignUnary`, Pattern: `(--|\+\+)`, Action: nil},
			{Name: `OperUnary`, Pattern: `!`, Action: nil}, // TODO: -はない？要確認
			{Name: `OperCalc`, Pattern: `(-|\+|/|\*|%)`, Action: nil},
			{Name: `OperLogic`, Pattern: `(\|\||&&)`, Action: nil},
			{Name: `OperComp`, Pattern: `(==|!=|>=|<=|>|<|_in_|!_in_)`, Action: nil},
			{Name: `OperAsign`, Pattern: `(=|:=|\+=|-=|\*=|/=|%=|\+:=|-:=|\*:=|/:=|%:=|,=)`, Action: nil},
			{Name: "FlowExpr", Pattern: `(while|if|elseif)`, Action: nil},
			{Name: `Function`, Pattern: `\{`, Action: lexer.Push("FuncRule")},
			{Name: `Separator`, Pattern: `--`, Action: nil},
			{Name: "Number", Pattern: `\d+(\.\d+)?`, Action: nil},
			{Name: "String", Pattern: `"`, Action: lexer.Push("StringRule")},
			lexer.Include("Common"),
			{Name: "Call", Pattern: `\(`, Action: lexer.Push("CallRule")},
			{Name: "FuncEnd", Pattern: `\}`, Action: lexer.Pop()},
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
			{Name: "Ident", Pattern: `\w[^ \n!"#$%&()*+,-/:;<=>?@\[\]{|}]*`, Action: nil}, //TODO:`も禁止文字に入れる
		},
	})
	parser = participle.MustBuild[Root](
		participle.Lexer(def),
		participle.Elide("Whitespace"),
	)
)

func main() {
	src := `
  FuncHoge
  {
    {
      e = f(g,func("h"))
    }
    a = b(c,"d")
  }

  FuncFuga
  {
    {
      e = f(g,func("h"))
    }
    a = b(c,"d") || vim == 1
    p++
  }

  `

	actual, err := parser.Parse("", strings.NewReader(src))
	repr.Println(actual)
	if err != nil {
		log.Fatal(err)
	}
}
