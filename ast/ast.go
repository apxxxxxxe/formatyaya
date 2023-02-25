package ast

import (
	"regexp"
	"strings"
)

var (
	Indent       string
	repBlankLine = regexp.MustCompile(`(?m)^[\t ]*$`)
	repIndents   = regexp.MustCompile(`(?m)^[\t ]+`)
)

func addIndent(s string) string {
	res := Indent + strings.TrimRight(strings.ReplaceAll(s, "\n", "\n"+Indent), Indent)
	res = repBlankLine.ReplaceAllString(res, "")
	return res
}

// 記述の半角スペースを削除する(ただし文字列内のスペースはそのまま)
func deleteSpace(str string) string {
	result := ""
	isinQuotes := false
	for _, s := range str {
		if s == '\'' || s == '"' {
			isinQuotes = !isinQuotes
		}

		if isinQuotes || s != ' ' {
			result += string(s)
		}
	}
	return result
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
	Definition   *Definition `(( @@ `
	Function     *Func       ` | @@`
	BlankLine    string      ` | @BlankLine)`
	CommentDep   *Comment    `   @@?)`
	CommentIndep *Comment    `| @@`
}

func (r RootEntity) String() string {
	result := ""
	if r.Definition != nil {
		result += r.Definition.String() + "\n"
	} else if r.Function != nil {
		result += r.Function.String()
	} else if r.CommentIndep != nil {
		result += r.CommentIndep.String()
	} else if r.BlankLine != "" {
		result += "\n"
	}
	comment := ""
	if r.CommentDep != nil {
		comment = r.CommentDep.String()
	}
	return result + comment
}

type Comment struct {
	CommentOneLine   string `( @CommentOneLine`
	CommentMultiLine string `| @CommentMultiLine)`
}

func (c Comment) String() string {
	return c.CommentOneLine + c.CommentMultiLine + "\n"
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
	OutputFixer  string        `(( @"--"`
	Flow         *Flow         ` | @@`
	PreValue     string        ` | @PreValue?`
	Value        *Expr         `   @@`
	ValueEnd     string        `   ";"?`
	Sub          []*FuncEntity ` | "{" @@* "}"`
	BlankLine    string        ` | @BlankLine)`
	CommentDep   *Comment      `   @@?)`
	CommentIndep *Comment      `| @@`
}

