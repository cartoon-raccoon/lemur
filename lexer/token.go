package lexer

type TokenType string

type Context struct {
	Line int
	Col  int
}

type Token struct {
	Type    string
	Literal string
	Pos     Context
}

func newToken(
	ttype string,
	lit string,
	line int,
	col int,
) Token {
	return Token{
		Type:    ttype,
		Literal: lit,
		Pos:     newContext(line, col),
	}
}

func newContext(line int, col int) Context {
	return Context{
		Line: line,
		Col:  col,
	}
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	//Identifiers and literals
	IDENT = "IDENT"
	INT   = "INT"

	//Operators
	ASSIGN = "="

	ADD = "+"
	SUB = "-"
	MUL = "*"
	DIV = "/"

	ADDASSIGN = "+="
	SUBASSIGN = "-="
	MULASSIGN = "*="
	DIVASSIGN = "/="

	BSL = "<<"
	BSR = ">>"

	BWOR  = "|"
	BWAND = "&"
	BWNOT = "~"

	BWOASSIGN = "|="
	BWAASSIGN = "&="
	BWNASSIGN = "~="

	LT = "<"
	GT = ">"
	LE = "<="
	GE = ">="
	EQ = "=="
	NE = "!="

	LOR  = "||"
	LAND = "&&"

	//Delimiters
	COMMA   = ","
	SEMICOL = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	//Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	RETURN   = "->"
)

var keywords = map[string]string{
	"fn":  FUNCTION,
	"let": LET,
}
