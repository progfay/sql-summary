package sqlsummary

import (
	"bytes"
	"errors"
	"io"
)

type StatementScanner struct {
	insideSingleQuote bool
	insideDoubleQuote bool
	insideBackQuote   bool
	escape            bool
	text              bytes.Buffer
	buf               []byte
	n                 int
	err               error
	r                 io.Reader
}

func (s *StatementScanner) Scan() bool {
	if s.err != nil {
		return false
	}
	s.text.Reset()

	for {
		for i := 0; i < s.n; i++ {
			if s.escape {
				s.escape = false
				continue
			}

			switch s.buf[i] {
			case ';':
				if s.insideSingleQuote || s.insideDoubleQuote || s.insideBackQuote {
					continue
				}
				s.n -= i + 1
				s.text.Write(s.buf[:i])
				s.buf = s.buf[i+1:]
				return true

			case '\\':
				s.escape = true

			case '\'':
				s.insideSingleQuote = !s.insideSingleQuote

			case '"':
				s.insideDoubleQuote = !s.insideDoubleQuote

			case '`':
				s.insideBackQuote = !s.insideBackQuote
			}
		}

		s.text.Write(s.buf[:s.n])
		s.n = 0

		if errors.Is(s.err, io.EOF) {
			return true
		}
		s.n, s.err = s.r.Read(s.buf)
		if s.err != nil && !errors.Is(s.err, io.EOF) {
			return false
		}
	}
}

func (s *StatementScanner) Text() string {
	return s.text.String()
}

func (s *StatementScanner) Err() error {
	return s.err
}

func NewStatementScanner(r io.Reader, maxCapacity int) *StatementScanner {
	return &StatementScanner{
		buf:  make([]byte, maxCapacity),
		r:    r,
	}
}
