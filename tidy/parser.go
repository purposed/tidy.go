package tidy

import (
	"errors"
	"fmt"
)

// Parser parses tokens.
type Parser struct {
	tokens <-chan string

	current string
	next    string
}

func newParser(tokenStream <-chan string) *Parser {
	p := &Parser{tokens: tokenStream}

	p.advance()

	return p
}

func (p *Parser) advance() {
	p.current = p.next
	p.next = <-p.tokens
}

func (p *Parser) acceptParenO() bool {
	x := p.next == "("
	if x {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) acceptParenC() bool {
	x := p.next == ")"
	if x {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) acceptFieldName() bool {
	for _, f := range Fields {
		if p.next == string(f) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) acceptOperator() bool {
	for _, f := range Operators {
		if p.next == string(f) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) acceptBoolOp() bool {
	for _, f := range BoolOperators {
		if p.next == string(f) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) eof() bool {
	return p.next == ""
}

func (p *Parser) parseExpression() (Condition, error) {
	b := BaseCondition{}
	if !p.acceptFieldName() {
		return nil, errors.New("must begin with field name")
	}
	b.Field = Field(p.current)

	if !p.acceptOperator() {
		return nil, errors.New("expecting operator after field name")
	}
	b.Operator = Operator(p.current)

	p.advance()
	b.Value = p.current
	return &b, nil
}

// ParseCondition parses into a simple condition.
func (p *Parser) ParseCondition() (Condition, error) {
	// Two valid formats.
	// EXPR = FIELD OP VAL
	// COND = EXPR
	// COND = (COND BOOL COND)
	wasParen := p.acceptParenO()

	if !wasParen {
		// Simple condition.
		return p.parseExpression()
	}

	// Bool expression.
	c := JoinedCondition{}
	leftMember, err := p.ParseCondition()
	if err != nil {
		return nil, err
	}
	c.LeftCond = leftMember

	if !p.acceptBoolOp() {
		return nil, errors.New("expected boolean operator")
	}
	c.Op = BoolOperator(p.current)

	rightMember, err := p.ParseCondition()
	if err != nil {
		return nil, err
	}
	c.RightCond = rightMember

	if !p.acceptParenC() {
		return nil, errors.New("incoherent parentheses")
	}

	return &c, nil
}

// Parse performs the full parse.
func (p *Parser) Parse() (Condition, error) {
	c, err := p.ParseCondition()
	if err != nil {
		return nil, err
	}

	if !p.eof() {
		return nil, fmt.Errorf("expected EOF, got [%s]", p.next)
	}

	return c, nil
}
