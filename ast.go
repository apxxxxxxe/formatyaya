package main

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type Root struct {
	Rootentities []*RootEntity `@@*`

	Tokens []lexer.Token
}

// Root内で出現しうる記述
type RootEntity struct {
	Definition *Definition `( @@ `
	Function   *Func       `| @@)`

	Tokens []lexer.Token
}

type Comment struct {
	CommentOneLine   string `( @CommentOneLine`
	CommentMultiLine string `| @CommentMultiLine)`

	Tokens []lexer.Token
}

type Definition struct {
	DefinitionSpace *DefinitionSpace `( @@`
	DefinitionTab   *DefinitionTab   `| @@)`

	Tokens []lexer.Token
}

type DefinitionSpace struct {
	DefSpaceKey   string `@DefinitionSpace`
	DefSpaceValue string `@DefinitionSpaceChar`

	Tokens []lexer.Token
}

type DefinitionTab struct {
	DefTabKey   string `@DefinitionTab`
	DefTabValue string `@DefinitionTabChar`

	Tokens []lexer.Token
}

type Func struct {
	FunctionName string `@FuncName`
	FunctionType string `(":" @("array"|"void"|"nonoverlap"|"sequential"))?` // @FuncTypeだと認識しないので即値で指定

	FuncEntitiesBegin string        `@"{"`
	FuncEntities      []*FuncEntity `@@*`
	FuncEntitiesEnd   string        `@"}"?`

	Tokens []lexer.Token
}

// Func内で出現しうる記述: 関数内に1行で取りうる式
type FuncEntity struct {
	OutputFixer string        `( @OutputFixer`
	Flow        *Flow         `| @@`
	PreValue    string        `| @PreValue?`
	Value       *Expr         `  @@`
	SubBegin    string        `| @"{"`
	Sub         []*FuncEntity `  @@*`
	SubEnd      string        `  @"}"?)`

	Tokens []lexer.Token
}

// フロー制御文
type Flow struct {
	FlowKey         string       `( @FlowKey`
	FlowKeyForEach  string       `| ( @FlowKeyForEach`
	FlowExprForEach *ExprForEach `    @@)`
	FlowKeyFor      string       `| ( @FlowKeyFor`
	FlowExprFor     *ExprFor     `    @@)`
	FlowKeyConst    string       `| ( @FlowKeyConst`
	FlowConst       []*Const     `    (@@ ","?)+)`
	FlowKeyExpr     string       `| ( @FlowKeyExpr`
	FlowExpr        *Expr        `    @@))`

	FlowOneLineSub        *FuncEntity   `( @@`
	FlowMultiLineSubBegin string        `| @"{"`
	FlowMultiLineSub      []*FuncEntity `  @@*`
	FlowMultiLineSubEnd   string        `  @"}")`

	Tokens []lexer.Token
}

type ExprFor struct {
	ForInitAsign *Asign `@@ ";"`
	ForEndExpr   *Expr  `@@ ";"`
	ForLoopAsign *Asign `@@`

	Tokens []lexer.Token
}

type ExprForEach struct {
	Array *Enum  `@@ ";"`
	Ident string `@Ident`

	Tokens []lexer.Token
}

// 条件式
type Expr struct {
	Enum *Enum `@@`

	Tokens []lexer.Token
}

type Enum struct {
	Asign *Asign `@@`

	OperEnum string `( @","`
	Right    *Enum  `  @@)?`

	Tokens []lexer.Token
}

// 代入式
type Asign struct {
	Or *Or `@@`

	OperAsign string `( @OperAsign`
	Right     *Asign `  @@)?`

	Tokens []lexer.Token
}

// 論理演算式Or
type Or struct {
	And *And `@@`

	OperOr string `( @"||"`
	Right  *Or    `  @@)?`

	Tokens []lexer.Token
}

// 論理演算式And
type And struct {
	Comparison *Comparison `@@`

	OperAnd string `( @"&&"`
	Right   *And   `  @@)?`

	Tokens []lexer.Token
}

// 比較演算式
type Comparison struct {
	Addition *Addition `@@`

	OperComp string      `( @("=" "="|"!""="|">="|"<="|">"|"<"|"_in_"|"!""_in_")`
	Right    *Comparison `  @@)?`

	Tokens []lexer.Token
}

// 加減法式
type Addition struct {
	Multipulation *Multipulation `@@`

	OperAdd string    `( @("+"|"-")`
	Right   *Addition `  @@)?`

	Tokens []lexer.Token
}

// 乗除法式
type Multipulation struct {
	Unary *Unary `@@`

	OperMulti string         `( @("*"|"/"|"%")`
	Right     *Multipulation `  @@)?`

	Tokens []lexer.Token
}

// 単項演算式
type Unary struct {
	Unary       string   `(@OperUnary)?`
	Primary     *Primary `@@`
	OperCalcOne string   `@("+" "+"|"--")?`

	Tokens []lexer.Token
}

// 単項: 左辺と右辺の両方になりうる式
type Primary struct {
	Const     *Const  `( @@`
	ArrayArgs []*Expr `  ("[" @@ "]")*`
	SubExpr   *Expr   `| "(" @@ ")")`

	Tokens []lexer.Token
}

type Const struct {
	String   *String   `( @@`
	FuncCall *FuncCall `| @@`
	Number   *Number   `| @@`
	Ident    string    `| @Ident)`

	Tokens []lexer.Token
}

type String struct {
	SingleQuote        string `( @"'" @SingleQuoteChar? @"'"`
	DoubleQuote        string `| @"\"" @DoubleQuoteChar? @"\""`
	HearDocumentSingle string `| @HearDocumentSingle`
	HearDocumentDouble string `| @HearDocumentDouble)`

	Tokens []lexer.Token
}

type FuncCall struct {
	FuncName string  `@Ident`
	FuncArgs []*Expr `"(" (@@ ","?)* ")"`

	Tokens []lexer.Token
}

// 目的はパースなのでstringでとっていい
type Number struct {
	Hex   string `( @HexNum`
	Bin   string `| @BinNum`
	Float string `| @Float`
	Int   string `| @Int)`

	Tokens []lexer.Token
}
