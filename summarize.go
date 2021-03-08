package sqlsummary

import (
	"fmt"
	"strings"

	"github.com/pingcap/parser/ast"
	"github.com/willf/pad"
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
			columnLen := len(node.Lists[0])
			columns := make([]string, columnLen)
			if node.Columns == nil {
				for i := 0; i < columnLen; i++ {
					columns[i] = fmt.Sprintf("column%d", i)
				}
			} else {
				for i := 0; i < columnLen; i++ {
					columns[i] = node.Columns[i].String()
				}
			}

			longestColumnNameLength := 0
			for _, column := range columns {
				if len(column) > longestColumnNameLength {
					longestColumnNameLength = len(column)
				}
			}

			summary.WriteString("-- Row example: (\n")
			for i, item := range node.Lists[0] {
				summary.WriteString(fmt.Sprintf("-- \t/* %s */\t", pad.Left(columns[i], longestColumnNameLength, " ")))
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
