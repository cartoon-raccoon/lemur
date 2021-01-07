package lexer

import (
	"fmt"
)

/*------------Lexer------------*/

// Lexer represents the FSM that tokenizes the input to the interpreter
type Lexer struct {
	input   string
	line    int
	col     int
	pos     int
	context string
	lastln  int
	readPos int
	ch      byte
}

// New returns a new uninitialized lexer
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.nextChar()
	return l
}

// Tokenize fully advances the lexer and returns a slice of tokens
func (l *Lexer) Tokenize() (tokens []Token, err error) {
	tokens = make([]Token, 0)
	tok, err := l.NextToken()
	for !tok.isEOF() {
		if err != nil {
			return
		}
		tok, err = l.NextToken()
		tokens = append(tokens, tok)
	}
	return
}

// NextToken advances the lexer and produces a token
func (l *Lexer) NextToken() (Token, error) {
	l.skipWhitespace()

	switch {

	case l.ch == '=':
		l.nextChar()
		if l.ch == '=' {
			l.nextChar()
			return newToken(EQ, EQ, l.line, l.col, l.context), nil
		}
		return newToken(ASSIGN, ASSIGN, l.line, l.col, l.context), nil

	case l.ch == '+':
		l.nextChar()
		if l.ch == '=' {
			l.nextChar()
			return newToken(ADDASSIGN, ADDASSIGN, l.line, l.col, l.context), nil
		}
		return newToken(ADD, ADD, l.line, l.col, l.context), nil

	case l.ch == '-':
		l.nextChar()
		switch l.ch {
		case '=':
			l.nextChar()
			return newToken(SUBASSIGN, SUBASSIGN, l.line, l.col, l.context), nil
		case '>':
			l.nextChar()
			return newToken(RETSIG, RETSIG, l.line, l.col, l.context), nil
		default:
			return newToken(SUB, SUB, l.line, l.col, l.context), nil
		}

	case l.ch == '*':
		l.nextChar()
		if l.ch == '=' {
			l.nextChar()
			return newToken(MULASSIGN, MULASSIGN, l.line, l.col, l.context), nil
		}
		return newToken(MUL, MUL, l.line, l.col, l.context), nil

	case l.ch == '/':
		l.nextChar()
		if l.ch == '=' {
			l.nextChar()
			return newToken(DIVASSIGN, DIVASSIGN, l.line, l.col, l.context), nil
		}
		return newToken(DIV, DIV, l.line, l.col, l.context), nil

	case l.ch == '>':
		l.nextChar()

		switch l.ch {
		case '>':
			l.nextChar()
			return newToken(BSR, BSR, l.line, l.col, l.context), nil
		case '=':
			l.nextChar()
			return newToken(GE, GE, l.line, l.col, l.context), nil
		default:
			return newToken(GT, GT, l.line, l.col, l.context), nil
		}

	case l.ch == '<':
		l.nextChar()
		switch l.ch {
		case '<':
			l.nextChar()
			return newToken(BSL, BSL, l.line, l.col, l.context), nil
		case '=':
			l.nextChar()
			return newToken(LE, LE, l.line, l.col, l.context), nil
		default:
			return newToken(LT, LT, l.line, l.col, l.context), nil
		}

	case l.ch == '&':
		l.nextChar()
		switch l.ch {
		case '&':
			l.nextChar()
			return newToken(LAND, LAND, l.line, l.col, l.context), nil
		case '=':
			l.nextChar()
			return newToken(BWAASSIGN, BWAASSIGN, l.line, l.col, l.context), nil
		default:
			return newToken(BWAND, BWAND, l.line, l.col, l.context), nil
		}

	case l.ch == '|':
		l.nextChar()
		switch l.ch {
		case '|':
			l.nextChar()
			return newToken(LOR, LOR, l.line, l.col, l.context), nil
		case '=':
			l.nextChar()
			return newToken(BWOASSIGN, BWOASSIGN, l.line, l.col, l.context), nil
		default:
			return newToken(BWOR, BWOR, l.line, l.col, l.context), nil
		}

	case l.ch == '~':
		l.nextChar()
		if l.ch == '=' {
			l.nextChar()
			return newToken(BWNASSIGN, BWNASSIGN, l.line, l.col, l.context), nil
		}
		return newToken(BWNOT, BWNOT, l.line, l.col, l.context), nil

	case l.ch == ',':
		l.nextChar()
		return newToken(COMMA, COMMA, l.line, l.col, l.context), nil

	case l.ch == ';':
		l.nextChar()
		return newToken(SEMICOL, SEMICOL, l.line, l.col, l.context), nil

	case l.ch == '(':
		l.nextChar()
		return newToken(LPAREN, LPAREN, l.line, l.col, l.context), nil

	case l.ch == ')':
		l.nextChar()
		return newToken(RPAREN, RPAREN, l.line, l.col, l.context), nil

	case l.ch == '[':
		l.nextChar()
		return newToken(LSBRKT, LSBRKT, l.line, l.col, l.context), nil

	case l.ch == ']':
		l.nextChar()
		return newToken(RSBRKT, RSBRKT, l.line, l.col, l.context), nil

	case l.ch == '{':
		l.nextChar()
		return newToken(LBRACE, LBRACE, l.line, l.col, l.context), nil

	case l.ch == '}':
		l.nextChar()
		return newToken(RBRACE, RBRACE, l.line, l.col, l.context), nil

	case l.ch == '.':
		l.nextChar()
		return newToken(DOT, DOT, l.line, l.col, l.context), nil

	case l.ch == '!':
		l.nextChar()
		if l.ch == '=' {
			l.nextChar()
			return newToken(NE, NE, l.line, l.col, l.context), nil
		}
		return newToken(BANG, BANG, l.line, l.col, l.context), nil

	case l.ch == '"':
		l.nextChar()
		return l.readStrLiteral(), nil

	case l.ch == 0:
		return newToken(EOF, EOF, l.line, l.col, l.context), nil

	default:
		if isLetter(l.ch) {
			return l.readIdent(), nil
		}
		if isNumber(l.ch) {
			return l.readNumLit(), nil
		}
		return newToken("", "", 0, 0, ""), Err{
			Msg: "Unknown token",
			Con: newContext(l.line, l.col, l.context),
		}
	}
}

