package main

type Root struct {
	Rootentities []*RootEntity `@@*`
}

// Root内で出現しうる記述
type RootEntity struct {
	Definition       *Definition `( @@ ("\r\n"|"\r"|"\n")`
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
	Asign            *Asign        `| @@`
	Logic            *Logic        `| @@`
	Value            *Primary      `| @@`
	Sub              []*FuncEntity `| "{" @@* "}")`
}

// フロー制御文
type Flow struct {
	Key        string   `( @FlowKey`
	KeyFor     string   `| @"for"`
	ExprFor    *ExprFor `  @@`
	KeyPrimery string   `| @FlowKeyPrimary`
	Primery    *Primary `  @@`
	KeyExpr    string   `| @FlowKeyExpr`
	Expr       *Expr    `  @@)`

	OneLineSub   *FuncEntity   `(  @@`
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
	Enum *Enum `@@`
}

type Enum struct {
	Logic *Logic `@@`

	OperEnum string `( @","`
	Right    *Enum  `  @@)?`
}

// 論理演算式
type Logic struct {
	Comparison *Comparison `@@`

	OperLogic string `( @("||"|"&&")`
	Right     *Logic `  @@)?`
}

// 比較演算式
type Comparison struct {
	Addition *Addition `@@`

	OperComp string      `( @("=" "="|"!="|">="|"<="|">"|"<"|"_in_"|"!_in_")`
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
	SingleQuote string     `( "'" @(SingleQuoteChar)? "'"`
	DoubleQuote string     `| "\"" @(DoubleQuoteChar)? "\""`
	FuncCall    *FuncCall  `| @@`
	ArrayCall   *ArrayCall `| @@`
	Number      *Number    `| @@`
	Ident       string     `| @Ident`
	SubExpr     *Expr      `| "(" @@ ")")`
}

type FuncCall struct {
	FuncName string     `@Ident`
	Args     []*Primary `"(" (@@ ","?)* ")"`
}

type ArrayCall struct {
	ArrayName string     `@Ident`
	Args      []*Primary `("[" @@ "]")+`
}

type Number struct {
	Value float64 `@Number`
}