func (f FuncEntity) String() string {
	result := ""
	if f.CommentIndep != nil {
		result += f.CommentIndep.String() + "piyo"
	} else if f.OutputFixer != "" {
		result = f.OutputFixer
	} else if f.Flow != nil {
		result = f.Flow.String()
	} else if f.Value != nil {
		if f.PreValue != "" {
			result = f.PreValue + " "
		}
		result += f.Value.String()
	} else if f.BlankLine != "" {
		result += "\n"
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
	comment := ""
	if f.CommentDep != nil {
		comment = f.CommentDep.String()
	}
	return result + comment
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
		upper = f.FlowKeyForEach + f.FlowExprForEach.String()
	} else if f.FlowKeyFor != "" {
		upper = f.FlowKeyFor + f.FlowExprFor.String()
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
	return deleteSpace(e.ForInitAsign.String()) + "; " + deleteSpace(e.ForEndExpr.String()) + "; " + e.ForLoopAsign.String()
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
	Slash    string `  @Slash?`
	Right    *Enum  `  @@)?`
}

func (e Enum) String() string {
	right := ""
	if e.Right != nil {
		space := " "
		if e.Slash != "" {
			space += e.Slash + Indent
		}
		right = e.OperEnum + space + e.Right.String()
	}
	return e.Asign.String() + right
}

// 代入式
type Asign struct {
	Or *Or `@@`

	OperAsign string `( @OperAsign`
	Slash     string `  @Slash?`
	Right     *Asign `  @@)?`
}

func (a Asign) String() string {
	right := ""
	if a.Right != nil {
		space := " "
		if a.Slash != "" {
			space += a.Slash + Indent
		}
		right = " " + a.OperAsign + space + a.Right.String()
	}
	return a.Or.String() + right
}

// 論理演算式Or
type Or struct {
	And *And `@@`

	OperOr string `( @"||"`
	Slash  string `  @Slash?`
	Right  *Or    `  @@)?`
}

func (o Or) String() string {
	right := ""
	if o.Right != nil {
		space := " "
		if o.Slash != "" {
			space += o.Slash + Indent
		}
		right = " " + o.OperOr + space + o.Right.String()
	}
	return o.And.String() + right
}

// 論理演算式And
type And struct {
	Comparison *Comparison `@@`

	OperAnd string `( @"&&"`
	Slash   string `  @Slash?`
	Right   *And   `  @@)?`
}

func (a And) String() string {
	right := ""
	if a.Right != nil {
		space := " "
		if a.Slash != "" {
			space += a.Slash + Indent
		}
		right = " " + a.OperAnd + space + a.Right.String()
	}
	return a.Comparison.String() + right
}

// 比較演算式
type Comparison struct {
	Addition *Addition `@@`

	OperComp string      `( @("=" "="|"!""="|">="|"<="|">"|"<"|"_in_"|"!""_in_")`
	Slash    string      `  @Slash?`
	Right    *Comparison `  @@)?`
}

func (c Comparison) String() string {
	right := ""
	if c.Right != nil {
		space := " "
		if c.Slash != "" {
			space += c.Slash + Indent
		}
		right = " " + c.OperComp + space + c.Right.String()
	}
	return c.Addition.String() + right
}

// 加減法式
type Addition struct {
	Multipulation *Multipulation `@@`

	OperAdd string    `( @("+"|"-")`
	Slash   string    `  @Slash?`
	Right   *Addition `  @@)?`
}

func (a Addition) String() string {
	right := ""
	if a.Right != nil {
		space := " "
		if a.Slash != "" {
			space += a.Slash + Indent
		}
		right = " " + a.OperAdd + space + a.Right.String()
	}
	return a.Multipulation.String() + right
}

// 乗除法式
type Multipulation struct {
	Unary *Unary `@@`

	OperMulti string         `( @("*"|"/"|"%")`
	Slash     string         `  @Slash?`
	Right     *Multipulation `  @@)?`
}

func (m Multipulation) String() string {
	right := ""
	if m.Right != nil {
		space := " "
		if m.Slash != "" {
			space += m.Slash + Indent
		}
		right = " " + m.OperMulti + space + m.Right.String()
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
	Const     *Const      `( @@`
	ArrayArgs []*ArrayArg `  @@*`
	SubExpr   *SubExpr    `| @@)`
}

func (p Primary) String() string {
	if p.Const != nil {
		args := ""
		for _, a := range p.ArrayArgs {
			args += a.String()
		}
		return p.Const.String() + args
	} else {
		return p.SubExpr.String()
	}
}

type ArrayArg struct {
	Start      string `"["`
	StartSlash string `@Slash?`
	Expr       *Expr  `@@`
	EndSlash   string `@Slash?`
	End        string `"]"`
}

func (a ArrayArg) String() string {
	start := a.StartSlash
	if start != "" {
		start = " " + start + Indent
	}
	end := a.EndSlash
	if end != "" {
		end = " " + end
	}
	return "[" + start + a.Expr.String() + end + "]"
}

type SubExpr struct {
	Start      string `"("`
	StartSlash string `@Slash?`
	Expr       *Expr  `@@`
	EndSlash   string `@Slash?`
	End        string `")"`
}

func (s SubExpr) String() string {
	start := s.StartSlash
	if start != "" {
		start = " " + start + Indent
	}
	end := s.EndSlash
	if end != "" {
		end = " " + end
	}
	return "(" + start + s.Expr.String() + end + ")"
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
	if s.SingleQuote != "" {
		return s.SingleQuote
	} else if s.DoubleQuote != "" {
		if !strings.Contains(s.DoubleQuote, "%") {
			return "'" + strings.Trim(s.DoubleQuote, "\"") + "'"
		} else {
			return s.DoubleQuote
		}
	} else if s.HearDocumentSingle != "" {
		return repIndents.ReplaceAllString(s.HearDocumentSingle, "")
	} else {
		return repIndents.ReplaceAllString(s.HearDocumentDouble, "")
	}
}

type FuncCall struct {
	FuncName string    `@Ident`
	FuncExpr *FuncExpr `@@`
}

func (f FuncCall) String() string {
	return f.FuncName + f.FuncExpr.String()
}

type FuncExpr struct {
	Start      string `"("`
	StartSlash string `@Slash?`
	Expr       *Expr  `@@?`
	EndSlash   string `@Slash?`
	End        string `")"`
}

func (f FuncExpr) String() string {
	start := f.StartSlash
	if start != "" {
		start = " " + start + Indent
	}
	expr := ""
	if f.Expr != nil {
		tmpIndent := "\t"
		exprBase := strings.ReplaceAll(f.Expr.String(), Indent, tmpIndent)
		for _, a := range strings.Split(exprBase, ", ") {
			expr += deleteSpace(a) + ", "
		}
		expr = strings.ReplaceAll(strings.TrimRight(expr, ", "), tmpIndent, Indent)
	}
	end := f.EndSlash
	if end != "" {
		end = " " + end
	}
	return "(" + start + expr + end + ")"
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
