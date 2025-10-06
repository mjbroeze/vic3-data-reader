package files

import (
	"io"
	"log"
	"strings"
	"testing"
)

const (
	DoesNotExist   DataFile = "testdata/DOES-NOT-EXIST.txt"
	Empty          DataFile = "testdata/empty.txt"
	SingleA        DataFile = `testdata/single-a.txt`
	NewlineSingleA DataFile = "testdata/newline-single-a.txt"
	SmokeSample    DataFile = "testdata/00_goods.txt"
)

func TestNewReader_noErrOnExistingFile(t *testing.T) {
	_, err := Empty.NewReader()
	if err != nil {
		t.Errorf("NewReader returned an error for existing file: %v", err)
	}
}

func TestNewReader_errOnMissingFile(t *testing.T) {
	_, err := DoesNotExist.NewReader()
	log.Printf("err: %v", err)
	if err == nil {
		t.Errorf("NewReader did not return an error for missing file")
	}
}

func TestNewReader_initPosIsSentinel(t *testing.T) {
	reader, err := Empty.NewReader()
	if err != nil {
		t.Errorf("could not open file")
	}

	pos := reader.Pos()
	if pos != -1 {
		t.Errorf("reader should be intiialized with pos of sentinel value -1; actual: %d", pos)
	}
}

func TestNewReader_initLineIsSentinel(t *testing.T) {
	reader, err := Empty.NewReader()
	if err != nil {
		t.Errorf("could not open file")
	}

	line := reader.Line()
	if line != 0 {
		t.Errorf("reader should be intiialized with line of sentinel value 0; actual: %d", line)
	}
}

func TestNewReader_initColIsSentinel(t *testing.T) {
	reader, err := Empty.NewReader()
	if err != nil {
		t.Errorf("could not open file")
	}

	col := reader.Col()
	if col != 0 {
		t.Errorf("reader should be intiialized with col of sentinel value 0; actual: %d", col)
	}
}

func TestClose(t *testing.T) {
	reader, err := Empty.NewReader()
	if err != nil {
		t.Errorf("could not open file")
	}

	err = reader.Close()
	if err != nil {
		t.Errorf("could not close reader: %v", err)
	}

	fileCloseErr := reader.file.Close()
	if fileCloseErr != nil {
		if !strings.HasSuffix(fileCloseErr.Error(), "file already closed") {
			t.Errorf("reader.Close failed with unexpected error %v", fileCloseErr)
		}
	} else {
		t.Errorf("reader.Close failed")
	}
}

func TestClose_closeTwiceReturnsError(t *testing.T) {
	reader, err := Empty.NewReader()
	if err != nil {
		t.Errorf("could not open file")
	}

	err = reader.Close()
	if err != nil {
		t.Errorf("could not close reader: %v", err)
	}

	err = reader.Close()
	if err != nil {
		if !strings.HasSuffix(err.Error(), "file already closed") {
			t.Errorf("second call to reader.Close failed with unexpected error %v", err)
		}
	} else {
		t.Errorf("second call to reader.Close did not result in an error")
	}
}

func TestPeek(t *testing.T) {
	reader, err := SingleA.NewReader()
	if err != nil {
		t.Errorf("could not open file")
	}

	ch, err := reader.Peek()
	if ch != 'a' {
		t.Errorf("peeked value should be 'a'; actual: %v", ch)
	} else if err != nil {
		t.Errorf("Peek returned unexpected error: %v", err)
	}
}

func TestPeek_posIsNotChanged(t *testing.T) {
	reader, err := SingleA.NewReader()
	if err != nil {
		t.Errorf("could not open file")
	}

	initPos := reader.Pos()
	_, err = reader.Peek()
	if err != nil {
		t.Errorf("Peek returned unexpected error: %v", err)
	}

	currPos := reader.Pos()
	if initPos != currPos {
		t.Errorf("Peek changed pos from %d to %d", initPos, currPos)
	}
}

func TestPeek_colIsNotChanged(t *testing.T) {
	reader, err := SingleA.NewReader()
	if err != nil {
		t.Errorf("could not open file")
	}

	initCol := reader.Col()

	_, err = reader.Peek()
	if err != nil {
		t.Errorf("Peek returned unexpected error: %v", err)
	}

	currCol := reader.Col()
	if initCol != currCol {
		t.Errorf("Peek changed col from %d to %d", initCol, currCol)
	}
}

func TestPeek_lineIsNotChanged(t *testing.T) {
	reader, err := NewlineSingleA.NewReader()
	if err != nil {
		t.Errorf("could not open file")
	}

	/// advance to first char (newline)
	ch, err := reader.Next()
	if err != nil {
		t.Errorf("Next returned unexpected error: %v", err)
	} else if ch != '\n' {
		t.Errorf("first rune should be '\\n'; actual: '%v'", ch)
	}

	// now peek will return the char on the next line
	initLine := reader.Line()

	_, err = reader.Peek()
	if err != nil {
		t.Errorf("Peek returned unexpected error: %v", err)
	}

	currLine := reader.Line()
	if initLine != currLine {
		t.Errorf("Peek changed line from %d to %d", initLine, currLine)
	}
}

func TestPeek_eofReturnsEOFErr(t *testing.T) {
	reader, err := Empty.NewReader()
	if err != nil {
		t.Errorf("could not open file")
	}

	_, err = reader.Peek()
	if err == nil {
		t.Errorf("Peek return nil error for EOF")
	} else if err != io.EOF {
		t.Errorf("Peek did not return EOF error: %v", err)
	}
}

