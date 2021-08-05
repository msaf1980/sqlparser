package sqlparser

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"text/template"

	"github.com/msaf1980/sqlparser/query"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	Name     string
	SQL      string
	Expected query.Query
	Err      error
	Ended    bool
}

type output struct {
	NoErrorExamples []testCase
	ErrorExamples   []testCase
	Types           []string
	Operators       []string
}

func TestSQL(t *testing.T) {
	ts := []testCase{
		{
			Name:     "empty query fails",
			SQL:      "",
			Expected: query.Query{},
			Err:      fmt.Errorf("query type cannot be empty"),
		},
		{
			Name:     "SELECT without FROM fails",
			SQL:      "SELECT",
			Expected: query.Query{Type: query.Select},
			Err:      fmt.Errorf("table name cannot be empty"),
		},
		{
			Name:     "SELECT without fields fails",
			SQL:      "SELECT FROM 'a'",
			Expected: query.Query{Type: query.Select},
			Err:      fmt.Errorf("at SELECT: expected field to SELECT"),
		},
		{
			Name:     "SELECT with comma and empty field fails",
			SQL:      "SELECT b, FROM 'a'",
			Expected: query.Query{Type: query.Select},
			Err:      fmt.Errorf("at SELECT: expected field to SELECT"),
		},
		{
			Name:     "SELECT with incomplete alias fails",
			SQL:      "SELECT a AS",
			Expected: query.Query{Type: query.Select},
			Err:      fmt.Errorf("at AS: expected alias for a"),
		},
		{
			Name:     "SELECT version() as version",
			SQL:      "SELECT version() as version",
			Expected: query.Query{Type: query.Select, Fields: []string{"version()"}, Aliases: []string{"version"}},
			Err:      nil,
		},
		{
			Name:     "SELECT works",
			SQL:      "SELECT a FROM b",
			Expected: query.Query{Type: query.Select, TableName: "b", Fields: []string{"a"}, Aliases: []string{""}},
			Err:      nil,
		},
		{
			Name:     "SELECT with alias works",
			SQL:      "SELECT a AS text FROM b",
			Expected: query.Query{Type: query.Select, TableName: "b", Fields: []string{"a"}, Aliases: []string{"text"}},
			Err:      nil,
		},
		{
			Name:     "SELECT with alias works",
			SQL:      "SELECT version(a) AS version FROM b",
			Expected: query.Query{Type: query.Select, TableName: "b", Fields: []string{"version(a)"}, Aliases: []string{"version"}},
			Err:      nil,
		},
		{
			Name:     "SELECT works with lowercase",
			SQL:      "select a fRoM b",
			Expected: query.Query{Type: query.Select, TableName: "b", Fields: []string{"a"}, Aliases: []string{""}},
			Err:      nil,
		},
		{
			Name:     "SELECT many fields works",
			SQL:      "SELECT a, c, d FROM b",
			Expected: query.Query{Type: query.Select, TableName: "b", Fields: []string{"a", "c", "d"}, Aliases: []string{"", "", ""}},
			Err:      nil,
		},
		{
			Name:     "SELECT with empty WHERE fails",
			SQL:      "SELECT a, c, d FROM b WHERE",
			Expected: query.Query{Type: query.Select, TableName: "b", Fields: []string{"a", "c", "d"}, Aliases: []string{"", "", ""}},
			Err:      fmt.Errorf("at WHERE: empty WHERE clause"),
		},
		{
			Name:     "SELECT with WHERE with only operand fails",
			SQL:      "SELECT a, c, d FROM b WHERE a",
			Expected: query.Query{Type: query.Select, TableName: "b", Fields: []string{"a", "c", "d"}, Aliases: []string{"", "", ""}},
			Err:      fmt.Errorf("at WHERE: condition without operator"),
		},
		{
			Name: "SELECT with WHERE with = works",
			SQL:  "SELECT a, c, d FROM b WHERE a = ''",
			Expected: query.Query{
				Type:      query.Select,
				TableName: "b",
				Fields:    []string{"a", "c", "d"}, Aliases: []string{"", "", ""},
				Conditions: []query.Condition{
					{Operand1: query.NewOperandField("a"), Operator: query.Eq, Operand2: query.NewOperandString("''")},
				},
			},
			Err: nil,
		},
		{
			Name: "SELECT with WHERE with < works",
			SQL:  "SELECT a, c, d FROM b WHERE a < '1'",
			Expected: query.Query{
				Type:      query.Select,
				TableName: "b",
				Fields:    []string{"a", "c", "d"}, Aliases: []string{"", "", ""},
				Conditions: []query.Condition{
					{Operand1: query.NewOperandField("a"), Operator: query.Lt, Operand2: query.NewOperandString("'1'")},
				},
			},
			Err: nil,
		},
		{
			Name: "SELECT with WHERE with <= works",
			SQL:  "SELECT a, c, d FROM b WHERE a <= '1'",
			Expected: query.Query{
				Type:      query.Select,
				TableName: "b",
				Fields:    []string{"a", "c", "d"}, Aliases: []string{"", "", ""},
				Conditions: []query.Condition{
					{Operand1: query.NewOperandField("a"), Operator: query.Lte, Operand2: query.NewOperandString("'1'")},
				},
			},
			Err: nil,
		},
		{
			Name: "SELECT with WHERE with > works",
			SQL:  "SELECT a, c, d FROM b WHERE a > '1'",
			Expected: query.Query{
				Type:      query.Select,
				TableName: "b",
				Fields:    []string{"a", "c", "d"}, Aliases: []string{"", "", ""},
				Conditions: []query.Condition{
					{Operand1: query.NewOperandField("a"), Operator: query.Gt, Operand2: query.NewOperandString("'1'")},
				},
			},
			Err: nil,
		},
		{
			Name: "SELECT with WHERE with >= works",
			SQL:  "SELECT a, c, d FROM b WHERE a >= '1'",
			Expected: query.Query{
				Type:      query.Select,
				TableName: "b",
				Fields:    []string{"a", "c", "d"}, Aliases: []string{"", "", ""},
				Conditions: []query.Condition{
					{Operand1: query.NewOperandField("a"), Operator: query.Gte, Operand2: query.NewOperandString("'1'")},
				},
			},
			Err: nil,
		},
		{
			Name: "SELECT with WHERE with != works",
			SQL:  "SELECT a, c, d FROM b WHERE a != '1'",
			Expected: query.Query{
				Type:      query.Select,
				TableName: "b",
				Fields:    []string{"a", "c", "d"}, Aliases: []string{"", "", ""},
				Conditions: []query.Condition{
					{Operand1: query.NewOperandField("a"), Operator: query.Ne, Operand2: query.NewOperandString("'1'")},
				},
			},
			Err: nil,
		},
		{
			Name: "SELECT with WHERE with != works (comparing field against another field)",
			SQL:  "SELECT a, c, d FROM b WHERE a != b",
			Expected: query.Query{
				Type:      query.Select,
				TableName: "b",
				Fields:    []string{"a", "c", "d"}, Aliases: []string{"", "", ""},
				Conditions: []query.Condition{
					{Operand1: query.NewOperandField("a"), Operator: query.Ne, Operand2: query.NewOperandField("b")},
				},
			},
			Err: nil,
		},
		{
			Name: "SELECT * works",
			SQL:  "SELECT * FROM b",
			Expected: query.Query{
				Type:       query.Select,
				TableName:  "b",
				Fields:     []string{"*"},
				Aliases:    []string{""},
				Conditions: nil,
			},
			Err: nil,
		},
		{
			Name: "SELECT a, * works",
			SQL:  "SELECT a, * FROM b",
			Expected: query.Query{
				Type:      query.Select,
				TableName: "b",
				Fields:    []string{"a", "*"}, Aliases: []string{"", ""},
				Conditions: nil,
			},
			Err: nil,
		},
		{
			Name: "SELECT with WHERE with two conditions using AND works",
			SQL:  "SELECT a, c, d FROM b WHERE a != '1' AND b = '2'",
			Expected: query.Query{
				Type:      query.Select,
				TableName: "b",
				Fields:    []string{"a", "c", "d"}, Aliases: []string{"", "", ""},
				Conditions: []query.Condition{
					{Operand1: query.NewOperandField("a"), Operator: query.Ne, Operand2: query.NewOperandString("'1'")},
					{Operand1: query.NewOperandField("b"), Operator: query.Eq, Operand2: query.NewOperandString("'2'")},
				},
			},
			Err: nil,
		},
		{
			Name:     "Empty UPDATE fails",
			SQL:      "UPDATE",
			Expected: query.Query{},
			Err:      fmt.Errorf("table name cannot be empty"),
		},
		{
			Name:     "Incomplete UPDATE with table name fails",
			SQL:      "UPDATE a",
			Expected: query.Query{},
			Err:      fmt.Errorf("at WHERE: WHERE clause is mandatory for UPDATE & DELETE"),
		},
		{
			Name:     "Incomplete UPDATE with table name and SET fails",
			SQL:      "UPDATE a SET",
			Expected: query.Query{},
			Err:      fmt.Errorf("at WHERE: WHERE clause is mandatory for UPDATE & DELETE"),
		},
		{
			Name:     "Incomplete UPDATE with table name, SET with a field but no value and WHERE fails",
			SQL:      "UPDATE a SET b WHERE",
			Expected: query.Query{},
			Err:      fmt.Errorf("at UPDATE: expected '='"),
		},
		{
			Name:     "Incomplete UPDATE with table name, SET with a field and = but no value and WHERE fails",
			SQL:      "UPDATE a SET b = WHERE",
			Expected: query.Query{},
			Err:      fmt.Errorf("at UPDATE: expected quoted value"),
		},
		{
			Name:     "Incomplete UPDATE due to no WHERE clause fails",
			SQL:      "UPDATE a SET b = 'hello' WHERE",
			Expected: query.Query{},
			Err:      fmt.Errorf("at WHERE: empty WHERE clause"),
		},
		{
			Name:     "Incomplete UPDATE due incomplete WHERE clause fails",
			SQL:      "UPDATE a SET b = 'hello' WHERE a",
			Expected: query.Query{},
			Err:      fmt.Errorf("at WHERE: condition without operator"),
		},
		{
			Name: "UPDATE works",
			SQL:  "UPDATE a SET b = 'hello' WHERE a = '1'",
			Expected: query.Query{
				Type:      query.Update,
				TableName: "a",
				Updates:   map[string]query.Operand{"b": query.NewOperandString("'hello'")},
				Conditions: []query.Condition{
					{Operand1: query.NewOperandField("a"), Operator: query.Eq, Operand2: query.NewOperandString("'1'")},
				},
			},
			Err: nil,
		},
		{
			Name: "UPDATE works with simple quote inside",
			SQL:  "UPDATE a SET b = 'hello\\'world' WHERE a = '1'",
			Expected: query.Query{
				Type:      query.Update,
				TableName: "a",
				Updates:   map[string]query.Operand{"b": query.NewOperandString("'hello\\'world'")},
				Conditions: []query.Condition{
					{Operand1: query.NewOperandField("a"), Operator: query.Eq, Operand2: query.NewOperandString("'1'")},
				},
			},
			Err: nil,
		},
		{
			Name: "UPDATE with multiple SETs works",
			SQL:  "UPDATE a SET b = 'hello', c = 'bye' WHERE a = '1'",
			Expected: query.Query{
				Type:      query.Update,
				TableName: "a",
				Updates:   map[string]query.Operand{"b": query.NewOperandString("'hello'"), "c": query.NewOperandString("'bye'")},
				Conditions: []query.Condition{
					{Operand1: query.NewOperandField("a"), Operator: query.Eq, Operand2: query.NewOperandString("'1'")},
				},
			},
			Err: nil,
		},
		{
			Name: "UPDATE with multiple SETs and multiple conditions works",
			SQL:  "UPDATE a SET b = 'hello', c = 'bye' WHERE a = '1' AND b = '789'",
			Expected: query.Query{
				Type:      query.Update,
				TableName: "a",
				Updates:   map[string]query.Operand{"b": query.NewOperandString("'hello'"), "c": query.NewOperandString("'bye'")},
				Conditions: []query.Condition{
					{Operand1: query.NewOperandField("a"), Operator: query.Eq, Operand2: query.NewOperandString("'1'")},
					{Operand1: query.NewOperandField("b"), Operator: query.Eq, Operand2: query.NewOperandString("'789'")},
				},
			},
			Err: nil,
		},
		{
			Name:     "Empty DELETE fails",
			SQL:      "DELETE FROM",
			Expected: query.Query{},
			Err:      fmt.Errorf("table name cannot be empty"),
		},
		{
			Name:     "DELETE without WHERE fails",
			SQL:      "DELETE FROM a",
			Expected: query.Query{},
			Err:      fmt.Errorf("at WHERE: WHERE clause is mandatory for UPDATE & DELETE"),
		},
		{
			Name:     "DELETE with empty WHERE fails",
			SQL:      "DELETE FROM a WHERE",
			Expected: query.Query{},
			Err:      fmt.Errorf("at WHERE: empty WHERE clause"),
		},
		{
			Name:     "DELETE with WHERE with field but no operator fails",
			SQL:      "DELETE FROM a WHERE b",
			Expected: query.Query{},
			Err:      fmt.Errorf("at WHERE: condition without operator"),
		},
		{
			Name: "DELETE with WHERE works",
			SQL:  "DELETE FROM a WHERE b = '1'",
			Expected: query.Query{
				Type:      query.Delete,
				TableName: "a",
				Conditions: []query.Condition{
					{Operand1: query.NewOperandField("b"), Operator: query.Eq, Operand2: query.NewOperandString("'1'")},
				},
			},
			Err: nil,
		},
		{
			Name:     "Empty INSERT fails",
			SQL:      "INSERT INTO",
			Expected: query.Query{},
			Err:      fmt.Errorf("table name cannot be empty"),
		},
		{
			Name:     "INSERT with no rows to insert fails",
			SQL:      "INSERT INTO a",
			Expected: query.Query{},
			Err:      fmt.Errorf("at INSERT INTO: need at least one row to insert"),
		},
		{
			Name:     "INSERT with incomplete value section fails",
			SQL:      "INSERT INTO a (",
			Expected: query.Query{},
			Err:      fmt.Errorf("at INSERT INTO: need at least one row to insert"),
		},
		{
			Name:     "INSERT with incomplete value section fails #2",
			SQL:      "INSERT INTO a (b",
			Expected: query.Query{},
			Err:      fmt.Errorf("at INSERT INTO: need at least one row to insert"),
		},
		{
			Name:     "INSERT with incomplete value section fails #3",
			SQL:      "INSERT INTO a (b)",
			Expected: query.Query{},
			Err:      fmt.Errorf("at INSERT INTO: need at least one row to insert"),
		},
		{
			Name:     "INSERT with incomplete value section fails #4",
			SQL:      "INSERT INTO a (b) VALUES",
			Expected: query.Query{},
			Err:      fmt.Errorf("at INSERT INTO: need at least one row to insert"),
		},
		{
			Name:     "INSERT with incomplete row fails",
			SQL:      "INSERT INTO a (b) VALUES (",
			Expected: query.Query{},
			Err:      fmt.Errorf("at INSERT INTO: value count doesn't match field count"),
		},
		{
			Name: "INSERT works",
			SQL:  "INSERT INTO a (b) VALUES ('1')",
			Expected: query.Query{
				Type:      query.Insert,
				TableName: "a",
				Fields:    []string{"b"},
				Inserts:   [][]query.Operand{{query.NewOperandString("'1'")}},
			},
			Err: nil,
		},
		{
			Name: "INSERT * fails",
			SQL:  "INSERT INTO a (*) VALUES ('1')",
			Expected: query.Query{
				Type:      query.Insert,
				TableName: "a",
				Fields:    []string{"*"},
				Inserts:   [][]query.Operand{{query.NewOperandString("'1'")}},
			},
			Err: fmt.Errorf("at INSERT INTO: expected at least one field to insert"),
		},
		{
			Name: "INSERT with multiple fields works",
			SQL:  "INSERT INTO a (b,c,    d) VALUES ('1','2' ,  '3' )",
			Expected: query.Query{
				Type:      query.Insert,
				TableName: "a",
				Fields:    []string{"b", "c", "d"},
				Inserts:   [][]query.Operand{{query.NewOperandString("'1'"), query.NewOperandString("'2'"), query.NewOperandString("'3'")}},
			},
			Err: nil,
		},
		{
			Name: "INSERT with multiple fields and multiple values works",
			SQL:  "INSERT INTO a (b,c,    d) VALUES ('1','2' ,  '3' ),('4','5' ,'6' )",
			Expected: query.Query{
				Type:      query.Insert,
				TableName: "a",
				Fields:    []string{"b", "c", "d"},
				Inserts: [][]query.Operand{
					{query.NewOperandString("'1'"), query.NewOperandString("'2'"), query.NewOperandString("'3'")},
					{query.NewOperandString("'4'"), query.NewOperandString("'5'"), query.NewOperandString("'6'")},
				},
			},
			Err: nil,
		},
	}

	output := output{Types: query.TypeString, Operators: query.OperatorString}
	for _, tc := range ts {
		t.Run(tc.Name, func(t *testing.T) {
			actual, err := ParseMany([]string{tc.SQL})
			if err != nil {
				if errPos, ok := err.(*ErrorWithPos); ok {
					fmt.Fprintln(os.Stderr, "")
					errPos.PrintPosError(tc.SQL, os.Stderr)
				}
			}
			if tc.Err != nil && err == nil {
				t.Errorf("Error should have been %v", tc.Err)
			}
			if tc.Err == nil && err != nil {
				t.Errorf("Error should have been nil but was %v", err)
			}
			if tc.Err != nil && err != nil {
				require.Equal(t, tc.Err.Error(), err.Error(), "Unexpected error")
			}
			if len(actual) > 0 {
				require.Equal(t, tc.Expected, actual[0], "Query didn't match expectation")
			}
			if tc.Err != nil {
				output.ErrorExamples = append(output.ErrorExamples, tc)
			} else {
				output.NoErrorExamples = append(output.NoErrorExamples, tc)
			}
		})
	}
	createReadme(output)
}

