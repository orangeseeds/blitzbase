package parser

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/orangeseeds/blitzbase/lang/ast"
	"github.com/orangeseeds/blitzbase/lang/lexer"
	"github.com/orangeseeds/blitzbase/lang/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	lexer     *lexer.Lexer
	errors    []string
	currToken token.Token
	peekToken token.Token
}

func New(lexer *lexer.Lexer) *Parser {
	p := Parser{lexer: lexer}
	p.nextToken()
	p.nextToken()
	return &p
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	for p.currToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statement = stmt
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	return p.parseExprStatement()
}

func (p *Parser) parseExprStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currToken}
	stmt.Expression = p.parseExpression()

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpression() ast.Expression {

	var leftExpr ast.Expression
	switch p.currToken.Type {
	case token.INT:
		leftExpr = p.parseIntegerLiteral()
	case token.IDENT:
		if p.peekTokenIs(token.ACCESSOR) {
			leftExpr = p.parseAccessorExpr()
		} else {
			leftExpr = p.parseIdentifier()
		}
	case token.COLLECTION:
		if p.peekTokenIs(token.ACCESSOR) {
			leftExpr = p.parseAccessorExpr()
		} else {
			msg := fmt.Sprintf("cannot use %s bare, needs to access property.", p.currToken.Literal)
			p.errors = append(p.errors, msg)
			leftExpr = nil
		}
	case token.DOLLAR:
		leftExpr = p.parseDollarReqExpr()
	case token.LPAREN:
		leftExpr = p.parseGroupExpr()
	case token.BANG:
		leftExpr = p.parsePrefixExpr()
	case token.QUESTION:
		leftExpr = p.parsePrefixExpr()
	case token.TRUE:
		leftExpr = p.parseBoolean()
	case token.FALSE:
		leftExpr = p.parseBoolean()
	default:
		leftExpr = nil
	}

	if leftExpr == nil {
		return nil
	}

	for !p.peekTokenIs(token.SEMICOLON) {

		var infix ast.Expression

		if slices.Contains(token.Operators(), p.peekToken.Type) {
			p.nextToken()
			infix = p.parseInfixExpr(leftExpr)
		} else {
			infix = nil
		}

		// case token.OR:
		// case token.LESSTHAN:
		// case token.GREATERTHAN:
		// case token.EQUAL:
		// case token.LIKE:
		// case token.BANG:
		// case token.QUESTION:
		// case token.AT_LEAST_ONE_EQ:
		// case token.NOT_EQUAL:
		// case token.NOT_LIKE:
		// case token.AT_LEAST_ONE_NOT_EQ:

		if infix == nil {
			return leftExpr
		}

		return infix
	}

	return leftExpr
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}
}

func (p *Parser) parseAccessorExpr() ast.Expression {
	root := &ast.AccessorExpresssion{
		Token: p.currToken,
	}
	// log.Println(p.currToken.Literal)
	// p.nextToken()
	// log.Println(p.currToken.Literal)

	successive := root
	for !p.curTokenIs(token.EOF) {

		// log.Println(p.currToken)

		// log.Println(p.currToken.Literal)
		// do something

		if p.currToken.Type == token.ACCESSOR {
			p.nextToken()
		} else {
			// if root.Token.Type == token.REQUEST {
			// 	switch p.currToken.Type {
			// 	case token.AUTH:
			// 	case token.METHOD:
			// 		if !slices.Contains(token.Operators(), p.peekToken.Type) {
			// 			msg := fmt.Sprintf("could add accessor %s to %s", p.peekToken.Literal, p.currToken.Literal)
			// 			p.errors = append(p.errors, msg)
			// 			return nil
			// 		}
			// 	case token.DATA:
			// 	default:
			// 		msg := fmt.Sprintf("Expected next token to be %s, %s or %s got %s instead",
			// 			token.AUTH, token.DATA, token.METHOD, p.currToken.Literal)
			// 		p.errors = append(p.errors, msg)
			// 		return nil
			// 	}
			//
			// }
			newExpr := &ast.AccessorExpresssion{
				Token: p.currToken,
			}
			successive.Expression = newExpr
			successive = newExpr
			p.nextToken()
		}

		if slices.Contains(token.Operators(), p.peekToken.Type) || p.peekTokenIs(token.SEMICOLON) || p.peekTokenIs(token.ILLEGAL) {
			newExpr := &ast.AccessorExpresssion{
				Token: p.currToken,
			}
			successive.Expression = newExpr
			successive = newExpr
			// log.Println(p.currToken.Literal)
			break
		}
	}

	return root
}
func (p *Parser) parseDollarReqExpr() ast.Expression {
	expr := ast.DollarExpression{
		Token: p.currToken,
	}
	if !p.expectPeek(token.REQUEST) {
		return nil
	}

	accessorExpr := p.parseAccessorExpr()
	if accessorExpr == nil {
		return nil
	}
	expr.Expression = accessorExpr

	return &expr
}

func (p *Parser) parsePrefixExpr() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression()

	return expression
}
func (p *Parser) parseInfixExpr(left ast.Expression) ast.Expression {

	expression := &ast.InfixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
		Left:     left,
	}

	p.nextToken()
	expression.Right = p.parseExpression()

	return expression
}
func (p *Parser) parseGroupExpr() ast.Expression {

	p.nextToken()

	exp := p.parseExpression()

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.currToken, Value: p.curTokenIs(token.TRUE)}
}
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.currToken}

	value, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

// helper functions

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.currToken.Type == t
}
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("Expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}
