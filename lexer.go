package main

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var (
	def = lexer.MustStateful(lexer.Rules{
		"Root": { // Rootは特別な名前　これが基点のRule
			{Name: "LF", Pattern: `\n`, Action: nil},
			{Name: `Space`, Pattern: ` +`, Action: nil},
			{Name: `TabSpace`, Pattern: `	+`, Action: nil},
			{Name: `Definition`, Pattern: `(#globaldefine|#define)`, Action: lexer.Push("DefinitionRule")},
			{Name: "FuncName", Pattern: `[a-zA-Z_]([^ \n!"#$%&()*+,-/:;<=>?@\[\]{|}]|\.)*`, Action: nil}, //TODO:`も禁止文字に入れる
			{Name: "Colon", Pattern: `:`, Action: nil},
			{Name: "FuncType", Pattern: `(array|void)`, Action: nil},
			{Name: `Function`, Pattern: `\{`, Action: lexer.Push("FuncRule")},
			lexer.Include("Comments"),
		},
		"DefinitionRule": {
			{Name: `Space`, Pattern: ` +`, Action: nil},
			{Name: `TabSpace`, Pattern: `	+`, Action: nil},
			{Name: `DefinitionChar`, Pattern: `[^\n\t ]+`, Action: nil},
			{Name: `ExprEnd`, Pattern: `;`, Action: nil},
			{Name: `LineEnd`, Pattern: `\n`, Action: lexer.Pop()},
		},
		"FuncRule": {
			{Name: `Space`, Pattern: ` +`, Action: nil},
			{Name: `TabSpace`, Pattern: `	+`, Action: nil},
			{Name: `ExprEnd`, Pattern: `;`, Action: nil},
			{Name: "LF", Pattern: `\n`, Action: nil},
			lexer.Include("Comments"),
			{Name: `OperAsign`, Pattern: `(=|:=|\+=|-=|\*=|/=|%=|\+:=|-:=|\*:=|/:=|%:=|,=)`, Action: nil},
			{Name: `OperAsignUnary`, Pattern: `[^ 	](--|\+\+)`, Action: nil},
			{Name: `Separator`, Pattern: `--`, Action: nil},
			{Name: `OperUnary`, Pattern: `(!|-)`, Action: nil},
			{Name: `OperCalc`, Pattern: `(-|\+|/|\*|%)`, Action: nil},
			{Name: `OperComp`, Pattern: `(==|!=|>=|<=|>|<|_in_|!_in_)`, Action: nil},
			{Name: `OperLogic`, Pattern: `(\|\||&&)`, Action: nil},
			{Name: "OperEnum", Pattern: `,`, Action: nil},
			{Name: "FlowKeyExpr", Pattern: `(while|if|elseif|case)`, Action: nil},
			{Name: "FlowKeyConst", Pattern: `when`, Action: nil},
			{Name: "FlowKeyFor", Pattern: `for`, Action: nil},
			{Name: "FlowKey", Pattern: `else`, Action: nil},
			{Name: `Function`, Pattern: `\{`, Action: lexer.Push("FuncRule")},
			{Name: "HexNum", Pattern: `0x[0-9A-za-z]+`, Action: nil},
			{Name: "Number", Pattern: `\d+(\.\d+)?`, Action: nil},
			{Name: "Ident", Pattern: `[a-zA-Z_]([^ \n!"#$%&()*+,-/:;<=>?@\[\]{|}]|\.)*`, Action: nil}, //TODO:`も禁止文字に入れる
			{Name: "SingleQuoteString", Pattern: `'`, Action: lexer.Push("SingleQuoteStringRule")},
			{Name: "DoubleQuoteString", Pattern: `"`, Action: lexer.Push("DoubleQuoteStringRule")},
			{Name: "FuncCall", Pattern: `\(`, Action: lexer.Push("FuncCallRule")},
			{Name: "ArrayCall", Pattern: `\[`, Action: lexer.Push("ArrayCallRule")},
			{Name: "FuncEnd", Pattern: `\}`, Action: lexer.Pop()},
		},
		"SingleQuoteStringRule": {
			{Name: "LF", Pattern: `\n`, Action: nil},
			{Name: "ContinuousLF", Pattern: `/$`, Action: nil},
			{Name: "SingleQuoteChar", Pattern: `[^'\n]+`, Action: nil},
			{Name: "SingleQuoteStringEnd", Pattern: `'`, Action: lexer.Pop()},
		},
		"DoubleQuoteStringRule": {
			{Name: "LF", Pattern: `\n`, Action: nil},
			{Name: "ContinuousLF", Pattern: `/$`, Action: nil},
			{Name: "DoubleQuoteChar", Pattern: `[^"\n]+`, Action: nil},
			{Name: "DoubleQuoteStringEnd", Pattern: `"`, Action: lexer.Pop()},
		},
		"FuncCallRule": {
			lexer.Include("FuncRule"),
			{Name: "FuncCallEnd", Pattern: `\)`, Action: lexer.Pop()},
		},
		"ArrayCallRule": {
			lexer.Include("FuncRule"),
			{Name: "ArrayCallEnd", Pattern: `\]`, Action: lexer.Pop()},
		},
		"Comments": {
			{Name: `CommentOneLine`, Pattern: `//.*\n`, Action: nil},
			{Name: "CommentMultiLine", Pattern: `/\*([\n	 ]|.)*\*/`, Action: nil},
		},
	})
	parser = participle.MustBuild[Root](
		participle.Lexer(def),
		participle.Elide("Space"),
		participle.Elide("TabSpace"),
		participle.Elide("LF"),
	)
)
