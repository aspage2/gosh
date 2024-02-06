package parse

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

type TokenType int

const (
	IDENT TokenType = iota
	VARIABLE
	WHITESPACE
	SEMICOLON
	STRING
	EOF
)

const (
	eof rune = '\ufffd'
)

type Token struct {
	Type  TokenType
	Value string
}


type Predicate func(rune) bool

func IsEOF(c rune) bool {
	return c == eof
}

func Not(b bool) bool {
	return !b
}

func Any(ps... Predicate) Predicate {
	return func(c rune) bool {
		for _, p := range ps {
			if p(c) { return true }
		}
		return false
	}
}

func Eq(c... rune) Predicate {
	return func(chr rune) bool {
		for _, ch := range c {
			if chr == c {
				return true
			}
		}
		return false
	}
}

func (t Token) String() string {
	var (
		hasValue bool
		typeName string
	)
	switch t.Type {
	case IDENT:
		hasValue = true
		typeName = "IDENT"
	case WHITESPACE:
		typeName = "WS"
	case VARIABLE:
		hasValue = true
		typeName = "VAR"
	case SEMICOLON:
		typeName = "SEMI"
	case STRING:
		typeName = "STRING"
	case EOF:
		typeName = "EOF"
	default:
		typeName = "???"
		hasValue = true
	}
	if hasValue {
		return fmt.Sprintf("%s(%s)", typeName, t.Value)
	} else {
		return typeName
	}
}

type Lexer struct {
	start   int
	pos     int
	width   int
	tokens  chan (Token)
	program string
}

func (l *Lexer) Next() rune {
	if l.pos >= len(l.program) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.program[l.pos:])
	l.pos += w
	l.width = w
	return r
}

func (l *Lexer) TakeWhile(p Predicate) int {
	n := 0
	char := l.Next()
	for char != eof && p {
		char = l.Next()
		n += 1
	}
	l.Backup()
	return n
}

func (l *Lexer) TakeUntil(p Predicate) int {
	return l.TakeWhile(Not(p))
}

func (l *Lexer) Backup() {
	l.pos -= l.width
	l.width = 0
}

func (l *Lexer) Peek() rune {
	r := l.Next()
	l.Backup()
	return r
}

func (l *Lexer) Ignore() {
	l.start = l.pos
	l.width = 0
}

func (l *Lexer) Emit(typ TokenType) {
	if l.start == l.pos {
		return
	}
	s := l.program[l.start:l.pos]
	l.start = l.pos
	l.tokens <- Token{typ, s}
}

func (l *Lexer) Run() {
	for s := Base; s != nil; s = s(l) {

	}
	close(l.tokens)
}

func Lex(program string) (*Lexer, chan Token) {
	l := &Lexer{
		program: program,
		tokens:  make(chan Token),
	}
	go l.Run()
	return l, l.tokens
}

// ------------------------------------------------

type State func(*Lexer) State

func Base(l *Lexer) State {
	char := l.Next()
	var (
		next State
	)
loop:
	for char != eof {
		switch char {
		case '"':
			l.Ignore()
			next = _String
			break loop
		case '$':
			l.Ignore()
			next = Variable
			break loop
		case ' ', '\n', '\t', '\r':
			next = WhiteSpace
			break loop
		}
		char = l.Next()
	}

	l.Backup()
	l.Emit(IDENT)

	if char == eof {
		return nil
	}
	return next
}

func WhiteSpace(l *Lexer) State {
	l.TakeWhile(unicode.IsSpace)
	l.Emit(WHITESPACE)
	if l.Peek() == eof {
		return nil
	}
	return Base
}

func isVarCharacter(c rune) bool {
	return c < unicode.MaxASCII && (unicode.IsLetter(c) || unicode.IsNumber(c)) || c == '_'
}

func Variable(l *Lexer) State {
	l.TakeWhile(isVarCharacter)
	l.Emit(VARIABLE)
	if l.Peek() == eof {
		return nil
	}
	return Base
}

func String(l *Lexer) State {
	l.TakeUntil(Eq('"'))
	if l.Peek() == eof {
		panic("string is missing a closing quote")
	}
	l.Backup()
	l.Emit(STRING)
	return Base
}
