package parser

import (
	"testing"

	"github.com/orangeseeds/blitzbase/lang/ast"
	"github.com/orangeseeds/blitzbase/lang/lexer"
	"github.com/orangeseeds/blitzbase/lang/token"
)

func TestExpressionStatement(t *testing.T) {
	cases := []struct {
		input string
	}{
		// {"apple = ball;"},
		// {"collection.name != 2;"},
		// {"$request.auth.name = collection.name;"},
		{"$request.data.is_valid_date & collection.exists | collection.id != 10;"},
		// {"$request.method = ;"},
	}

	for _, tc := range cases {
		l := lexer.New(tc.input)
		p := New(l)
		program := p.ParseProgram()
		if program == nil {
			t.Error("Some error occured!")
			t.FailNow()
		}
		checkParseErrors(t, p)

		t.Log(program.Statement.String())

		// checkSomeOther(t, program)
	}
}

func checkSomeOther(t *testing.T, program *ast.Program) {
	req := struct {
		method string
		auth   struct {
			id    string
			name  string
			field string
		}
		data struct{}
	}{
		method: "GET",
		auth: struct {
			id    string
			name  string
			field string
		}{
			id:    "1",
			name:  "auth_collection",
			field: "field",
		},
		data: struct{}{},
	}

	col := struct{ name string }{
		name: "collection_name",
	}

	var leftVal string
	var rightVal string

	stmt := program.Statement.(*ast.ExpressionStatement).Expression

	switch node := stmt.(type) {
	case *ast.InfixExpression:
		switch left := node.Left.(type) {
		case *ast.DollarExpression:
			// request vars check
			re, _ := left.Expression.(*ast.AccessorExpresssion)
			if re.Token.Type != token.REQUEST {
				t.Error("not req after $")
				t.FailNow()
			}

			acc := re.Expression.(*ast.AccessorExpresssion).Expression.(*ast.AccessorExpresssion) // auth
			switch acc.Token.Type {
			case token.AUTH:
				switch acc.Expression.TokenLiteral() {
				case "id":
					leftVal = req.auth.id
				case "name":
					leftVal = req.auth.name
				case "field":
					leftVal = req.auth.field
				default:
					t.Errorf("Property %s doesnot exist in auth", acc.Expression.TokenLiteral())
				}
			case token.METHOD:
				leftVal = req.method
			default:
				t.Errorf("Invalid property %s on Request", acc)
			}
		default:
			t.Error("NOT FOUND LEFT!")
		}

		switch right := node.Right.(type) {
		case *ast.AccessorExpresssion:
			// request vars check
			acc := right.Expression.(*ast.AccessorExpresssion)

			if right.Token.Type != token.COLLECTION {
				t.Error("error not collection")
				t.FailNow()
			}

			switch acc.Expression.TokenLiteral() {
			case "name":
				rightVal = col.name
			default:
				t.Errorf("Invalid property %s on Collection", acc.Expression.TokenLiteral())
			}
		default:
			t.Error("NOT FOUND RIGHT!")
		}
	}

	t.Log("left", leftVal)
	t.Log("right", rightVal)

}

func checkParseErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("Parser error %q", msg)
	}

	t.FailNow()
}
