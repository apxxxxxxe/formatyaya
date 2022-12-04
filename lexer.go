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
			{Name: `DefinitionSpace`, Pattern: `(#globaldefine|#define) +`, Action: lexer.Push("DefinitionSpaceRule")},
			{Name: `DefinitionTab`, Pattern: `(#globaldefine|#define)	+`, Action: lexer.Push("DefinitionTabRule")},
			{Name: "FuncName", Pattern: `([a-zA-Z_\x{3041}-\x{3096}\x{30A1}-\x{30FA}々〇〻\x{3400}-\x{9FFF}\x{F900}-\x{FAFF}]|[\x{D840}-\x{D87F}][\x{DC00}-\x{DFFF}])([^ \n!"#$%&()*+,-/:;<=>?@\[\]{|}]|\.)*`, Action: nil}, //TODO:`も禁止文字に入れる
			{Name: "Colon", Pattern: `:`, Action: nil},
			{Name: "FuncType", Pattern: `(array|void|nonoverlap|sequential)`, Action: nil},
			{Name: `Function`, Pattern: `\{`, Action: lexer.Push("FuncRule")},
			lexer.Include("Comments"),
		},
		"DefinitionSpaceRule": {
			{Name: `Space`, Pattern: ` +`, Action: nil},
			{Name: `DefinitionSpaceChar`, Pattern: `[^\n ]+`, Action: nil},
			{Name: `LineEnd`, Pattern: `(\n|;)`, Action: lexer.Pop()},
		},
		"DefinitionTabRule": {
			{Name: `TabSpace`, Pattern: `	+`, Action: nil},
			{Name: `DefinitionTabChar`, Pattern: `[^\n\t]+`, Action: nil},
			{Name: `LineEnd`, Pattern: `(\n|;)`, Action: lexer.Pop()},
		},
		"FuncRule": {
			{Name: `PreValue`, Pattern: `(void)`, Action: nil},
			{Name: `BlankLine`, Pattern: `^[ 	]*\n`, Action: nil},
			{Name: `Space`, Pattern: ` +`, Action: nil},
			{Name: `TabSpace`, Pattern: `	+`, Action: nil},
			{Name: `ExprEnd`, Pattern: `;`, Action: nil},
			{Name: "LF", Pattern: `\n`, Action: nil},
			lexer.Include("Comments"),
			{Name: "SingleQuoteString", Pattern: `'`, Action: lexer.Push("SingleQuoteStringRule")},
			{Name: "DoubleQuoteString", Pattern: `"`, Action: lexer.Push("DoubleQuoteStringRule")},
			{Name: "HearDocumentSingle", Pattern: `<<'([\n	 ]|.)*'>>`, Action: nil},
			{Name: "HearDocumentDouble", Pattern: `<<"([\n	 ]|.)*">>`, Action: nil},
			{Name: `OperAsign`, Pattern: `(=|:=|\+=|-=|\*=|/=|%=|\+:=|-:=|\*:=|/:=|%:=|,=)`, Action: nil},
			{Name: `OperAsignUnary`, Pattern: `[^ 	](--|\+\+)`, Action: nil},
			{Name: `OutputFixer`, Pattern: `--`, Action: nil},
			{Name: `OperUnary`, Pattern: `(!|-)`, Action: nil},
			{Name: `OperCalc`, Pattern: `(-|\+|/|\*|%)`, Action: nil},
			{Name: `OperComp`, Pattern: `(==|!=|>=|<=|>|<|_in_|!_in_)`, Action: nil},
			{Name: `OperLogic`, Pattern: `(\|\||&&)`, Action: nil},
			{Name: "OperEnum", Pattern: `,`, Action: nil},
			{Name: "FlowKeyExpr", Pattern: `(while|if|elseif|case|switch)`, Action: nil},
			{Name: "FlowKeyConst", Pattern: `when`, Action: nil},
			{Name: "FlowKeyForEach", Pattern: `foreach `, Action: nil},
			{Name: "FlowKeyFor", Pattern: `for `, Action: nil},
			{Name: "FlowKey", Pattern: `else`, Action: nil},
			{Name: `Function`, Pattern: `\{`, Action: lexer.Push("FuncRule")},
			{Name: "HexNum", Pattern: `0x[0-9A-za-z]+`, Action: nil},
			{Name: "BinNum", Pattern: `0b[01]+`, Action: nil},
			{Name: "Float", Pattern: `\d+\.\d+`, Action: nil},
			{Name: "Int", Pattern: `\d+`, Action: nil},
			{Name: "Ident", Pattern: `([a-zA-Z_\x{3041}-\x{3096}\x{30A1}-\x{30FA}々〇〻\x{3400}-\x{9FFF}\x{F900}-\x{FAFF}]|[\x{D840}-\x{D87F}][\x{DC00}-\x{DFFF}])([^ \n!"#$%&()*+,-/:;<=>?@\[\]{|}]|\.)*`, Action: nil}, //TODO:`も禁止文字に入れる

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
		"FuncCallRule": {
			lexer.Include("FuncRule"),
			{Name: "FuncCallEnd", Pattern: `\)`, Action: lexer.Pop()},
		},
		"ArrayCallRule": {
			lexer.Include("FuncRule"),
			{Name: "ArrayCallEnd", Pattern: `\]`, Action: lexer.Pop()},
		},
		"Comments": {
			{Name: `CommentOneLine`, Pattern: `//[^\n]*`, Action: nil},
			{Name: "CommentMultiLine", Pattern: `/\*(\n|.)*\*/`, Action: nil},
		},
	})
	parser = participle.MustBuild[Root](
		participle.Lexer(def),
		participle.Elide("Space"),
		participle.Elide("TabSpace"),
		participle.Elide("CommentOneLine"),
		participle.Elide("CommentMultiLine"),
	)
)
