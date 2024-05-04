package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT = "IDENT" // identifier
	INT   = "INT"   // integers 1234567890
    SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	DOLLAR = "$"

	INV_COMMA = "\""

	AND = "&"
	OR  = "|"

	LESSTHAN            = "<"
	GREATERTHAN         = ">"
	EQUAL               = "="
	LIKE                = "~"
	BANG                = "!"
	QUESTION            = "?"
	ACCESSOR            = "."
	OPERATOR            = ":"
	NOT_EQUAL           = "!="
	NOT_LIKE            = "!~"
	AT_LEAST_ONE_NOT_EQ = "!?="
	AT_LEAST_ONE_EQ     = "?="

	// Keywords
	REQUEST    = "REQUEST"
	AUTH       = "AUTH"
	METHOD     = "METHOD"
	DATA       = "DATA"
	COLLECTION = "COLLECTION"

	ISSET = "ISSET"
	EACH  = "EACH"
	ONE   = "ONE"
	TRUE  = "TRUE"
	FALSE = "FALSE"
)

var operators = []TokenType{
	AND,
	OR,

	LESSTHAN,
	GREATERTHAN,
	EQUAL,
	LIKE,
	BANG,
	QUESTION,
	// ACCESSOR,
	// OPERATOR,
	AT_LEAST_ONE_EQ,
	NOT_EQUAL,
	NOT_LIKE,
	AT_LEAST_ONE_NOT_EQ,
}

var reqMacros = map[string]TokenType{
	"auth":   AUTH,
	"method": METHOD,
	"data":   DATA,
}

func Operators() []TokenType {
	return operators
}
func ReqMacros() map[string]TokenType {
	return reqMacros
}

func Keywords() map[string]TokenType {
	return keywords
}

var keywords = map[string]TokenType{
	"true":       TRUE,
	"false":      FALSE,
	"auth":       AUTH,
	"method":     METHOD,
	"data":       DATA,
	"isset":      ISSET,
	"each":       EACH,
	"one":        ONE,
	"request":    REQUEST,
	"collection": COLLECTION,
}
