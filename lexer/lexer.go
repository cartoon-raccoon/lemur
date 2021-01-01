package lexer

/*------------Lexer------------*/

type Lexer struct {
	input   string
	line    int
	col     int
	pos     int
	readPos int
	ch      byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.nextChar()
	return l
}

func (l *Lexer) Tokenize() (tokens []Token) {
	tokens = make([]Token, 0)
	tok := l.nextToken()
	for !tok.isEOF() {
		tok = l.nextToken()
		tokens = append(tokens, tok)
	}
	return
}

func (l *Lexer) nextToken() Token {
	l.skipWhitespace()

	switch {

	case l.ch == '=':
		l.nextChar()
		if l.ch == '=' {
			l.nextChar()
			return newToken(EQ, EQ, l.line, l.col)
		}
		return newToken(ASSIGN, ASSIGN, l.line, l.col)

	case l.ch == '+':
		l.nextChar()
		if l.ch == '=' {
			l.nextChar()
			return newToken(ADDASSIGN, ADDASSIGN, l.line, l.col)
		}
		return newToken(ADD, ADD, l.line, l.col)

	case l.ch == '-':
		l.nextChar()
		if l.ch == '=' {
			l.nextChar()
			return newToken(SUBASSIGN, SUBASSIGN, l.line, l.col)
		}
		return newToken(SUB, SUB, l.line, l.col)

	case l.ch == '*':
		l.nextChar()
		if l.ch == '=' {
			l.nextChar()
			return newToken(MULASSIGN, MULASSIGN, l.line, l.col)
		}
		return newToken(MUL, MUL, l.line, l.col)

	case l.ch == '/':
		l.nextChar()
		if l.ch == '=' {
			l.nextChar()
			return newToken(DIVASSIGN, DIVASSIGN, l.line, l.col)
		}
		return newToken(DIV, DIV, l.line, l.col)

	case l.ch == '>':
		l.nextChar()

		switch l.ch {
		case '>':
			l.nextChar()
			return newToken(BSR, BSR, l.line, l.col)
		case '=':
			l.nextChar()
			return newToken(GE, GE, l.line, l.col)
		default:
			return newToken(GT, GT, l.line, l.col)
		}

	case l.ch == '<':
		l.nextChar()
		switch l.ch {
		case '<':
			l.nextChar()
			return newToken(BSL, BSL, l.line, l.col)
		case '=':
			l.nextChar()
			return newToken(LE, LE, l.line, l.col)
		default:
			return newToken(LT, LT, l.line, l.col)
		}

	case l.ch == '&':
		l.nextChar()
		switch l.ch {
		case '&':
			l.nextChar()
			return newToken(LAND, LAND, l.line, l.col)
		case '=':
			l.nextChar()
			return newToken(BWAASSIGN, BWAASSIGN, l.line, l.col)
		default:
			return newToken(BWAND, BWAND, l.line, l.col)
		}

	case l.ch == '|':
		l.nextChar()
		switch l.ch {
		case '|':
			l.nextChar()
			return newToken(LOR, LOR, l.line, l.col)
		case '=':
			l.nextChar()
			return newToken(BWOASSIGN, BWOASSIGN, l.line, l.col)
		default:
			return newToken(BWOR, BWOR, l.line, l.col)
		}

	case l.ch == '~':
		l.nextChar()
		if l.ch == '=' {
			l.nextChar()
			return newToken(BWNASSIGN, BWNASSIGN, l.line, l.col)
		}
		return newToken(BWNOT, BWNOT, l.line, l.col)

	case l.ch == ',':
		l.nextChar()
		return newToken(COMMA, COMMA, l.line, l.col)

	case l.ch == ';':
		l.nextChar()
		return newToken(SEMICOL, SEMICOL, l.line, l.col)

	case l.ch == '(':
		l.nextChar()
		return newToken(LPAREN, LPAREN, l.line, l.col)

	case l.ch == ')':
		l.nextChar()
		return newToken(RPAREN, RPAREN, l.line, l.col)

	case l.ch == '[':
		l.nextChar()
		return newToken(LSBRKT, LSBRKT, l.line, l.col)

	case l.ch == ']':
		l.nextChar()
		return newToken(RSBRKT, RSBRKT, l.line, l.col)

	case l.ch == '{':
		l.nextChar()
		return newToken(LBRACE, LBRACE, l.line, l.col)

	case l.ch == '}':
		l.nextChar()
		return newToken(RBRACE, RBRACE, l.line, l.col)

	case isAlnum(l.ch):
		return l.readIdent()

	case l.ch == 0:
		return newToken(EOF, EOF, l.line, l.col)

	default:
		return newToken(ILLEGAL, ILLEGAL, l.line, l.col)
	}
}

func (l *Lexer) readIdent() Token {
	position := l.pos
	for isAlnum(l.ch) {
		l.nextChar()
	}
	token := l.input[position:l.pos]
	if tok, is_kw := lookupKeyword(token); is_kw {
		return newToken(tok, tok, l.line, l.col)
	} else {
		return newToken(IDENT, token, l.line, l.col)
	}
}

func (l *Lexer) nextChar() {
	//have reached EOF
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
		l.pos = l.readPos
		l.readPos++
		if l.ch == '\n' {
			l.col = 0
			l.line++
		} else {
			l.col++
		}
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\r' || l.ch == '\t' {
		if l.ch == '\n' {
			l.col = 0
			l.line++
		} else {
			l.col++
		}
		l.nextChar()
	}
}

func lookupKeyword(ident string) (string, bool) {
	if tok, ok := keywords[ident]; ok {
		return tok, true
	}
	return ident, false
}

func isAlnum(ch byte) bool {
	return 'a' <= ch && ch <= 'z' ||
		'A' <= ch && ch <= 'Z' ||
		'0' <= ch && ch <= '9' ||
		ch == '_'
}
