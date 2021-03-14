package sqlsummary

import (
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/pingcap/parser"
	_ "github.com/pingcap/tidb/types/parser_driver"
	"github.com/progfay/sqlsummary/summarizer"
	// _ "github.com/pingcap/parser/test_driver"
)

func Run(w io.Writer, src io.Reader, maxCapacity int) {
	scanner := NewStatementScanner(src, maxCapacity)
	s := summarizer.New()

	for scanner.Scan() {
		statement := scanner.Text()

		p := parser.New()
		node, warns, err := p.Parse(statement, "", "")
		if err != nil {
			log.Printf("error occuerred, skip summarizing: %v\n", err)
			fmt.Fprint(w, statement+";")
			continue
		}

		for _, warn := range warns {
			log.Println(warn)
		}

		if len(node) == 0 {
			// statement only has comments.
			fmt.Fprint(w, statement)
			continue
		}

		if len(node) > 1 {
			log.Printf("StatementScanner.Text() return SQL Query with multiple statements, skip summarizing: %q\n", statement)
			fmt.Fprint(w, statement+";")
			continue
		}

		summary, err := s.Summarize(node[0])
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
