package statement_test

import (
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/progfay/sqlsummary/statement"
)

type Testcase struct {
	title string
	in    string
	want  []string
}

func Test_StatementScanner(t *testing.T) {
	for _, testcase := range []Testcase{
		{
			title: "single statement",
			in:    "SELECT * FROM users",
			want: []string{
				"SELECT * FROM users",
			},
		},
		{
			title: "single statement with tail semicolon",
			in:    "SELECT * FROM users;",
			want: []string{
				"SELECT * FROM users",
				"",
			},
		},
		{
			title: "single statement with multi-lines",
			in: `
				SELECT
					*
				FROM
					users
			`,
			want: []string{
				`
				SELECT
					*
				FROM
					users
			`,
			},
		},
		{
			title: "many statement in single line",
			in:    "SELECT * FROM users;SELECT * FROM users;SELECT * FROM users",
			want: []string{
				"SELECT * FROM users",
				"SELECT * FROM users",
				"SELECT * FROM users",
			},
		},
		{
			title: "many statement in each lines",
			in: `
				SELECT * FROM users;
				SELECT * FROM users;
				SELECT * FROM users;
			`,
			want: []string{
				`
				SELECT * FROM users`,
				`
				SELECT * FROM users`,
				`
				SELECT * FROM users`,
				`
			`,
			},
		},
		{
			title: "statement with single-quote",
			in:    `SELECT * FROM users WHERE name = 'progfay;'`,
			want: []string{
				`SELECT * FROM users WHERE name = 'progfay;'`,
			},
		},
		{
			title: "statement with double-quote",
			in:    `SELECT * FROM users WHERE name = "progfay;"`,
			want: []string{
				`SELECT * FROM users WHERE name = "progfay;"`,
			},
		},
		{
			title: "statement with back-quote",
			in:    "SELECT `text`, `semicolon;` FROM `table`",
			want: []string{
				"SELECT `text`, `semicolon;` FROM `table`",
			},
		},
		{
			title: "statement with escaped double-quote",
			in:    `SELECT * FROM users WHERE description = "\";"`,
			want: []string{
				`SELECT * FROM users WHERE description = "\";"`,
			},
		},
	} {
		t.Run(testcase.title, func(t *testing.T) {
			scanner := statement.NewScanner(strings.NewReader(testcase.in))

			got := make([]string, 0)
			for scanner.Scan() {
				got = append(got, scanner.Text())
			}

			if !cmp.Equal(testcase.want, got) {
				t.Errorf(cmp.Diff(testcase.want, got))
			}

			err := scanner.Err()
			if err != nil && err != io.EOF {
				t.Error(err)
			}
		})
	}
}
