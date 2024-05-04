package ast

import (
	"bytes"

	"github.com/orangeseeds/blitzbase/lang/token"
)

type Node interface {
	// the literal of the Node Token, like identifier's literal, symbols literal etc.
	TokenLiteral() string
	String() string
}
type Statement interface {
	Node
	stmtNode()
}
type Expression interface {
	Node
	exprNode()
}

type Program struct {
	Statement Statement
}

// Example ExpressionStatements:
// x = y;
// a.attr + 10 = 20;
type ExpressionStatement struct {
	Token token.Token
	// TODO: Expressions are:
	// - Prefix Expressions
	// - Infix Expressions
	Expression Expression
}

func (s *ExpressionStatement) stmtNode() {}
func (s *ExpressionStatement) String() string {
	if s.Expression != nil {
		return s.Expression.String()
	}

	return ""
}
func (s *ExpressionStatement) TokenLiteral() string { return s.Token.Literal }

type InfixExpression struct {
	Token    token.Token
	Operator string
	Left     Expression
	Right    Expression
}

func (s *InfixExpression) exprNode() {}
func (s *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(s.Left.String())
	out.WriteString(" " + s.Operator + " ")
	out.WriteString(s.Right.String())
	out.WriteString(")")

	return out.String()
}
func (s *InfixExpression) TokenLiteral() string { return s.Token.Literal }

type Identifier struct {
	Token token.Token // the token.IDENT Token
	Value string
}

func (i *Identifier) exprNode()            {}
func (i *Identifier) String() string       { return i.Value }
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

type AccessorExpresssion struct {
	Token      token.Token // starting token
	Expression Expression
}

func (i *AccessorExpresssion) exprNode()            {}
func (i *AccessorExpresssion) TokenLiteral() string { return i.Token.Literal }
func (i *AccessorExpresssion) String() string {
	var out bytes.Buffer
	successive := i.Expression.(*AccessorExpresssion)
	out.WriteString(i.TokenLiteral())
	for {
		if successive.Expression == nil {
			break
		}

		out.WriteString(".")
		out.WriteString(successive.Expression.TokenLiteral())
		successive = successive.Expression.(*AccessorExpresssion)
	}

	return out.String()
}

type DollarExpression struct {
	Token      token.Token // starting token
	Expression Expression
}

func (i *DollarExpression) exprNode()            {}
func (i *DollarExpression) TokenLiteral() string { return i.Token.Literal }
func (i *DollarExpression) String() string {
	var out bytes.Buffer
	out.WriteString(i.TokenLiteral())
	out.WriteString(i.Expression.String())
	return out.String()
}

type PrefixExpression struct {
	Token    token.Token // the prefix token, e.g., !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) exprNode()            {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (s *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(s.Operator)
	out.WriteString(s.Right.String())
	out.WriteString(")")

	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) exprNode()            {}
func (b *Boolean) String() string       { return b.Token.Literal }
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) exprNode()            {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }
