package main

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var (
	def = lexer.MustStateful(lexer.Rules{
		"Root": { // Rootは特別な名前　これが基点のRule
			{Name: "LF", Pattern: `(\r\n|\r|\n)`, Action: nil},
			{Name: `Space`, Pattern: ` +`, Action: nil},
			{Name: `TabSpace`, Pattern: `	+`, Action: nil},
			{Name: `Definition`, Pattern: `(#globaldefine|#define)`, Action: lexer.Push("DefinitionRule")},
			{Name: "FuncName", Pattern: `[a-zA-Z_.][^ \r\n!"#$%&()*+,-/:;<=>?@\[\]{|}]*`, Action: nil}, //TODO:`も禁止文字に入れる
			{Name: "Colon", Pattern: `:`, Action: nil},
			{Name: "FuncType", Pattern: `(array|void)`, Action: nil},
			{Name: `Function`, Pattern: `\{`, Action: lexer.Push("FuncRule")},
			lexer.Include("Comments"),
		},
		"DefinitionRule": {
			{Name: `Space`, Pattern: ` +`, Action: nil},
			{Name: `TabSpace`, Pattern: `	+`, Action: nil},
			{Name: `DefinitionChar`, Pattern: `[^\r\n ]+`, Action: nil},
			{Name: `DefinitionEnd`, Pattern: `(\r\n|\r|\n|;)`, Action: lexer.Pop()},
		},
		"FuncRule": {
			{Name: `Space`, Pattern: ` +`, Action: nil},
			{Name: `TabSpace`, Pattern: `	+`, Action: nil},
			{Name: "LF", Pattern: `(\r\n|\r|\n)`, Action: nil},
			lexer.Include("Comments"),
			{Name: `OperAsign`, Pattern: `(=|:=|\+=|-=|\*=|/=|%=|\+:=|-:=|\*:=|/:=|%:=|,=)`, Action: nil},
			{Name: `OperAsignUnary`, Pattern: `[^ 	](--|\+\+)`, Action: nil},
			{Name: `Separator`, Pattern: `--`, Action: nil},
			{Name: `OperUnary`, Pattern: `!`, Action: nil}, // TODO: -はない？要確認
			{Name: `OperCalc`, Pattern: `(-|\+|/|\*|%)`, Action: nil},
			{Name: `OperComp`, Pattern: `(==|!=|>=|<=|>|<|_in_|!_in_)`, Action: nil},
			{Name: `OperLogic`, Pattern: `(\|\||&&)`, Action: nil},
			{Name: "FlowKeyExpr", Pattern: `(while|if|elseif|case)`, Action: nil},
			{Name: "FlowKeyPrimary", Pattern: `when`, Action: nil},
			{Name: "FlowKey", Pattern: `else`, Action: nil},
			{Name: `Function`, Pattern: `\{`, Action: lexer.Push("FuncRule")},
			{Name: "Number", Pattern: `\d+(\.\d+)?`, Action: nil},
			{Name: "Ident", Pattern: `[a-zA-Z_.][^ \r\n!"#$%&()*+,-/:;<=>?@\[\]{|}]*`, Action: nil}, //TODO:`も禁止文字に入れる
			{Name: "SingleQuoteString", Pattern: `'`, Action: lexer.Push("SingleQuoteStringRule")},
			{Name: "DoubleQuoteString", Pattern: `"`, Action: lexer.Push("DoubleQuoteStringRule")},
			{Name: `ExprEnd`, Pattern: `;`, Action: nil},
			{Name: "FuncCall", Pattern: `\(`, Action: lexer.Push("FuncCallRule")},
			{Name: "ArrayCall", Pattern: `\[`, Action: lexer.Push("ArrayCallRule")},
			{Name: "FuncEnd", Pattern: `\}`, Action: lexer.Pop()},
		},
		"SingleQuoteStringRule": {
			{Name: "SingleQuoteChar", Pattern: `[^']+`, Action: nil},
			{Name: "SingleQuoteStringEnd", Pattern: `'`, Action: lexer.Pop()},
		},
		"DoubleQuoteStringRule": {
			{Name: "DoubleQuoteChar", Pattern: `[^"]+`, Action: nil},
			{Name: "DoubleQuoteStringEnd", Pattern: `"`, Action: lexer.Pop()},
		},
		"FuncCallRule": { // 条件式のブラケットと関数呼び出しのカッコが混在している
			lexer.Include("FuncRule"),
			{Name: "Delim", Pattern: `,`, Action: nil}, //関数呼び出し用
			{Name: "FuncCallEnd", Pattern: `\)`, Action: lexer.Pop()},
		},
		"ArrayCallRule": {
			lexer.Include("FuncRule"),
			{Name: "ArrayCallEnd", Pattern: `\]`, Action: lexer.Pop()},
		},
		"Comments": {
			{Name: `CommentOneLine`, Pattern: `//.*(\r\n|\r|\n)`, Action: nil},
			{Name: "CommentMultiLine", Pattern: `/\*([\r\n	 ]|.)*\*/`, Action: nil},
		},
	})
	parser = participle.MustBuild[Root](
		participle.Lexer(def),
		participle.Elide("Space"),
		participle.Elide("TabSpace"),
		participle.Elide("LF"),
	)
)
