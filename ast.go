package main

type Root struct {
	Rootentities []*RootEntity `@@*`
}

// Root内で出現しうる記述
type RootEntity struct {
	Definition       *Definition `( @@ "\n"`
	Functions        *Func       `| @@`
	CommentOneLine   string      `| @CommentOneLine`
	CommentMultiLine string      `| @CommentMultiLine)`
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

	FuncEntity []*FuncEntity `"{" @@* "}"`
}

// Func内で出現しうる記述: 関数内に1行で取りうる式
type FuncEntity struct {
	Separator        string        `( @Separator`
	CommentOneLine   string        `| @CommentOneLine`
	CommentMultiLine string        `| @CommentMultiLine`
	Flow             *Flow         `| @@`
	Value            *Expr         `| @@`
	Sub              []*FuncEntity `| "{" @@* "}")`
}

// フロー制御文
type Flow struct {
	Key         string       `( @FlowKey`
	KeyForEach  string       `| ( @"for""each"`
	ExprForEach *ExprForEach `    @@)`
	KeyFor      string       `| ( @"for"`
	ExprFor     *ExprFor     `    @@)`
	KeyConst    string       `| ( @FlowKeyConst`
	Const       []*Const     `    (@@ ","?)+)`
	KeyExpr     string       `| ( @FlowKeyExpr`
	Expr        *Expr        `    @@))`

	ExprEnd      string        `( @("\n"|";")`
	OneLineSub   *FuncEntity   `  @@`
	MultiLineSub []*FuncEntity `| "{" @@* "}")`
}

type ExprFor struct {
	InitAsign *Asign `@@ ";"`
	EndExpr   *Expr  `@@ ";"`
	LoopAsign *Asign `@@`
}

type ExprForEach struct {
	Array *Enum  `@@ ";"`
	Ident string `@Ident`
}

// 代入式
type Asign struct {
	LeftIdent string `@Ident`

	OperAsign      string `( @OperAsign`
	Right          *Expr  `  @@`
	OperAsignUnary string `| @("+" "+"|"-" "-"))`
}

// 条件式
type Expr struct {
	Asign *Asign `( @@`
	Enum  *Enum  `| @@)`
}

type Enum struct {
	Logic *Or `@@`

	OperEnum string `( @","`
	Right    *Enum  `  @@)?`
}

// 論理演算式And
type Or struct {
	And *And `@@`

	OperOr string `( @"||"`
	Right  *Or    `  @@)?`
}

// 論理演算式Or
type And struct {
	Comparison *Comparison `@@`

	OperAnd string `( @"&&"`
	Right   *And   `  @@)?`
}

// 比較演算式
type Comparison struct {
	Addition *Addition `@@`

	OperComp string      `( @("=" "="|"!""="|">="|"<="|">"|"<"|"_in_"|"!_in_")`
	Right    *Comparison `  @@)?`
}

// 加減法式
type Addition struct {
	Multipulation *Multipulation `@@`

	Op    string    `( @("+"|"-")`
	Right *Addition `  @@)?`
}

// 乗除法式
type Multipulation struct {
	Unary *Unary `@@`

	Op    string         `( @("*"|"/"|"%")`
	Right *Multipulation `  @@)?`
}

// 単項演算式
type Unary struct {
	Unary   string   `(@OperUnary)?`
	Primary *Primary `@@`
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
	Args     []*Expr `"(" (@@ ","?)* ")"`
}

type ArrayCall struct {
	ArrayName string  `@Ident`
	Args      []*Expr `("[" @@ "]")+`
}

type Number struct {
	Number float64 `( @Number`
	Hex    int     `| @HexNum)`
}