func TestNext(t *testing.T) {
	reader, err := SingleA.NewReader()
	if err != nil {
		t.Errorf("could not open file")
	}

	ch, err := reader.Next()
	if ch != 'a' {
		t.Errorf("Next value should be 'a'; actual: %v", ch)
	} else if err != nil {
		t.Errorf("Next returned unexpected error: %v", err)
	}
}

func TestPeek_posIsAdvanced(t *testing.T) {
	reader, err := SingleA.NewReader()
	if err != nil {
		t.Errorf("could not open file")
	}

	initPos := reader.Pos()
	if err != nil {
		t.Errorf("Pos returned unexpected error: %v", err)
	}

	_, err = reader.Next()
	if err != nil {
		t.Errorf("Next returned unexpected error: %v", err)
	}

	currPos := reader.Pos()
	if err != nil {
		t.Errorf("Pos returned unexpected error: %v", err)
	}

	expectedPos := initPos + 1
	if expectedPos != currPos {
		t.Errorf("Next changed pos from %d to %d, %d expected", initPos, currPos, expectedPos)
	}
}

func TestNext_colIsAdvanced(t *testing.T) {
	reader, err := SingleA.NewReader()
	if err != nil {
		t.Errorf("could not open file")
	}

	initCol := reader.Col()
	if err != nil {
		t.Errorf("Col returned unexpected error: %v", err)
	}

	_, err = reader.Next()
	if err != nil {
		t.Errorf("Next returned unexpected error: %v", err)
	}

	currCol := reader.Col()
	if err != nil {
		t.Errorf("Col returned unexpected error: %v", err)
	}

	expectedCol := initCol + 1
	if expectedCol != currCol {
		t.Errorf("Next changed col from %d to %d, %d expected", initCol, currCol, expectedCol)
	}
}

func TestNext_lineIsAdvanced(t *testing.T) {
	reader, err := NewlineSingleA.NewReader()
	if err != nil {
		t.Errorf("could not open file")
	}

	/// advance to first char (newline)
	ch, err := reader.Next()
	if err != nil {
		t.Errorf("Next returned unexpected error: %v", err)
	} else if ch != '\n' {
		t.Errorf("first rune should be '\\n'; actual: '%v'", ch)
	}

	// now next will return the char on the next line
	initLine := reader.Line()
	if err != nil {
		t.Errorf("Line returned unexpected error: %v", err)
	}

	ch, err = reader.Next()
	if err != nil {
		t.Errorf("Next returned unexpected error: %v", err)
	}

	currLine := reader.Line()
	if err != nil {
		t.Errorf("Line returned unexpected error: %v", err)
	}

	expectedLine := initLine + 1
	if expectedLine != currLine {
		t.Errorf("Next changed line from %d to %d, %d expected", initLine, currLine, expectedLine)
	}
}

func TestNext_eofReturnsEOFErr(t *testing.T) {
	reader, err := Empty.NewReader()
	if err != nil {
		t.Errorf("could not open file")
	}

	_, err = reader.Next()
	if err == nil {
		t.Errorf("Next returned nil error for EOF")
	} else if err != io.EOF {
		t.Errorf("Next did not return EOF error: %v", err)
	}
}

func TestNextPosition_pos(t *testing.T) {
	reader, err := SmokeSample.NewReader()
	if err != nil {
		t.Errorf("could not open file")
	}

	nextPos, err := reader.NextPosition()
	if err != nil {
		t.Errorf("NextPosition returned unexpected error: %v", err)
	}

	_, err = reader.Next()
	if err != nil {
		t.Errorf("Next returned unexpected error: %v", err)
	}

	if reader.Pos() != nextPos.Pos() {
		t.Errorf("NextPosition expected pos %d, actual: %d", reader.Pos(), nextPos.Pos())
	}
}

func TestNextPosition_line(t *testing.T) {
	reader, err := SmokeSample.NewReader()
	if err != nil {
		t.Errorf("could not open file")
	}

	nextPos, err := reader.NextPosition()
	if err != nil {
		t.Errorf("NextPosition returned unexpected error: %v", err)
	}

	_, err = reader.Next()
	if err != nil {
		t.Errorf("Next returned unexpected error: %v", err)
	}

	if reader.Line() != nextPos.Line() {
		t.Errorf("NextPosition expected line %d, actual: %d", reader.Line(), nextPos.Line())
	}
}

func TestNextPosition_col(t *testing.T) {
	reader, err := SmokeSample.NewReader()
	if err != nil {
		t.Errorf("could not open file")
	}

	nextPos, err := reader.NextPosition()
	if err != nil {
		t.Errorf("NextPosition returned unexpected error: %v", err)
	}

	_, err = reader.Next()
	if err != nil {
		t.Errorf("Next returned unexpected error: %v", err)
	}

	if reader.Col() != nextPos.col {
		t.Errorf("NextPosition expected col %d, actual: %d", reader.Col(), nextPos.Col())
	}
}

func TestSmoke(t *testing.T) {
	reader, err := SmokeSample.NewReader()
	if err != nil {
		t.Errorf("could not open file")
	}

	var contents []rune

	var ch rune
	var nextErr error
	for nextErr == nil {
		ch, nextErr = reader.Next()
		if nextErr != nil && nextErr != io.EOF {
			t.Errorf("Next returned unexpected error: %v", err)
		}
		contents = append(contents, ch)
	}

	t.Log(string(contents))

	if nextErr != io.EOF {
		t.Errorf("Next did not return EOF error: %v", err)
	}
}
