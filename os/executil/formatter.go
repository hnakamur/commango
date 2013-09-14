package executil

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"unicode/utf8"
)

func FormatCommand(cmd *exec.Cmd) (string, error) {
	var out bytes.Buffer
	for i, arg := range cmd.Args {
		if i > 0 {
			out.WriteRune(' ')
		}
		err := WriteWord(&out, arg)
		if err != nil {
			return "", err
		}
	}
	return out.String(), nil
}

func WriteWord(output *bytes.Buffer, word string) error {
	return newEscaper(word, output).run()
}

// escaper

type escaper struct {
	input      string
	start      int
	pos        int
	width      int
	output     *bytes.Buffer
	state      stateFn
	quoteStack runeStack
	err        error
}

type stateFn func(*escaper) stateFn

func newEscaper(input string, output *bytes.Buffer) *escaper {
	return &escaper{
		input:   input,
		output:  output,
	}
}

func (e *escaper) run() error {
	for e.state = scanText; e.state != nil; {
		e.state = e.state(e)
	}
	return e.err
}

func (e *escaper) readRune() (r rune, size int, err error) {
	if e.pos >= len(e.input) {
		e.width = 0
		err = io.EOF
		return
	}
	r, size = utf8.DecodeRuneInString(e.input[e.pos:])
	e.width = size
	e.pos += size
	return
}

func (e *escaper) unreadRune() error {
	if e.pos <= e.start {
		return errors.New("cannot unreadRune()")
	}
	e.pos -= e.width
	return nil
}

func (e *escaper) writeRune(r rune) (n int, err error) {
	return e.output.WriteRune(r)
}

func scanText(e *escaper) stateFn {
	r, _, err := e.readRune()
	if err == io.EOF {
		return nil
	}
	switch r {
	case ' ':
		e.writeRune('\\')
		e.writeRune(r)
		return scanText
	case '"', '\'':
		e.quoteStack.push(r)
		e.writeRune(r)
		return scanInsideQuote
	case '\\':
		e.writeRune('\\')
		r, _, err = e.readRune()
		if err == io.EOF {
			e.err = fmt.Errorf("character needed after backslash: %s.", e.input)
			return nil
		}
		e.writeRune(r)
		return scanText
	default:
		e.writeRune(r)
		return scanText
	}
}

func scanInsideQuote(e *escaper) stateFn {
	r, _, err := e.readRune()
	if err == io.EOF {
		e.err = fmt.Errorf("quote not closed: %s.", e.input)
		return nil
	}
	switch r {
	case e.quoteStack.peek():
		e.quoteStack.pop()
		e.writeRune(r)
		if e.quoteStack.isEmpty() {
			return scanText
		} else {
			return scanInsideQuote
		}
	case '\\':
		e.writeRune(r)
		r, _, err = e.readRune()
		if err == io.EOF {
			e.err = fmt.Errorf("character needed after backslash: %s.", e.input)
			return nil
		}
		e.writeRune(r)
		return scanInsideQuote
	default:
		e.writeRune(r)
		return scanInsideQuote
	}
}

// runeStack

type runeStack struct {
	runes []rune
}

func (s *runeStack) isEmpty() bool {
	return len(s.runes) == 0
}

func (s *runeStack) push(r rune) {
	s.runes = append(s.runes, r)
}

func (s *runeStack) pop() rune {
	r := s.peek()
	s.runes = s.runes[:len(s.runes)-1]
	return r
}

func (s *runeStack) peek() rune {
	return s.runes[len(s.runes)-1]
}
