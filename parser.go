package gqlparser

import (
	"fmt"
	"strconv"

	"github.com/matthewmueller/gqlparser/errors"
	"github.com/matthewmueller/gqlparser/lexer"
)

type parser struct {
	lexer lexer.Lexer
	err   *errors.Syntax

	peeked    bool
	peekToken lexer.Token
	peekError *errors.Syntax

	prev lexer.Token
}

func (p *parser) peek() lexer.Token {
	if p.err != nil {
		return p.prev
	}

	if !p.peeked {
		p.peekToken, p.peekError = p.lexer.ReadToken()
		p.peeked = true
	}

	return p.peekToken
}

func (p *parser) error(tok lexer.Token, format string, args ...interface{}) {
	if p.err != nil {
		return
	}
	p.err = &errors.Syntax{
		Message: fmt.Sprintf(format, args...),
		Locations: []errors.Location{
			{Line: tok.Line, Column: tok.Column},
		},
	}
}

func (p *parser) next() lexer.Token {
	if p.err != nil {
		return p.prev
	}
	if p.peeked {
		p.peeked = false
		p.prev, p.err = p.peekToken, p.peekError
	} else {
		p.prev, p.err = p.lexer.ReadToken()
	}
	return p.prev
}

func (p *parser) expectKeyword(value string) lexer.Token {
	tok := p.peek()
	if tok.Kind == lexer.Name && tok.Value == value {
		return p.next()
	}

	p.error(tok, "Expected %s, found %s", strconv.Quote(value), tok.String())
	return tok
}

func (p *parser) expect(kind lexer.Type) lexer.Token {
	tok := p.peek()
	if tok.Kind == kind {
		return p.next()
	}

	p.error(tok, "Expected %s, found %s", kind, tok.Kind.String())
	return tok
}

func (p *parser) skip(kind lexer.Type) bool {
	tok := p.peek()

	if tok.Kind != kind {
		return false
	}
	p.next()
	return true
}

func (p *parser) unexpectedError() {
	p.unexpectedToken(p.peek())
}

func (p *parser) unexpectedToken(tok lexer.Token) {
	p.error(tok, "Unexpected %s", tok.String())
}

func (p *parser) many(start lexer.Type, end lexer.Type, cb func()) {
	hasDef := p.skip(start)
	if !hasDef {
		return
	}

	for p.peek().Kind != end && p.err == nil {
		cb()
	}
	p.next()
}
