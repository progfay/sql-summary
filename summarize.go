package sqlsummary

import (
	"fmt"
	"strings"

	"github.com/pingcap/parser/ast"
)

func summarize(node ast.StmtNode) (string, error) {
	var summary strings.Builder

	switch node.(type) {
	case *ast.InsertStmt:
		node := node.(*ast.InsertStmt)
		summary.WriteString(fmt.Sprintf("-- Insert rows count: %d\n", len(node.Lists)))

		summary.WriteString("-- Row example: (\n")
		for _, item := range node.Lists[0] {
			summary.WriteString("-- \t")
			item.Format(&summary)
			summary.WriteString("\n")
		}
		summary.WriteString("-- )\n")

		node.Lists = nil
	}

	err := restore(&summary, node)
	if err != nil {
		return "", err
	}

	return summary.String(), nil
}