func (l *Lexer) readIdent() Token {
	position := l.pos
	for isAlnum(l.ch) {
		l.nextChar()
	}
	token := l.input[position:l.pos]
	if tok, isKw := lookupKeyword(token); isKw {
		return newToken(tok, tok, l.line, l.col, l.context)
	}
	return newToken(IDENT, token, l.line, l.col, l.context)
}

func (l *Lexer) readStrLiteral() Token {
	position := l.pos
	for {
		if l.ch == '"' {
			token := l.input[position:l.pos]
			l.nextChar()
			return newToken(STRLIT, token, l.line, l.col, l.context)
		} else if l.ch == '\\' { //handling escaped characters and backslashes
			l.nextChar()
			// switch l.ch {
			// case 'n':
			// 	l.input = strings.Replace(l.input, "\\n", "\n", 1)
			// case 'r':
			// 	l.input = strings.Replace(l.input, "\\r", "\r", 1)
			// case 't':
			// 	l.input = strings.Replace(l.input, "\\t", "\t", 1)
			// default:
			// 	l.input = strings.Replace(l.input, "\\", "", 1)
			// 	l.col, l.context--
			// }
		} else if l.ch == 0 {
			return newToken(EOF, EOF, l.line, l.col, l.context)
		}
		l.nextChar()
	}
}

func (l *Lexer) readNumLit() Token {
	position := l.pos
	isFloat := false
	for isNumber(l.ch) || l.ch == '_' || l.ch == '.' {
		if l.ch == '.' {
			isFloat = true
		}
		l.nextChar()
	}
	token := l.input[position:l.pos]
	if isFloat {
		return newToken(FLTLIT, token, l.line, l.col, l.context)
	}
	return newToken(INTLIT, token, l.line, l.col, l.context)
}

func (l *Lexer) nextChar() {
	//have reached EOF
	if l.readPos >= len(l.input) {
		l.ch = 0
		if l.readPos == len(l.input) {
			l.pos = l.readPos
		}
	} else {
		l.ch = l.input[l.readPos]
		l.pos = l.readPos
		l.readPos++
		if l.ch == '\n' {
			l.col = 0
			l.line++
			l.context = l.input[l.lastln:l.pos]
			l.lastln = l.pos
		} else {
			l.col++
		}
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\r' || l.ch == '\t' {
		l.nextChar()
	}
}

func (l *Lexer) skipBlockComment() {
	for {
		l.nextChar()
		if l.ch == '*' {
			l.nextChar()
			if l.ch == '/' {
				break
			}
			continue
		}
	}
}

func (l *Lexer) skipLineComment() {
	for {
		l.nextChar()
		if l.ch == '\n' {
			break
		}
	}
}

func lookupKeyword(ident string) (string, bool) {
	if tok, ok := keywords[ident]; ok {
		return tok, true
	}
	return ident, false
}

func isAlnum(ch byte) bool {
	return isLetter(ch) || isNumber(ch)
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' ||
		'A' <= ch && ch <= 'Z' ||
		ch == '_'
}

func isNumber(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// Err represents the an error that the lexer can return
type Err struct {
	Msg string
	Con Context
}

func (err Err) Error() string {
	return fmt.Sprintf(
		"%s: line %d, col %d",
		err.Msg,
		err.Con.Line,
		err.Con.Col,
	)
}
