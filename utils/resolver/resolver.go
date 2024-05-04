package resolver

import (
	"errors"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/orangeseeds/blitzbase/lang/ast"
	"github.com/orangeseeds/blitzbase/lang/lexer"
	"github.com/orangeseeds/blitzbase/lang/parser"
	"github.com/orangeseeds/blitzbase/lang/token"
	model "github.com/orangeseeds/blitzbase/models"
)

type RequestInfo struct {
	Context    echo.Context
	Data       map[string]any
	Method     string
	Collection *model.Collection
}

type Resolver struct {
}

func CheckPermission(permission string, request RequestInfo) bool {
	l := lexer.New(permission)
	p := parser.New(l)
	program := p.ParseProgram()
	if program == nil {
		return false
	}
	log.Println(program.Statement.String())

	// infixExpr := stmt.(*ast.InfixExpression)
	//
	// switch infixExpr.Token {
	// case token.LESSTHAN:
	// case token.GREATERTHAN:
	// case token.EQUAL:
	// case token.LIKE:
	// case token.NOT_EQUAL:
	// case token.AT_LEAST_ONE_EQ:
	// case token.AT_LEAST_ONE_NOT_EQ:
	// case token.AND:
	// case token.OR:
	// default:
	// 	return false
	// }

	left, right, err := GetLeftAndRight(program, request)
	if err != nil {
		log.Println(err)
		return false
	}

	return left == right
}

func GetLeftAndRight(program *ast.Program, request RequestInfo) (string, string, error) {

	stmt := program.Statement.(*ast.ExpressionStatement).Expression
	log.Printf("STMT: %#v", stmt)
	switch node := stmt.(type) {
	case *ast.InfixExpression:
		left, err := getLeft(node, &request)
		if err != nil {
			return "", "", err
		}

		right, err := getRight(node, &request)
		if err != nil {
			return "", "", err
		}
		return left, right, nil
	default:
		return "", "", fmt.Errorf("NOT FOUND RIGHT!%s", "")
	}
}

func getLeft(node *ast.InfixExpression, request *RequestInfo) (string, error) {

	var leftVal string
	switch left := node.Left.(type) {
	case *ast.DollarExpression:
		// request vars check
		re, _ := left.Expression.(*ast.AccessorExpresssion)
		if re.Token.Type != token.REQUEST {
			return "", errors.New("not req after $")
		}
		acc := re.Expression.(*ast.AccessorExpresssion).Expression.(*ast.AccessorExpresssion) // auth
		switch acc.Token.Type {
		case token.AUTH:
			auth, ok := request.Context.Get("authRecord").(*model.Record) // api.CtxAuthKey
			if !ok {
				return "", errors.New("Error converting interface to AuthRecord!")
			}
			leftVal = auth.GetString(acc.Expression.TokenLiteral())
			if leftVal == "" {
				return "", fmt.Errorf("Property %s doesnot exist in auth", acc.Expression.TokenLiteral())
			}
		case token.METHOD:
			leftVal = request.Method
		default:
			return "", fmt.Errorf("Invalid property %s on Request", acc)
		}
	case *ast.AccessorExpresssion:
		// request vars check
		acc := left.Expression.(*ast.AccessorExpresssion)

		if left.Token.Type != token.COLLECTION {
			return "", errors.New("not req after $")
		}

		switch acc.Expression.TokenLiteral() {
		case "name":
			leftVal = request.Collection.Name
		default:
			fmt.Errorf("Invalid property %s on Collection", acc.Expression.TokenLiteral())
		}
	default:
		return "", fmt.Errorf("NOT FOUND LEFT inner %s!", "")
	}
	return leftVal, nil
}

func getRight(node *ast.InfixExpression, request *RequestInfo) (string, error) {

	var rightVal string
	switch right := node.Right.(type) {
	case *ast.AccessorExpresssion:
		// request vars check
		acc := right.Expression.(*ast.AccessorExpresssion)

		if right.Token.Type != token.COLLECTION {
			return "", errors.New("not req after $")
		}

		switch acc.Expression.TokenLiteral() {
		case "name":
			rightVal = request.Collection.Name
		default:
			fmt.Errorf("Invalid property %s on Collection", acc.Expression.TokenLiteral())
		}

	}
	return rightVal, nil
}
