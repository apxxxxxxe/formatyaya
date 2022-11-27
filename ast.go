package main

type Root struct {
	Rootentities []*RootEntity `@@*`
}

// Root内で出現しうる記述
type RootEntity struct {
	Definition       *Definition `( @@ ("\r\n"|"\r"|"\n")`
	Functions        *Func       `| @@`
	CommentOneLine   string      `| "//" @CommentOneLineChar ("\r\n"|"\r"|"\n")`
	CommentMultiLine string      `| "/" "*" @CommentMultiLineChar "*" "/"`
	BlankLine        struct{}    `| LF)`
}

type Definition struct {
	Key    string `@("#define"|"#globaldefine") Space+`
	Before string `@DefinitionChar Space+`
	After  string `@DefinitionChar`
}

type Func struct {
	FunctionName string `@Ident`
	FunctionType string `(":" @FuncType)?`

	Begin      struct{}      `("\r\n"|"\r"|"\n")? "{" ("\r\n"|"\r"|"\n")*`
	FuncEntity []*FuncEntity `@@*`
	End        struct{}      `"}"`
}

// Func内で出現しうる記述
type FuncEntity struct {
	Main *Line         `( @@ ("\r\n"|"\r"|"\n"|";")`
	Sub  []*FuncEntity `| "{" ("\r\n"|"\r"|"\n")* @@* ("\r\n"|"\r"|"\n")* "}" ("\r\n"|"\r"|"\n")*)`
}

// 関数内に1行で取りうる式
type Line struct {
	Separator        string   `( @Separator`
	CommentOneLine   string   `| "//" @CommentOneLineChar ("\r\n"|"\r"|"\n")`
	CommentMultiLine string   `| "/" "*" @CommentMultiLineChar "*" "/"`
	Flow             *Flow    `| @@`
	Asign            *Asign   `| @@`
	Logic            *Logic   `| @@`
	Value            *Primary `| @@)`
}

// フロー制御文
type Flow struct {
	Key     string `( @FlowKey`
	KeyExpr string `| @FlowKeyExpr`
	Expr    *Logic `  @@)`

	OneLineSub   *FuncEntity   `( ("\r\n"|"\r"|"\n") @@`
	MultiLineSub []*FuncEntity `| ("\r\n"|"\r"|"\n")? "{" ("\r\n"|"\r"|"\n")* @@* "}")`
}

// 代入式
type Asign struct {
	Left string `@Ident`

	OperAsign      string `( @OperAsign`
	Right          *Logic `  @@`
	OperAsignUnary string `| @OperAsignUnary)`
}

// 論理演算式
type Logic struct {
	Comparison *Comparison `@@`

	OperLogic string `[ @("||"|"&&")`
	Right     *Logic `  @@]`
}

// 比較演算式
type Comparison struct {
	Addition *Addition `@@`

	OperComp string      `[ @("=" "="|"!" "="|">" "="|"<" "="|">"|"<"|"_in_"|"!_in_")`
	Right    *Comparison `  @@]`
}

// 加減法式
type Addition struct {
	Multipulation *Multipulation `@@`

	Op    string    `[ @("+"|"-")`
	Right *Addition `  @@]`
}

// 乗除法式
type Multipulation struct {
	Unary *Unary `@@`

	Op    string         `[ @("*"|"/"|"%")`
	Right *Multipulation `  @@]`
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
	Ident       string     `| @Ident)`
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
