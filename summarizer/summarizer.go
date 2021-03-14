package summarizer

import (
	"fmt"

	"github.com/pingcap/parser/ast"
)

var (
	NoChangeErr = fmt.Errorf("no change")
)

type Summarizer struct {
	alreadyInsertedTableMap map[string]struct{}
}

func New() *Summarizer {
	return &Summarizer{
		alreadyInsertedTableMap: make(map[string]struct{}),
	}
}

func (s *Summarizer) Summarize(node ast.StmtNode) (string, error) {
	switch node.(type) {
	case *ast.InsertStmt:
		node := node.(*ast.InsertStmt)
		return s.summarizeInsertStmt(node)

	default:
		return "", NoChangeErr
	}
}
