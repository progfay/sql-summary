package sqlsummary

import (
	"bufio"
	"io"
)

func splitStatement(data []byte, atEOF bool) (int, []byte, error) {
	var insideSingleQuote, insideDoubleQuote, insideBackQuote, escape bool
	for i := 0; i < len(data); i++ {
		if escape {
			escape = false
			continue
		}

		switch data[i] {
		case ';':
			if insideSingleQuote || insideDoubleQuote || insideBackQuote {
				continue
			}
			return i + 1, data[:i], nil

		case '\\':
			escape = true

		case '\'':
			insideSingleQuote = !insideSingleQuote

		case '"':
			insideDoubleQuote = !insideDoubleQuote

		case '`':
			insideBackQuote = !insideBackQuote
		}
	}

	return 0, data, bufio.ErrFinalToken
}

func NewStatementScanner(r io.Reader, bufferSize int) *bufio.Scanner {
	scanner := bufio.NewScanner(r)
	scanner.Split(splitStatement)
	buf := make([]byte, bufferSize)
	scanner.Buffer(buf, bufferSize)
	return scanner
}
