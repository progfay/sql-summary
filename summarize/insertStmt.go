package summarize

import (
	"fmt"
	"strings"

	"github.com/pingcap/parser/ast"
	"github.com/willf/pad"
)

var (
	insertStmtTableNameNotFoundErr = fmt.Errorf("table name is not found in *ast.InsertStmt")

	alreadyInsertedTableMap = make(map[string]struct{})
)

func getInsertStmtTableName(node *ast.InsertStmt) (string, error) {
	tableSource, ok := node.Table.TableRefs.Left.(*ast.TableSource)
	if !ok {
		return "", insertStmtTableNameNotFoundErr
	}

	tableName, ok := tableSource.Source.(*ast.TableName)
	if !ok {
		return "", insertStmtTableNameNotFoundErr
	}

	return tableName.Name.String(), nil
}

func summarizeInsertStmt(node *ast.InsertStmt) (string, error) {
	var summary strings.Builder
	summary.WriteString(fmt.Sprintf("-- Insert rows count: %d\n", len(node.Lists)))

	var visited bool = false
	table, err := getInsertStmtTableName(node)
	if err == nil {
		_, visited = alreadyInsertedTableMap[table]
	}

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

	err = restore(&summary, node)
	if err != nil {
		return "", err
	}

	return summary.String(), nil
}
