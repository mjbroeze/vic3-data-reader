package files

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// DataFile is a filepath with some convenience methods.
type DataFile string

// Reader loads a DataFile, and provides the runes in that file
type Reader struct {
	file   *os.File
	reader *bufio.Reader
	pos    Position
	curr   *rune
	err    error
}

func (df DataFile) NewReader() (*Reader, error) {
	file, err := os.Open(string(df))
	if err != nil {
		return nil, err
	}

	rdr := bufio.NewReader(file)
	pos := newPosition()

	return &Reader{file: file, reader: rdr, pos: pos, curr: nil, err: nil}, nil
}

func (r *Reader) Close() error {
	if r.file == nil {
		log.Fatalf("no file to close")
	}
	return r.file.Close()
}

func (r *Reader) Peek() (rune, error) {
	rn, _, err := r.reader.ReadRune()
	if err == nil {
		unReadErr := r.reader.UnreadRune()
		if unReadErr != nil {
			return 0, fmt.Errorf("could not unread rune; reader is in an invalid state: %s", unReadErr)
		}
	}
	return rn, err
}

func (r *Reader) Next() (rune, error) {
	r.pos.advance(r.curr, r.err)
	curr, _, err := r.reader.ReadRune()
	r.curr = &curr
	r.err = err
	return curr, r.err
}

func (r *Reader) NextPosition() (Position, error) {
	ch, err := r.Peek()
	pos := r.pos
	pos.advance(&ch, err)
	return pos, err
}

type Position struct {
	pos       int
	line, col int
}

func (p *Position) Pos() int {
	return p.pos
}

func (p *Position) Line() int {
	return p.line
}

func (p *Position) Col() int {
	return p.col
}

func newPosition() Position {
	return Position{pos: -1, line: 0, col: 0}
}

func (p *Position) resetToInvalidState() {
	p.pos = -1
	p.line, p.col = 0, 0
}

// advance params are the *current* rune/error *prior* to advancing
// these will be nil on the first call to advance
func (p *Position) advance(r *rune, err error) {
	if err == nil {
		p.pos = p.pos + 1
		if r != nil && *r == '\n' {
			p.line = p.line + 1
			p.col = 1
		} else {
			p.col = p.col + 1
		}
	} else { // error (either eof or other) so pos/line/col are invalid for the rune
		p.resetToInvalidState()
	}
}

// getters

func (r *Reader) Err() error {
	return r.err
}

func (r *Reader) Pos() int {
	return r.pos.pos
}

func (r *Reader) Line() int {
	return r.pos.line
}

func (r *Reader) Col() int {
	return r.pos.col
}
