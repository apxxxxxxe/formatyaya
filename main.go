package main

import (
	"log"
	"os"

	"github.com/alecthomas/repr"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Root struct {
	Rootentities []*RootEntity `@@*`
}

type RootEntity struct {
	Functions      *Func    `( @@`
	CommentOneLine string   `| "//" @CommentOneLineChar ("\r\n"|"\r"|"\n")`
	BlankLine      struct{} `| LF)`
}

type Func struct {
	FunctionName string `@Ident`

	Begin      struct{}      `("\r\n"|"\r"|"\n")? "{" ("\r\n"|"\r"|"\n")*`
	FuncEntity []*FuncEntity `@@*`
	End        struct{}      `"}"`
}

type FuncEntity struct {
	Main *Line         `( @@ ("\r\n"|"\r"|"\n"|";")`
	Sub  []*FuncEntity `| "{" ("\r\n"|"\r"|"\n")* @@* ("\r\n"|"\r"|"\n")* "}" ("\r\n"|"\r"|"\n")*)`
}

// 関数内に1行で取りうる式
type Line struct {
	Separator      string   `( @Separator`
	CommentOneLine string   `| @CommentOneLine`
	Flow           *Flow    `| @@`
	Asign          *Asign   `| @@`
	Logic          *Logic   `| @@`
	Value          *Primary `| @@)`
}

// 代入式
type Asign struct {
	Left string `@Ident`

	OperAsign      string `( @OperAsign`
	Right          *Logic `  @@`
	OperAsignUnary string `| @OperAsignUnary)`
}

type Logic struct {
	Comparison *Comparison `@@`

	OperLogic string `[ @("||"|"&&")`
	Right     *Logic `  @@]`
}

type Comparison struct {
	Addition *Addition `@@`

	OperComp string      `[ @("=""="|"!""="|">""="|"<""="|">"|"<"|"_in_"|"!_in_")`
	Right    *Comparison `  @@]`
}

type Addition struct {
	Multipulation *Multipulation `@@`

	Op    string    `[ @("+"|"-")`
	Right *Addition `  @@]`
}

type Multipulation struct {
	Unary *Unary `@@`

	Op    string         `[ @("*"|"/"|"%")`
	Right *Multipulation `  @@]`
}

type Unary struct {
	Unary   string   `(@OperUnary)?`
	Primary *Primary `@@`
}

// 単項: 左辺と右辺の両方になりうる式
type Primary struct {
	SingleQuote string    `( "'" @(SingleQuoteChar)? "'"`
	DoubleQuote string    `| "\"" @(DoubleQuoteChar)? "\""`
	FuncCall    *FuncCall `| @@`
	Number      *Number   `| @@`
	Ident       string    `| @Ident)`
}

type Flow struct {
	Key     string `( @FlowKey`
	KeyExpr string `| @FlowKeyExpr`
	Expr    *Logic `  @@)`

	OneLineSub   *FuncEntity   `( ("\r\n"|"\r"|"\n") @@`
	MultiLineSub []*FuncEntity `| ("\r\n"|"\r"|"\n")? "{" ("\r\n"|"\r"|"\n")* @@* "}")`
}

type FuncCall struct {
	FuncName string     `@Ident`
	Args     []*Primary `"(" (@@ ","?)* ")"`
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
			lexer.Include("Comments"),
		},
		"FuncRule": {
			{Name: `OperAsign`, Pattern: `(=|:=|\+=|-=|\*=|/=|%=|\+:=|-:=|\*:=|/:=|%:=|,=)`, Action: nil},
			{Name: `OperAsignUnary`, Pattern: `[^\s](--|\+\+)`, Action: nil},
			{Name: `Separator`, Pattern: `--`, Action: nil},
			{Name: `OperUnary`, Pattern: `!`, Action: nil}, // TODO: -はない？要確認
			{Name: `ExprEnd`, Pattern: `;`, Action: nil},
			{Name: `OperCalc`, Pattern: `(-|\+|/|\*|%)`, Action: nil},
			{Name: `OperComp`, Pattern: `(==|!=|>=|<=|>|<|_in_|!_in_)`, Action: nil},
			{Name: `OperLogic`, Pattern: `(\|\||&&)`, Action: nil},
			{Name: "FlowKeyExpr", Pattern: `(while|if|elseif)`, Action: nil},
			{Name: "FlowKey", Pattern: `else`, Action: nil},
			{Name: `Function`, Pattern: `\{`, Action: lexer.Push("FuncRule")},
			{Name: "Number", Pattern: `\d+(\.\d+)?`, Action: nil},
			{Name: "SingleQuoteString", Pattern: `'`, Action: lexer.Push("SingleQuoteStringRule")},
			{Name: "DoubleQuoteString", Pattern: `"`, Action: lexer.Push("DoubleQuoteStringRule")},
			lexer.Include("Common"),
			lexer.Include("Comments"),
			{Name: "Call", Pattern: `\(`, Action: lexer.Push("CallRule")},
			{Name: "FuncEnd", Pattern: `\}`, Action: lexer.Pop()},
		},
		"SingleQuoteStringRule": {
			{Name: "SingleQuoteChar", Pattern: `[^']+`, Action: nil},
			{Name: "SingleQuoteStringEnd", Pattern: `'`, Action: lexer.Pop()},
		},
		"DoubleQuoteStringRule": {
			{Name: "DoubleQuoteChar", Pattern: `[^"]+`, Action: nil},
			{Name: "DoubleQuoteStringEnd", Pattern: `"`, Action: lexer.Pop()},
		},
		"CallRule": {
			lexer.Include("FuncRule"),
			{Name: "Delim", Pattern: `,`, Action: nil},
			{Name: "CallEnd", Pattern: `\)`, Action: lexer.Pop()},
		},
		"CommentOneLineRule": {
			{Name: "CommentOneLineEnd", Pattern: `(\r\n|\r|\n)`, Action: lexer.Pop()},
			{Name: `CommentOneLineChar`, Pattern: `[^\r\n]*`, Action: nil},
		},
		"CommentMultiLineRule": {
			{Name: "CommentMultiLineEnd", Pattern: `\*/`, Action: lexer.Pop()},
			{Name: `CommentMultiLineChar`, Pattern: `[^*/]*`, Action: nil}, // TODO:要修正
		},
		"Common": {
			{Name: "LF", Pattern: `(\r\n|\r|\n)`, Action: nil},
			{Name: `Whitespace`, Pattern: `\s+`, Action: nil},
			{Name: "Ident", Pattern: `\w[^ \r\n!"#$%&()*+,-/:;<=>?@\[\]{|}]*`, Action: nil}, //TODO:`も禁止文字に入れる
		},
		"Comments": {
			{Name: `CommentOneLine`, Pattern: `//`, Action: lexer.Push("CommentOneLineRule")},
			{Name: `CommentMultiLine`, Pattern: `/\*`, Action: lexer.Push("CommentMultiLineRule")},
		},
	})
	parser = participle.MustBuild[Root](
		participle.Lexer(def),
		participle.Elide("Whitespace"),
		participle.Elide("LF"),
	)
)

func main() {
	/*
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
		    a = b(c,"d") || vim == 1 * 2
		    p++
		  }

		  `
			actual, err := parser.Parse("", strings.NewReader(src))
	*/

	fp, err := os.Open("./yaya_bootend.dic")
	if err != nil {
		log.Fatal(err)
	}

	actual, err := parser.Parse("", fp)
	repr.Println(actual)
	if err != nil {
		log.Fatal(err)
	}
}
