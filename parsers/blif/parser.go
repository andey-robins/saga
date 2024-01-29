package blif

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

func NewBlifParser() (*participle.Parser[BlifFile], error) {
	basicLexer := lexer.MustSimple([]lexer.SimpleRule{
		{`Comment`, `# [^\n]*`},
		{`Ident`, `\w*[a-zA-Z_0-9/.]\w*`},
		{`whitespace`, `[ \t]+`},
		{`Equality`, `=`},
		{`NewLine`, `[\r\n]+`},
	})

	return participle.Build[BlifFile](
		participle.Lexer(basicLexer),
	)
}
