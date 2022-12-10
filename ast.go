package main

import (
	"regexp"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
)

var (
	Indent       = "  "
	repBlankLine = regexp.MustCompile(`(?m)^[\t ]*$`)
)

func addIndent(s string) string {
	res := Indent + strings.TrimRight(strings.ReplaceAll(s, "\n", "\n"+Indent), Indent)
	res = repBlankLine.ReplaceAllString(res, "")
	return res
}

type Root struct {
	Rootentities []*RootEntity `@@*`
}

func (r Root) String() string {
	result := ""
	for _, r := range r.Rootentities {
		result += r.String()
	}
	return result
}

// Root内で出現しうる記述
type RootEntity struct {
	Tokens     []lexer.Token
	Definition *Definition `( @@ `
	Function   *Func       `| @@)`
}

func (r RootEntity) String() string {
	comments := ""
	inSub := 0
	for _, t := range r.Tokens {
		if dict[t.Type] == "Function" {
			inSub++
		} else if dict[t.Type] == "FuncEnd" {
			inSub--
		}

		if inSub == 0 {
			if dict[t.Type] == "CommentOneLine" {
				comments += t.Value + "\n"
			} else if dict[t.Type] == "CommentMultiLine" {
				comments += t.Value + "\n"
			} else if dict[t.Type] == "BlankLine" {
				comments += "\n"
			}
		}
	}

	result := ""
	if r.Definition != nil {
		result += r.Definition.String() + "\n"
	} else {
		result += r.Function.String()
	}
	return comments + result
}

type Comment struct {
	CommentOneLine   string `( @CommentOneLine`
	CommentMultiLine string `| @CommentMultiLine)`
}

func (c Comment) String() string {
	return c.CommentOneLine + c.CommentMultiLine
}

type Definition struct {
	DefinitionSpace *DefinitionSpace `( @@`
	DefinitionTab   *DefinitionTab   `| @@)`
}

func (d Definition) String() string {
	if d.DefinitionSpace != nil {
		return d.DefinitionSpace.String()
	} else {
		return d.DefinitionTab.String()
	}
}

type DefinitionSpace struct {
	DefSpaceKey   string `@DefinitionSpace`
	DefSpaceValue string `@DefinitionSpaceChar`
}

func (d DefinitionSpace) String() string {
	return d.DefSpaceKey + d.DefSpaceValue
}

type DefinitionTab struct {
	DefTabKey   string `@DefinitionTab`
	DefTabValue string `@DefinitionTabChar`
}

func (d DefinitionTab) String() string {
	return d.DefTabKey + d.DefTabValue
}

type Func struct {
	FunctionName string `@FuncName`
	FunctionType string `(":" @("array"|"void"|"nonoverlap"|"sequential"))?` // @FuncTypeだと認識しないので即値で指定

	FuncEntities []*FuncEntity `"{" @@* "}"`
}

func (f Func) String() string {
	funcName := f.FunctionName
	if f.FunctionType != "" {
		funcName += " : " + f.FunctionType
	}
	if len(f.FuncEntities) == 0 {
		return funcName + "\n{\n}\n"
	} else {
		entities := ""
		for _, e := range f.FuncEntities {
			entities += e.String() + "\n"
		}
		return funcName + "\n{\n" + addIndent(entities) + "}\n"
	}
}

// Func内で出現しうる記述: 関数内に1行で取りうる式
type FuncEntity struct {
	Tokens      []lexer.Token
	OutputFixer string        `( @"--"`
	Flow        *Flow         `| @@`
	PreValue    string        `| @PreValue?`
	Value       *Expr         `  @@`
	ValueEnd    string        `  ";"?`
	Sub         []*FuncEntity `| "{" @@* "}")`
}

