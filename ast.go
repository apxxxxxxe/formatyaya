package main

type Root struct {
	Rootentities []*RootEntity `@@*`
}

// Root内で出現しうる記述
type RootEntity struct {
	Definition       *Definition `( @@ `
	DefLF            string      `  "\n"`
	Function         *Func       `| @@`
	CommentOneLine   string      `| @CommentOneLine`
	CommentMultiLine string      `| @CommentMultiLine`
	LF               string      `| @"\n")`
}

type Definition struct {
	DefinitionSpace *DefinitionSpace `( @@`
	DefinitionTab   *DefinitionTab   `| @@)`
}

type DefinitionSpace struct {
	Key    string `@DefinitionSpace`
	Before string `@DefinitionSpaceChar Space`
	After  string `@DefinitionSpaceChar`
}

type DefinitionTab struct {
	Key    string `@DefinitionTab`
	Before string `@DefinitionTabChar TabSpace`
	After  string `@DefinitionTabChar`
}

type Func struct {
	FunctionName string `@FuncName`
	FunctionType string `(":" @("array"|"void"))?`

	LF string `| "\n"?`

	FuncEntitiesBegin string        `@"{"`
	FuncEntities      []*FuncEntity `@@*`
	FuncEntitiesEnd   string        `@"}"`
}

// Func内で出現しうる記述: 関数内に1行で取りうる式
type FuncEntity struct {
	Separator        string        `( @Separator`
	CommentOneLine   string        `| @CommentOneLine`
	CommentMultiLine string        `| @CommentMultiLine`
	Flow             *Flow         `| @@`
	Value            *Expr         `| @@`
	LF               string        `| "\n"`
	SubBegin         string        `| @"{"`
	Sub              []*FuncEntity `  @@*`
	SubEnd           string        `  @"}")`
}

// フロー制御文
type Flow struct {
	KeyFlow     string       `( @FlowKey`
	KeyForEach  string       `| ( @"for""each"`
	ExprForEach *ExprForEach `    @@)`
	KeyFor      string       `| ( @"for"`
	ExprFor     *ExprFor     `    @@)`
	KeyConst    string       `| ( @FlowKeyConst`
	Const       []*Const     `    (@@ ","?)+)`
	KeyExpr     string       `| ( @FlowKeyExpr`
	Expr        *Expr        `    @@))`

	ExprEnd               string        `( ("\n"|";")`
	FlowOneLineSub        *FuncEntity   `  @@`
	FlowMultiLineSubBegin string        `| @"{"`
	FlowMultiLineSub      []*FuncEntity `  @@*`
	FlowMultiLineSubEnd   string        `  @"}")`
}

type ExprFor struct {
	ForInitAsign *Asign `@@ ";"`
	ForEndExpr   *Expr  `@@ ";"`
	ForLoopAsign *Asign `@@`
}

type ExprForEach struct {
	Array *Enum  `@@ ";"`
	Ident string `@Ident`
}

// 条件式
type Expr struct {
	Enum *Enum `@@`
}

type Enum struct {
	Asign *Asign `@@`

	OperEnum string `( @","`
	Right    *Enum  `  @@)?`
}

// 代入式
type Asign struct {
	Or *Or `@@`

	OperAsign string `( @OperAsign`
	Right     *Asign `  @@)?`
}

// 論理演算式Or
type Or struct {
	And *And `@@`

	OperOr string `( @"||"`
	Right  *Or    `  @@)?`
}

// 論理演算式And
type And struct {
	Comparison *Comparison `@@`

	OperAnd string `( @"&&"`
	Right   *And   `  @@)?`
}

// 比較演算式
type Comparison struct {
	Addition *Addition `@@`

	OperComp string      `( @("=" "="|"!""="|">="|"<="|">"|"<"|"_in_"|"!""_in_")`
	Right    *Comparison `  @@)?`
}

// 加減法式
type Addition struct {
	Multipulation *Multipulation `@@`

	OperAdd string    `( @("+"|"-")`
	Right   *Addition `  @@)?`
}

// 乗除法式
type Multipulation struct {
	Unary *Unary `@@`

	OperMulti string         `( @("*"|"/"|"%")`
	Right     *Multipulation `  @@)?`
}

// 単項演算式
type Unary struct {
	Unary       string   `(@OperUnary)?`
	Primary     *Primary `@@`
	OperCalcOne string   `@("+" "+"|"--")?`
}

// 単項: 左辺と右辺の両方になりうる式
type Primary struct {
	Const   *Const `( @@`
	SubExpr *Expr  `| "(" @@ ")")`
}

type Const struct {
	String    *String    `( @@`
	FuncCall  *FuncCall  `| @@`
	ArrayCall *ArrayCall `| @@`
	Number    *Number    `| @@`
	Ident     string     `| @Ident)`
}

type String struct {
	SingleQuote        string `( "'" @SingleQuoteChar? "'"`
	DoubleQuote        string `| "\"" @DoubleQuoteChar? "\""`
	HearDocumentSingle string `| @HearDocumentSingle`
	HearDocumentDouble string `| @HearDocumentDouble)`
}

type FuncCall struct {
	FuncName string  `@Ident`
	FuncArgs []*Expr `"(" (@@ ","?)* ")"`
}

type ArrayCall struct {
	ArrayName string  `@Ident`
	ArrayArgs []*Expr `("[" @@ "]")+`
}

// 目的はパースなのでstringでとっていい
type Number struct {
	Hex   string `( @HexNum`
	Float string `| @Float`
	Int   string `| @Int)`
}
