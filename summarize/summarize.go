package summarize

import (
	"github.com/pingcap/parser/ast"
)

func Summarize(node ast.StmtNode) (string, error) {
	switch node.(type) {
	case *ast.InsertStmt:
		node := node.(*ast.InsertStmt)
		return summarizeInsertStmt(node)

	default:
		return "", nil
	}
}
