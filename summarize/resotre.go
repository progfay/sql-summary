package summarize

import (
	"io"

	"github.com/pingcap/parser/format"
)

const restoreFlags = format.RestoreStringDoubleQuotes | format.RestoreKeyWordUppercase | format.RestoreNameLowercase | format.RestoreNameBackQuotes

func restore(w io.Writer, stmt interface {
	Restore(*format.RestoreCtx) error
}) error {
	ctx := format.NewRestoreCtx(restoreFlags, w)
	err := stmt.Restore(ctx)
	return err
}
