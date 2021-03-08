package sqlsummary

import (
	"fmt"
	"io"
	"log"

	"github.com/pingcap/parser"
	_ "github.com/pingcap/tidb/types/parser_driver"
	// _ "github.com/pingcap/parser/test_driver"
)

func Run(w io.Writer, src io.Reader, maxCapacity int) {
	scanner := NewStatementScanner(src, maxCapacity)

	for scanner.Scan() {
		statement := scanner.Text()

		p := parser.New()
		nodes, _, err := p.Parse(statement, "", "")
		if err != nil {
			log.Println(err)
		}

		for _, node := range nodes {
			summary, err := summarize(node)
			if err != nil {
				log.Println(err)
				continue
			}

			fmt.Fprintln(w, summary + ";")
		}
	}
}
