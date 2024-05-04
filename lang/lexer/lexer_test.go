package lexer

import (
	"testing"

	"github.com/orangeseeds/blitzbase/lang/token"
)

// go test -run <function> <path>
// go test <path>

func TestNewToken(t *testing.T) {
	input := `$request.data.name = "apple" | $collection.auth = $request.auth.id data:each = 10`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.DOLLAR, "$"},
		{token.REQUEST, "request"},
		{token.ACCESSOR, "."},
		{token.DATA, "data"},
		{token.ACCESSOR, "."},
		{token.IDENT, "name"},
		{token.EQUAL, "="},
		{token.INV_COMMA, "\""},
		{token.IDENT, "apple"},
		{token.INV_COMMA, "\""},
		{token.OR, "|"},
		{token.DOLLAR, "$"},
		{token.COLLECTION, "collection"},
		{token.ACCESSOR, "."},
		{token.AUTH, "auth"},
		{token.EQUAL, "="},
		{token.DOLLAR, "$"},
		{token.REQUEST, "request"},
		{token.ACCESSOR, "."},
		{token.AUTH, "auth"},
		{token.ACCESSOR, "."},
		{token.IDENT, "id"},
		{token.DATA, "data"},
		{token.OPERATOR, ":"},
		{token.EACH, "each"},
		{token.EQUAL, "="},
		{token.INT, "10"},
	}

	l := NewLexer(input)
	for i, tt := range tests {
		tok := l.NextToken()

		// t.Logf("%v\n", tok)
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expeced %q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
