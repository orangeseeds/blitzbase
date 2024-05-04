package lexer

import "github.com/orangeseeds/blitzbase/lang/token"

type Lexer struct {
	input   string
	pos     int
	readPos int
	char    byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input: input,
	}
	l.readChar() // first read to fill pos,readPos and char
	return l
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos += 1
}

// Literal is string instead of byte because we donot have any single char tok
func NewToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()
	switch l.char {
	case ';':
		tok = NewToken(token.SEMICOLON, l.char)
	case '$':
		tok = NewToken(token.DOLLAR, l.char)
	case '&':
		tok = NewToken(token.AND, l.char)
	case '|':
		tok = NewToken(token.OR, l.char)
	case '<':
		tok = NewToken(token.GREATERTHAN, l.char)
	case '>':
		tok = NewToken(token.LESSTHAN, l.char)
	case '"':
		tok = NewToken(token.INV_COMMA, l.char)
	case '=':
		tok = NewToken(token.EQUAL, l.char)
	case '.':
		tok = NewToken(token.ACCESSOR, l.char)
	case ':':
		tok = NewToken(token.OPERATOR, l.char)
	case '!':
		switch l.peekChar() {
		case '=':
			bangChar := l.char
			l.readChar() // load '=' into l.char
			tok = token.Token{
				Type:    token.NOT_EQUAL,
				Literal: string(bangChar) + string(l.char),
			}
		case '~':
			bangChar := l.char
			l.readChar()
			tok = token.Token{
				Type:    token.NOT_LIKE,
				Literal: string(bangChar) + string(l.char),
			}
		case '?':
			bangChar := l.char
			l.readChar() // reads the '?' char
			if l.peekChar() == '=' {
				questionChar := l.char
				l.readChar()
				tok = token.Token{
					Type:    token.AT_LEAST_ONE_NOT_EQ,
					Literal: string(bangChar) + string(questionChar) + string(l.char),
				}
			} else {
				tok = token.Token{
					Type:    token.ILLEGAL,
					Literal: string(bangChar) + string(l.char),
				}
			}
		default:
			tok = NewToken(token.BANG, l.char)
		}
	case '~':
		tok = NewToken(token.LIKE, l.char)
	case '?':
		if l.peekChar() == '=' {
			questionChar := l.char
			l.readChar()
			tok = token.Token{
				Type:    token.AT_LEAST_ONE_EQ,
				Literal: string(questionChar) + string(l.char),
			}
		} else {
			tok = NewToken(token.QUESTION, l.char)
		}
	case 0:
		tok = token.Token{
			Type:    token.EOF,
			Literal: "",
		}
	default:
		if isLetter(l.char) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.char) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = NewToken(token.ILLEGAL, l.char)
		}
	}
	l.readChar()

	return tok
}

func (l *Lexer) readIdentifier() string {
	postion := l.pos
	for isLetter(l.char) {
		l.readChar()
	}

	return l.input[postion:l.pos]
}

func (l *Lexer) readNumber() string {
	position := l.pos
	for isDigit(l.char) {
		l.readChar()
	}
	return l.input[position:l.pos]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) skipWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		l.readChar()
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPos]
	}
}

func LookupIdent(ident string) token.TokenType {
	if tok, ok := token.Keywords()[ident]; ok {
		return tok
	}
	return token.IDENT
}