func TestWhere(t *testing.T) {
	ts := []testCase{
		{
			Name:     "empty query fails",
			SQL:      "",
			Expected: query.Query{},
			Err:      fmt.Errorf("at WHERE: empty WHERE clause"),
			Ended:    true,
		},
		{
			Name: "WHERE a",
			SQL:  "a ",
			Expected: query.Query{
				Conditions: []query.Condition{
					{Operand1: query.NewOperandField("a"), Operator: query.UnknownOperator, Operand2: nil},
				},
			},
			Err:   nil,
			Ended: true,
		},
		{
			Name: "WHERE a = ''",
			SQL:  "a = ''",
			Expected: query.Query{
				Conditions: []query.Condition{
					{Operand1: query.NewOperandField("a"), Operator: query.Eq, Operand2: query.NewOperandString("''")},
				},
			},
			Err:   nil,
			Ended: true,
		},
		{
			Name: "WHERE a = 1",
			SQL:  "a>=1",
			Expected: query.Query{
				Conditions: []query.Condition{
					{Operand1: query.NewOperandField("a"), Operator: query.Gte, Operand2: query.NewOperandNumber("1")},
				},
			},
			Err:   nil,
			Ended: true,
		},
		{
			Name: "WHERE a >= 1.24",
			SQL:  "a>= 1.24",
			Expected: query.Query{
				Conditions: []query.Condition{
					{Operand1: query.NewOperandField("a"), Operator: query.Gte, Operand2: query.NewOperandNumber("1.24")},
				},
			},
			Err:   nil,
			Ended: true,
		},
		{
			Name: "WHERE a >= -1.21",
			SQL:  "a>=-1.21",
			Expected: query.Query{
				Conditions: []query.Condition{
					{Operand1: query.NewOperandField("a"), Operator: query.Gte, Operand2: query.NewOperandNumber("-1.21")},
				},
			},
			Err:   nil,
			Ended: true,
		},
		{
			Name: "WHERE a = 1 AND b > a1",
			SQL:  "a = 1 AND b > a1",
			Expected: query.Query{
				Conditions: []query.Condition{
					{Operand1: query.NewOperandField("a"), Operator: query.Eq, Operand2: query.NewOperandNumber("1")},
					{Operand1: query.NewOperandField("b"), Operator: query.Gt, Operand2: query.NewOperandField("a1")},
				},
			},
			Err:   nil,
			Ended: true,
		},
		{
			Name: "ERROR (a1) WHERE a = 1a",
			SQL:  "a = 1a",
			Expected: query.Query{
				Conditions: []query.Condition{
					{Operand1: query.NewOperandField("a"), Operator: query.Eq, Operand2: nil},
				},
			},
			Err:   fmt.Errorf("at WHERE: expected quoted value"),
			Ended: false,
		},
	}

	for _, tc := range ts {
		t.Run(tc.Name, func(t *testing.T) {
			var p parser
			// init parser internals
			p.step = stepWhereField
			p.sql = tc.SQL
			p.sqlUpper = strings.ToUpper(tc.SQL)

			ended, err := p.parseWhere()
			if err != nil {
				if errPos, ok := err.(*ErrorWithPos); ok {
					fmt.Fprintln(os.Stderr, "")
					errPos.PrintPosError(tc.SQL, os.Stderr)
				}
			}
			if tc.Err != nil && err == nil {
				t.Errorf("Error should have been %v", tc.Err)
			}
			if tc.Err == nil && err != nil {
				t.Errorf("Error should have been nil but was %v", err)
			}
			if tc.Ended != ended {
				t.Errorf("End not detected")
			}
			if tc.Err != nil && err != nil {
				require.Equal(t, tc.Err.Error(), err.Error(), "Unexpected error")
			}
			require.Equal(t, tc.Expected, p.query, "Query didn't match expectation")
		})
	}
}

func BenchmarkSQLSelect(b *testing.B) {
	sql := "SELECT a AS text FROM b WHERE c = 'c' AND d = 'd'"
	for i := 0; i < b.N; i++ {
		q, err := Parse(sql)
		if err != nil {
			b.Errorf("Error should have been %v: %v", err, q)
		}
	}
}

func BenchmarkSQLInsert(b *testing.B) {
	sql := "INSERT INTO a (b,c,    d) VALUES ('1','2' ,  '3' )"
	for i := 0; i < b.N; i++ {
		q, err := Parse(sql)
		if err != nil {
			b.Errorf("Error should have been %v: %v", err, q)
		}
	}
}

func createReadme(out output) {
	content, err := ioutil.ReadFile("README.template")
	if err != nil {
		log.Fatal(err)
	}
	t := template.Must(template.New("").Parse(string(content)))
	f, err := os.Create("README.md")
	if err != nil {
		log.Fatal(err)
	}
	if err := t.Execute(f, out); err != nil {
		log.Fatal(err)
	}
}
