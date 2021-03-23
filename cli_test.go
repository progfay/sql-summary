package sqlsummary

import (
	"errors"
	"testing"
)

func Test_parseStatement (t *testing.T) {
	for _, testcase := range []struct{
		title string
		in string
		want error
	}{
		{
			title: "One statement",
			in: "SELECT * FROM users",
			want: nil,
		},
		{
			title: "Many Statements",
			in: "SELECT * FROM users; SELECT * FROM users",
			want: &multiStatementErr{S: "SELECT * FROM users; SELECT * FROM users"},
		},
		{
			title: "Only comments",
			in: "-- Comments\n/* Comments */",
			want: onlyCommentErr,
		},
	}{
		t.Run(testcase.title, func(t *testing.T) {
			_, err := parseStatement(testcase.in)

			if !errors.Is(err, testcase.want) {
				t.Errorf("want err %v, got %v", testcase.want, err)
			}
		})
	}
}
