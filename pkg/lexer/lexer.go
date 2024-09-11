package lexer

import (
	"errors"
	"unicode"
)

type Lexer struct {
	content []rune
}

func NewLexer(content string) *Lexer {
	return &Lexer{
		content: []rune(content),
	}
}

func (l *Lexer) Tokens() []string {
	tokens := []string{}
	for {
		token, err := l.NextToken()
		if err != nil {
			break
		}
		tokens = append(tokens, token)
	}
	return tokens
}

func (l *Lexer) isEmpty() bool {
	if len(l.content) < 1 {
		return true
	}
	return false
}

func (l *Lexer) NextToken() (string, error) {
	l.trimLeft()
	if l.isEmpty() {
		return "", errors.New("EOF")
	}

	if unicode.IsDigit(l.content[0]) {
		return l.chopDigits(), nil
	}

	if unicode.IsLetter(l.content[0]) {
		return l.chopAlpha(), nil
	}

	return l.chop(1), nil
}

func (l *Lexer) chop(n int) string {
	token := l.content[0:n]
	l.content = l.content[n:]
	return string(token)
}

func (l *Lexer) chopDigits() string {
	n := 0
	for n < len(l.content) && unicode.IsDigit(l.content[n]) {
		n++
	}
	return l.chop(n)
}

func (l *Lexer) chopAlpha() string {
	n := 0
	for n < len(l.content) && unicode.IsLetter(l.content[n]) {
		n++
	}
	return l.chop(n)
}

func (l *Lexer) trimLeft() {
	for !l.isEmpty() && unicode.IsSpace(l.content[0]) {
		l.content = l.content[1:]
	}
}
