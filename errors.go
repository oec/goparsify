package goparsify

import (
	"bufio"
	"fmt"
	"strings"
)

// Error represents a parse error. These will often be set, the parser will back up a little and
// find another viable path. In general when combining errors the longest error should be returned.
type Error struct {
	pos      int
	expected string
}

// Pos is the offset into the document the error was found
func (e *Error) Pos() int { return e.pos }

// Error satisfies the golang error interface
func (e *Error) Error() string { return fmt.Sprintf("offset %d: expected %s", e.pos, e.expected) }

// UnparsedInputError is returned by Run when not all of the input was consumed. There may still be a valid result
type UnparsedInputError struct {
	Remaining string
}

// Error satisfies the golang error interface
func (e UnparsedInputError) Error() string {
	return "left unparsed: " + e.Remaining
}

// LocalError locates the error position in the input string s and returns the
// error description along with a cursor to the input.
func (e *Error) LocateError(s string) string {
	if len(s) < e.Pos() {
		return e.Error()
	}

	pos := e.Pos()

	// find the line.
	var (
		lino, prev, off int
		line, indent    []byte
	)
	byline := bufio.NewScanner(strings.NewReader(s))
	for byline.Scan() {
		lino++
		line = byline.Bytes()
		prev = off
		off += len(line) + 1
		// log.Printf("lino:%3d len:%4d prev:%4d  pos:%4d off:%4d rel:%4d line: %s",
		// 	lino, len(line), prev, pos, off, (pos - prev), string(line))
		if prev <= pos && pos <= off {
			break
		}
	}
	indent = make([]byte, len(line[:pos-prev]))
	for i, c := range indent {
		if c != '\t' {
			indent[i] = ' '
		}
	}
	off = pos - prev
	if off > 40 {
		indent = indent[off-30:]
		line = line[off-30:]
		line[0] = '.'
		line[1] = '.'
		line[2] = '.'
	}
	if len(line) > 70 {
		line = line[:70]
		line[69], line[68], line[67] = '.', '.', '.'
	}
	return fmt.Sprintf("Parsing error in line %d:\n%s\n%s^\n%v\n", lino, string(line), string(indent), e.Error())
}
