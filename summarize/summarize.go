package summarize

import (
	"fmt"

	"github.com/pingcap/parser/ast"
)

var (
	NoChangeErr = fmt.Errorf("no change")
)

func Summarize(node ast.StmtNode) (string, error) {
	switch node.(type) {
	case *ast.InsertStmt:
		node := node.(*ast.InsertStmt)
		return summarizeInsertStmt(node)

	default:
		return "", NoChangeErr
	}
}
