package lexer

// Context is the place of that particular token in the text
type Context struct {
	Line int
	Col  int
	Ctxt string
}

// Token represents a single word in Monkey
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
	ctx string,
) Token {
	return Token{
		Type:    ttype,
		Literal: lit,
		Pos:     newContext(line, col, ctx),
	}
}

func (tok *Token) isEOF() bool {
	return tok.Type == EOF
}

func newContext(line int, col int, ctx string) Context {
	return Context{
		Line: line,
		Col:  col,
		Ctxt: ctx,
	}
}

//Types of tokens
const (
	// ILLEGAL - Unknown token
	ILLEGAL = "ILLEGAL"

	// EOF - End of file
	EOF = "EOF"

	//Identifiers and literals

	// IDENT - Identifier (variable/function name)
	IDENT = "IDENT"
	// STRLIT - String literal
	STRLIT = "STRLIT"
	// INTLIT - Integer literal
	INTLIT = "INTLIT"
	// FLTLIT - Float literal
	FLTLIT = "FLTLIT"

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

	DOT = "."

	LT = "<"
	GT = ">"
	LE = "<="
	GE = ">="
	EQ = "=="
	NE = "!="

	LOR  = "||"
	LAND = "&&"

	BANG = "!"

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

	FUNCTION = "fn"
	LET      = "let"
	RETURN   = "return"
	IF       = "if"
	ELSE     = "else"
	WHILE    = "while"
	FOR      = "for"
	IN       = "in"
	LOOP     = "loop"
	STRING   = "str"
	INT      = "int"
	FLOAT    = "flt"
	CLASS    = "class"
	BOOL     = "bool"
	TRUE     = "true"
	FALSE    = "false"
)

var keywords = map[string]string{
	"fn":     FUNCTION,
	"let":    LET,
	"return": RETURN,
	"if":     IF,
	"else":   ELSE,
	"while":  WHILE,
	"for":    FOR,
	"in":     IN,
	"loop":   LOOP,
	"int":    INT,
	"float":  FLOAT,
	"class":  CLASS,
	"bool":   BOOL,
	"true":   TRUE,
	"false":  FALSE,
}