func (f FuncEntity) String() string {
	comments := ""
	inSub := 0
	isInFlowOneLine := false
	isInFlowPre := false
	for _, t := range f.Tokens {
		switch dict[t.Type] {
		case "Function":
			inSub++
		case "FuncEnd":
			inSub--
		case "LF", "ExprEnd":
			isInFlowOneLine = false
		case "FlowKey", "FlowKeyFor", "FlowKeyForEach", "FlowKeyConst", "FlowKeyExpr":
			isInFlowPre = true
		}

		if isInFlowPre && (dict[t.Type] == "ExprEnd" || dict[t.Type] == "LF") {
			isInFlowOneLine = true
			isInFlowPre = false
		}

		if !isInFlowOneLine && inSub == 0 {
			if dict[t.Type] == "CommentOneLine" {
				comments += t.Value + "\n"
			} else if dict[t.Type] == "CommentMultiLine" {
				comments += t.Value + "\n"
			} else if dict[t.Type] == "BlankLine" {
				comments += "\n"
			}
		}
	}

	result := ""
	if f.OutputFixer != "" {
		result = f.OutputFixer
	} else if f.Flow != nil {
		result = f.Flow.String()
	} else if f.Value != nil {
		if f.PreValue != "" {
			result = f.PreValue + " "
		}
		result += f.Value.String()
	} else {
		if len(f.Sub) == 0 {
			result = "{\n}"
		} else {
			lines := ""
			for _, s := range f.Sub {
				lines += s.String() + "\n"
			}
			result = "{\n" + addIndent(lines) + "}"
		}
	}
	return comments + result
}

// フロー制御文
type Flow struct {
	FlowKey         string       `( @FlowKey`
	FlowKeyForEach  string       `| ( @FlowKeyForEach`
	FlowExprForEach *ExprForEach `    @@)`
	FlowKeyFor      string       `| ( @FlowKeyFor`
	FlowExprFor     *ExprFor     `    @@)`
	FlowKeyConst    string       `| ( @FlowKeyConst`
	FlowConst       []*Const     `    @@ ("," @@)*)`
	FlowKeyExpr     string       `| ( @FlowKeyExpr`
	FlowExpr        *Expr        `    @@))`

	FlowMultiLineSub []*FuncEntity `( "{" @@* "}"`
	FlowOneLineSub   *FuncEntity   `| ("\n"|";") @@)`
}

func (f Flow) String() string {
	upper := ""
	if f.FlowKey != "" {
		upper = f.FlowKey
	} else if f.FlowKeyForEach != "" {
		upper = f.FlowKeyForEach + " " + f.FlowExprForEach.String()
	} else if f.FlowKeyFor != "" {
		upper = f.FlowKeyFor + " " + f.FlowExprFor.String()
	} else if f.FlowKeyConst != "" {
		consts := ""
		for _, c := range f.FlowConst {
			consts += c.String() + ","
		}
		consts = strings.TrimRight(consts, ", ")
		upper = f.FlowKeyConst + " " + consts
	} else {
		upper = f.FlowKeyExpr + " " + f.FlowExpr.String()
	}

	lower := ""
	if f.FlowOneLineSub != nil {
		lower = "\n" + addIndent(f.FlowOneLineSub.String())
	} else {
		if len(f.FlowMultiLineSub) == 0 {
			lower = " {\n}"
		} else {
			lower = ""
			for _, f := range f.FlowMultiLineSub {
				lower += f.String() + "\n"
			}
			lower = " {\n" + addIndent(lower) + "}"
		}
	}

	return upper + lower
}

type ExprFor struct {
	ForInitAsign *Asign `@@ ";"`
	ForEndExpr   *Expr  `@@ ";"`
	ForLoopAsign *Asign `@@`
}

func (e ExprFor) String() string {
	return e.ForInitAsign.String() + "; " + e.ForEndExpr.String() + "; " + e.ForLoopAsign.String()
}

type ExprForEach struct {
	Array *Enum  `@@ ";"`
	Ident string `@Ident`
}

func (e ExprForEach) String() string {
	return e.Array.String() + "; " + e.Ident
}

// 条件式
type Expr struct {
	Enum *Enum `@@`
}

func (e Expr) String() string {
	return e.Enum.String()
}

type Enum struct {
	Asign *Asign `@@`

	OperEnum string `( @","`
	Right    *Enum  `  @@)?`
}

func (e Enum) String() string {
	right := ""
	if e.Right != nil {
		right = e.OperEnum + " " + e.Right.String()
	}
	return e.Asign.String() + right
}

// 代入式
type Asign struct {
	Or *Or `@@`

	OperAsign string `( @OperAsign`
	Right     *Asign `  @@)?`
}

func (a Asign) String() string {
	right := ""
	if a.Right != nil {
		right = " " + a.OperAsign + " " + a.Right.String()
	}
	return a.Or.String() + right
}

