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
	buf := ""

	return
}

func (l *Lexer) nextToken() Token {
	l.skipWhitespace()
	switch {
	case l.ch == '=':
		l.nextChar()
		return newToken(ASSIGN, ASSIGN, l.line, l.col)
	case l.ch == '+':
		l.nextChar()
		return newToken(ADD, ADD, l.line, l.col)
	case l.ch == '-':
		l.nextChar()
		return newToken(SUB, SUB, l.line, l.col)
	case l.ch == '*':
		l.nextChar()
		return newToken(MUL, MUL, l.line, l.col)
	case l.ch == '/':
		l.nextChar()
		return newToken(DIV, DIV, l.line, l.col)
	case l.ch == '>':
		switch l.peek() {
		case '>':
			return newToken(BSR, BSR, l.line, l.col)
		case '=':
			return newToken(GE, GE, l.line, l.col)
		default:
			return newToken(ILLEGAL, ILLEGAL, l.line, l.col)
		}
		return newToken(DIV, DIV, l.line, l.col)
	case l.ch == ',':
		l.nextChar()
		return newToken(COMMA, COMMA, l.line, l.col)
	case l.ch == ';':
		l.nextChar()
		return newToken(SEMICOL, SEMICOL, l.line, l.col)
	case l.ch == '(':
		l.nextChar()
		return newToken(LPAREN, LPAREN, l.line, l.col)
	case isLetter(l.ch):
		return l.readIdent()
	default:
		return newToken(ILLEGAL, ILLEGAL, l.line, l.col)
	}
}

func (l *Lexer) readIdent() Token {
	position := l.pos
	for isLetter(l.ch) {
		l.nextChar()
	}
	token := l.input[position:l.pos]
	if tok := lookupKeyword(token); tok != IDENT {
		return newToken(tok, tok, l.line, l.col)
	} else {
		return newToken(IDENT, token, l.line, l.col)
	}
}

func (l *Lexer) peek() byte {
	return l.input[l.readPos]
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

func lookupKeyword(ident string) string {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
