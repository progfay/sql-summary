package sqlsummary

import (
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	_ "github.com/pingcap/tidb/types/parser_driver"
	"github.com/progfay/sqlsummary/summarizer"
	// _ "github.com/pingcap/parser/test_driver"
)

var (
	onlyCommentErr = fmt.Errorf("comment only")
)

func Run(w io.Writer, src io.Reader, maxCapacity int) {
	scanner := NewStatementScanner(src, maxCapacity)
	s := summarizer.New()

	for scanner.Scan() {
		statement := scanner.Text()

		node, err := parseStatement(statement)
		if err != nil {
			if errors.Is(err, onlyCommentErr) {
				fmt.Fprint(w, statement)
				continue
			}

			log.Println(err)
			fmt.Fprint(w, statement+";")
			continue
		}

		summary, err := s.Summarize(node)
		if err != nil {
			if errors.Is(err, summarizer.NoChangeErr) {
				fmt.Fprint(w, statement+";")
				continue
			}

			log.Println(err)
			continue
		}

		fmt.Fprintln(w, "\n"+summary+";")
	}
}

func parseStatement(statement string) (ast.StmtNode, error) {
	p := parser.New()
	nodes, _, err := p.Parse(statement, "", "")
	if err != nil {
		return nil, fmt.Errorf("error occurred, skip summarizing: %w", err)
	}

	if len(nodes) == 0 {
		return nil, onlyCommentErr
	}

	if len(nodes) > 1 {
		return nil, fmt.Errorf("StatementScanner.Text() return SQL Query with multiple statements, skip summarizing: %q", statement)
	}

	return nodes[0], nil
}
