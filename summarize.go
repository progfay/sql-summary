package sqlsummary

import (
	"fmt"
	"strings"

	"github.com/pingcap/parser/ast"
)

var (
	alreadyInsertedTableMap = make(map[string]struct{})
)

func summarize(node ast.StmtNode) (string, error) {
	var summary strings.Builder

	switch node.(type) {
	case *ast.InsertStmt:
		node := node.(*ast.InsertStmt)
		summary.WriteString(fmt.Sprintf("-- Insert rows count: %d\n", len(node.Lists)))

		table := node.Table.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName).Name.String()
		_, visited := alreadyInsertedTableMap[table]
		if !visited {
			summary.WriteString("-- Row example: (\n")
			for _, item := range node.Lists[0] {
				summary.WriteString("-- \t")
				item.Format(&summary)
				summary.WriteString("\n")
			}
			summary.WriteString("-- )\n")
			alreadyInsertedTableMap[table] = struct{}{}
		}

		node.Lists = nil
	}

	err := restore(&summary, node)
	if err != nil {
		return "", err
	}

	return summary.String(), nil
}