// 論理演算式Or
type Or struct {
	And *And `@@`

	OperOr string `( @"||"`
	Right  *Or    `  @@)?`
}

func (o Or) String() string {
	right := ""
	if o.Right != nil {
		right = " " + o.OperOr + " " + o.Right.String()
	}
	return o.And.String() + right
}

// 論理演算式And
type And struct {
	Comparison *Comparison `@@`

	OperAnd string `( @"&&"`
	Right   *And   `  @@)?`
}

func (a And) String() string {
	right := ""
	if a.Right != nil {
		right = " " + a.OperAnd + " " + a.Right.String()
	}
	return a.Comparison.String() + right
}

// 比較演算式
type Comparison struct {
	Addition *Addition `@@`

	OperComp string      `( @("=" "="|"!""="|">="|"<="|">"|"<"|"_in_"|"!""_in_")`
	Right    *Comparison `  @@)?`
}

func (c Comparison) String() string {
	right := ""
	if c.Right != nil {
		right = " " + c.OperComp + " " + c.Right.String()
	}
	return c.Addition.String() + right
}

// 加減法式
type Addition struct {
	Multipulation *Multipulation `@@`

	OperAdd string    `( @("+"|"-")`
	Right   *Addition `  @@)?`
}

func (a Addition) String() string {
	right := ""
	if a.Right != nil {
		right = " " + a.OperAdd + " " + a.Right.String()
	}
	return a.Multipulation.String() + right
}

// 乗除法式
type Multipulation struct {
	Unary *Unary `@@`

	OperMulti string         `( @("*"|"/"|"%")`
	Right     *Multipulation `  @@)?`
}

func (m Multipulation) String() string {
	right := ""
	if m.Right != nil {
		right = " " + m.OperMulti + " " + m.Right.String()
	}
	return m.Unary.String() + right
}

// 単項演算式
type Unary struct {
	Unary       string   `(@OperUnary)?`
	Primary     *Primary `@@`
	OperCalcOne string   `@((?! Space|TabSpace|LF|BlankLine) ("+" "+"|"--"))?`
}

func (u Unary) String() string {
	return u.Unary + u.Primary.String() + u.OperCalcOne
}

// 単項: 左辺と右辺の両方になりうる式
type Primary struct {
	Const     *Const  `( @@`
	ArrayArgs []*Expr `  ("[" @@ "]")*`
	SubExpr   *Expr   `| "(" @@ ")")`
}

func (p Primary) String() string {
	if p.Const != nil {
		args := ""
		for _, e := range p.ArrayArgs {
			args += "[" + e.String() + "]"
		}
		return p.Const.String() + args
	} else {
		return "(" + p.SubExpr.String() + ")"
	}
}

type Const struct {
	CString  *String   `( @@`
	FuncCall *FuncCall `| @@`
	Number   *Number   `| @@`
	Ident    string    `| @Ident)`
}

func (c Const) String() string {
	if c.CString != nil {
		return c.CString.String()
	} else if c.FuncCall != nil {
		return c.FuncCall.String()
	} else if c.Number != nil {
		return c.Number.String()
	} else {
		return c.Ident
	}
}

type String struct {
	SingleQuote        string `( @"'" @SingleQuoteChar? @"'"`
	DoubleQuote        string `| @"\"" @DoubleQuoteChar? @"\""`
	HearDocumentSingle string `| @HearDocumentSingle`
	HearDocumentDouble string `| @HearDocumentDouble)`
}

func (s String) String() string {
	return s.SingleQuote + s.DoubleQuote + s.HearDocumentSingle + s.HearDocumentDouble
}

type FuncCall struct {
	FuncName string `@Ident`
	FuncArgs *Expr  `"(" @@? ")"`
}

func (f FuncCall) String() string {
	result := f.FuncName + "("

	args := ""
	if f.FuncArgs != nil {
		for _, a := range strings.Split(f.FuncArgs.String(), ", ") {
			args += strings.ReplaceAll(a, " ", "") + ", "
		}
		args = strings.TrimRight(args, ", ")
	}

	result += args + ")"
	return result
}

// 目的はパースなのでstringでとっていい
type Number struct {
	Hex   string `( @HexNum`
	Bin   string `| @BinNum`
	Float string `| @Float`
	Int   string `| @Int)`
}

func (n Number) String() string {
	return n.Hex + n.Bin + n.Float + n.Int
}
