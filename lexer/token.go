package lexer

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

func (tok *Token) isEOF() bool {
	return tok.Type == EOF
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

	IDENT  = "IDENT"
	INT    = "INT"
	STRLIT = "STRLIT"

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
	LSBRKT = "["
	RSBRKT = "]"

	//Keywords

	FUNCTION = "FUNCTION"
	LET      = "LET"
	RETURN   = "RETURN"
)

var keywords = map[string]string{
	"fn":     FUNCTION,
	"let":    LET,
	"return": RETURN,
}
