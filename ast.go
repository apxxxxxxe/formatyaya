package main

type Root struct {
	Rootentities []*RootEntity `@@*`
}

// Root内で出現しうる記述
type RootEntity struct {
	Definition *Definition `( @@ `
	DefLF      string      `  "\n"`
	Comment    *Comment    `| @@`
	Function   *Func       `| @@`
	LF         string      `| @"\n")`
}

type Comment struct {
	CommentOneLine   string `( @CommentOneLine`
	CommentMultiLine string `| @CommentMultiLine)`
}

type Definition struct {
	DefinitionSpace *DefinitionSpace `( @@`
	DefinitionTab   *DefinitionTab   `| @@)`
}

type DefinitionSpace struct {
	DefSpaceKey    string `@DefinitionSpace`
	DefSpaceBefore string `@DefinitionSpaceChar Space`
	DefSpaceAfter  string `@DefinitionSpaceChar`
}

type DefinitionTab struct {
	DefTabKey    string `@DefinitionTab`
	DefTabBefore string `@DefinitionTabChar TabSpace`
	DefTabAfter  string `@DefinitionTabChar`
}

type Func struct {
	FunctionName string `@FuncName`
	FunctionType string `(":" @("array"|"void"))?`

	LF string `| "\n"?`

	FuncEntitiesBegin string        `@"{" "\n"?`
	FuncEntities      []*FuncEntity `@@*`
	FuncEntitiesEnd   string        `@"}"?`
}

// Func内で出現しうる記述: 関数内に1行で取りうる式
type FuncEntity struct {
	Separator         string        `( @Separator "\n"`
	Comment           *Comment      `| @@`
	Flow              *Flow         `| @@`
	PreValue          string        `| @PreValue?`
	Value             *Expr         `  @@`
	CommentAfterValue *Comment      `  @@? ("}"|"\n"|";")`
	BlankLine         string        `| @BlankLine`
	SubBegin          string        `| @"{" "\n"?`
	Sub               []*FuncEntity `  @@*`
	SubEnd            string        `  @"}"?)`
}

// フロー制御文
type Flow struct {
	FlowKey         string       `( @FlowKey`
	FlowKeyForEach  string       `| ( @"for""each"`
	FlowExprForEach *ExprForEach `    @@)`
	FlowKeyFor      string       `| ( @"for"`
	FlowExprFor     *ExprFor     `    @@)`
	FlowKeyConst    string       `| ( @FlowKeyConst`
	FlowConst       []*Const     `    (@@ ","?)+)`
	FlowKeyExpr     string       `| ( @FlowKeyExpr`
	FlowExpr        *Expr        `    @@))`

	FlowExprEnd           string        `( ("\n"|";")`
	FlowComment           *Comment      `  @@*`
	FlowOneLineSub        *FuncEntity   `  @@`
	FlowOneLineSubEnd     string        `  @("\n"|";")?`
	FlowMultiLineSubBegin string        `| @"{" "\n"?`
	FlowMultiLineSub      []*FuncEntity `  @@*`
	FlowMultiLineSubEnd   string        `  @"}"?)`
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
	SingleQuote        string `( @"'" @SingleQuoteChar? @"'"`
	DoubleQuote        string `| @"\"" @DoubleQuoteChar? @"\""`
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
	Bin   string `| @BinNum`
	Float string `| @Float`
	Int   string `| @Int)`
}
