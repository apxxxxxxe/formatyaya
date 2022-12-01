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
	Key    string `@("#define"|"#globaldefine")`
	Before string `@DefinitionChar`
	After  string `@DefinitionChar`
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
	Key        string   `( @FlowKey`
	KeyFor     string   `| ( @"for"`
	ExprFor    *ExprFor `    @@)`
	KeyPrimery string   `| ( @FlowKeyPrimary`
	Primery    *Primary `    @@)`
	KeyExpr    string   `| ( @FlowKeyExpr`
	Expr       *Expr    `    @@))`

	ExprEnd      string        `( @("\n"|";")`
	OneLineSub   *FuncEntity   `  @@`
	MultiLineSub []*FuncEntity `| "{" @@* "}")`
}

type ExprFor struct {
	InitAsign *Asign `@@ ";"`
	EndExpr   *Expr  `@@ ";"`
	LoopAsign *Asign `@@`
}

// 代入式
type Asign struct {
	Left string `@Ident`

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
	String    *String    `( @@`
	FuncCall  *FuncCall  `| @@`
	ArrayCall *ArrayCall `| @@`
	Number    *Number    `| @@`
	Ident     string     `| @Ident`
	SubExpr   *Expr      `| "(" @@ ")")`
}

type String struct {
	SingleQuote *SingleQuote `( @@`
	DoubleQuote *DoubleQuote `| @@)`
}

type SingleQuote struct {
	Lines []*SingleQuoteLine `"'" @@* "'"`
}

type SingleQuoteLine struct {
	Line string `@(SingleQuoteChar ContinuousLF?)`
}

type DoubleQuote struct {
	Lines []*DoubleQuoteLine `"\"" @@* "\""`
}

type DoubleQuoteLine struct {
	Line string `@(DoubleQuoteChar ContinuousLF?)`
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
