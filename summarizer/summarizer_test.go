package summarizer_test

import (
	"errors"
	"testing"

	"github.com/pingcap/parser"
	_ "github.com/pingcap/tidb/types/parser_driver"
	"github.com/progfay/sqlsummary/summarizer"
)

func Test_Summarize(t *testing.T) {
	for _, testcase := range []struct {
		title string
		in    string
		want  struct {
			summary string
			err     error
		}
	}{
		{
			title: "SELECT Statement",
			in:    "SELECT name from users WHERE id = 1",
			want: struct {
				summary string
				err     error
			}{
				summary: "",
				err:     summarizer.NoChangeErr,
			},
		},
		{
			title: "UPDATE Statement",
			in:    "UPDATE users SET password = 'pa55w0rd' WHERE id = 1",
			want: struct {
				summary string
				err     error
			}{
				summary: "",
				err:     summarizer.NoChangeErr,
			},
		},
		{
			title: "INSERT Statement",
			in:    "INSERT INTO users VALUES (1, 'admin', 'pa55w0rd')",
			want: struct {
				summary string
				err     error
			}{
				summary: "-- Insert rows count: 1\n-- Row example: (\n-- \t/* column0 */\t1\n-- \t/* column1 */\t\"admin\"\n-- \t/* column2 */\t\"pa55w0rd\"\n-- )\nINSERT INTO `users`",
				err:     nil,
			},
		},
		{
			title: "INSERT Statement with column names",
			in:    "INSERT INTO users(id, name, password) VALUES (1, 'admin', 'pa55w0rd')",
			want: struct {
				summary string
				err     error
			}{
				summary: "-- Insert rows count: 1\n-- Row example: (\n-- \t/*       id */\t1\n-- \t/*     name */\t\"admin\"\n-- \t/* password */\t\"pa55w0rd\"\n-- )\nINSERT INTO `users` (`id`,`name`,`password`)",
				err:     nil,
			},
		},
	} {
		t.Run(testcase.title, func(t *testing.T) {
			p := parser.New()
			node, err := p.ParseOneStmt(testcase.in, "", "")
			if err != nil {
				t.Error(err)
				return
			}

			s := summarizer.New()
			out, err := s.Summarize(node)

			if testcase.want.summary != out {
				t.Errorf("want %q, got %q", testcase.want.summary, out)
			}

			if !errors.Is(testcase.want.err, err) {
				t.Errorf("want error %v, got %v", testcase.want.err, err)
			}
		})
	}
}

func Test_Summarize_InsertRowExample(t *testing.T) {
	want := struct {
		first  string
		second string
	}{
		first:  "-- Insert rows count: 1\n-- Row example: (\n-- \t/* column0 */\t1\n-- \t/* column1 */\t\"admin\"\n-- \t/* column2 */\t\"pa55w0rd\"\n-- )\nINSERT INTO `users`",
		second: "-- Insert rows count: 1\nINSERT INTO `users`",
	}

	p := parser.New()
	s := summarizer.New()

	node, err := p.ParseOneStmt("INSERT INTO users VALUES (1, 'admin', 'pa55w0rd')", "", "")
	if err != nil {
		t.Error(err)
		return
	}
	got, err := s.Summarize(node)
	if err != nil {
		t.Error(err)
		return
	}
	if got != want.first {
		t.Errorf("first: want %q, got %q", want.first, got)
	}

	node, err = p.ParseOneStmt("INSERT INTO users VALUES (2, 'guest', 'password')", "", "")
	if err != nil {
		t.Error(err)
		return
	}
	got, err = s.Summarize(node)
	if err != nil {
		t.Error(err)
		return
	}
	if got != want.second {
		t.Errorf("first: want %q, got %q", want.second, got)
	}
}
