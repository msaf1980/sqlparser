# sqlparser - simple SQL parser

Based on https://github.com/marianogappa/sqlparser

Documentation on https://godoc.org/github.com/msaf1980/sqlparser
### Usage

```
package main

import (
	"fmt"
	"log"

	"github.com/msaf1980/sqlparser"
)

func main() {
	query, err := sqlparser.Parse("SELECT a, b, c FROM 'd' WHERE e = '1' AND f > '2'")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+#v", query)
}
```


### Example: SELECT version() as version

```
query, err := sqlparser.Parse(`SELECT version() as version`)

query.Query {
	Type: Select
	TableName: 
	Conditions: []
	Updates: map[]
	Inserts: []
	Fields: [version()]
}
```

### Example: SELECT works

```
query, err := sqlparser.Parse(`SELECT a FROM 'b'`)

query.Query {
	Type: Select
	TableName: b
	Conditions: []
	Updates: map[]
	Inserts: []
	Fields: [a]
}
```

### Example: SELECT with alias works

```
query, err := sqlparser.Parse(`SELECT a AS text FROM 'b'`)

query.Query {
	Type: Select
	TableName: b
	Conditions: []
	Updates: map[]
	Inserts: []
	Fields: [a]
}
```

### Example: SELECT with alias works

```
query, err := sqlparser.Parse(`SELECT version(a) AS version FROM 'b'`)

query.Query {
	Type: Select
	TableName: b
	Conditions: []
	Updates: map[]
	Inserts: []
	Fields: [version(a)]
}
```

### Example: SELECT works with lowercase

```
query, err := sqlparser.Parse(`select a fRoM 'b'`)

query.Query {
	Type: Select
	TableName: b
	Conditions: []
	Updates: map[]
	Inserts: []
	Fields: [a]
}
```

### Example: SELECT many fields works

```
query, err := sqlparser.Parse(`SELECT a, c, d FROM 'b'`)

query.Query {
	Type: Select
	TableName: b
	Conditions: []
	Updates: map[]
	Inserts: []
	Fields: [a c d]
}
```

### Example: SELECT with WHERE with = works

```
query, err := sqlparser.Parse(`SELECT a, c, d FROM 'b' WHERE a = ''`)

query.Query {
	Type: Select
	TableName: b
	Conditions: [
        {
            Operand1: a,
            Operand1Type: 1,
            Operator: Eq,
            Operand2: ,
            Operand2Type: 2,
        }]
	Updates: map[]
	Inserts: []
	Fields: [a c d]
}
```

### Example: SELECT with WHERE with < works

```
query, err := sqlparser.Parse(`SELECT a, c, d FROM 'b' WHERE a < '1'`)

query.Query {
	Type: Select
	TableName: b
	Conditions: [
        {
            Operand1: a,
            Operand1Type: 1,
            Operator: Lt,
            Operand2: 1,
            Operand2Type: 2,
        }]
	Updates: map[]
	Inserts: []
	Fields: [a c d]
}
```

### Example: SELECT with WHERE with <= works

```
query, err := sqlparser.Parse(`SELECT a, c, d FROM 'b' WHERE a <= '1'`)

query.Query {
	Type: Select
	TableName: b
	Conditions: [
        {
            Operand1: a,
            Operand1Type: 1,
            Operator: Lte,
            Operand2: 1,
            Operand2Type: 2,
        }]
	Updates: map[]
	Inserts: []
	Fields: [a c d]
}
```

### Example: SELECT with WHERE with > works

```
query, err := sqlparser.Parse(`SELECT a, c, d FROM 'b' WHERE a > '1'`)

query.Query {
	Type: Select
	TableName: b
	Conditions: [
        {
            Operand1: a,
            Operand1Type: 1,
            Operator: Gt,
            Operand2: 1,
            Operand2Type: 2,
        }]
	Updates: map[]
	Inserts: []
	Fields: [a c d]
}
```

### Example: SELECT with WHERE with >= works

```
query, err := sqlparser.Parse(`SELECT a, c, d FROM 'b' WHERE a >= '1'`)

query.Query {
	Type: Select
	TableName: b
	Conditions: [
        {
            Operand1: a,
            Operand1Type: 1,
            Operator: Gte,
            Operand2: 1,
            Operand2Type: 2,
        }]
	Updates: map[]
	Inserts: []
	Fields: [a c d]
}
```

### Example: SELECT with WHERE with != works

```
query, err := sqlparser.Parse(`SELECT a, c, d FROM 'b' WHERE a != '1'`)

query.Query {
	Type: Select
	TableName: b
	Conditions: [
        {
            Operand1: a,
            Operand1Type: 1,
            Operator: Ne,
            Operand2: 1,
            Operand2Type: 2,
        }]
	Updates: map[]
	Inserts: []
	Fields: [a c d]
}
```

### Example: SELECT with WHERE with != works (comparing field against another field)

```
query, err := sqlparser.Parse(`SELECT a, c, d FROM 'b' WHERE a != b`)

query.Query {
	Type: Select
	TableName: b
	Conditions: [
        {
            Operand1: a,
            Operand1Type: 1,
            Operator: Ne,
            Operand2: b,
            Operand2Type: 1,
        }]
	Updates: map[]
	Inserts: []
	Fields: [a c d]
}
```

### Example: SELECT * works

```
query, err := sqlparser.Parse(`SELECT * FROM 'b'`)

query.Query {
	Type: Select
	TableName: b
	Conditions: []
	Updates: map[]
	Inserts: []
	Fields: [*]
}
```

### Example: SELECT a, * works

```
query, err := sqlparser.Parse(`SELECT a, * FROM 'b'`)

query.Query {
	Type: Select
	TableName: b
	Conditions: []
	Updates: map[]
	Inserts: []
	Fields: [a *]
}
```

### Example: SELECT with WHERE with two conditions using AND works

```
query, err := sqlparser.Parse(`SELECT a, c, d FROM 'b' WHERE a != '1' AND b = '2'`)

query.Query {
	Type: Select
	TableName: b
	Conditions: [
        {
            Operand1: a,
            Operand1Type: 1,
            Operator: Ne,
            Operand2: 1,
            Operand2Type: 2,
        }
        {
            Operand1: b,
            Operand1Type: 1,
            Operator: Eq,
            Operand2: 2,
            Operand2Type: 2,
        }]
	Updates: map[]
	Inserts: []
	Fields: [a c d]
}
```

### Example: UPDATE works

```
query, err := sqlparser.Parse(`UPDATE 'a' SET b = 'hello' WHERE a = '1'`)

query.Query {
	Type: Update
	TableName: a
	Conditions: [
        {
            Operand1: a,
            Operand1Type: 1,
            Operator: Eq,
            Operand2: 1,
            Operand2Type: 2,
        }]
	Updates: map[b:hello]
	Inserts: []
	Fields: []
}
```

### Example: UPDATE works with simple quote inside

```
query, err := sqlparser.Parse(`UPDATE 'a' SET b = 'hello\'world' WHERE a = '1'`)

query.Query {
	Type: Update
	TableName: a
	Conditions: [
        {
            Operand1: a,
            Operand1Type: 1,
            Operator: Eq,
            Operand2: 1,
            Operand2Type: 2,
        }]
	Updates: map[b:hello\'world]
	Inserts: []
	Fields: []
}
```

### Example: UPDATE with multiple SETs works

```
query, err := sqlparser.Parse(`UPDATE 'a' SET b = 'hello', c = 'bye' WHERE a = '1'`)

query.Query {
	Type: Update
	TableName: a
	Conditions: [
        {
            Operand1: a,
            Operand1Type: 1,
            Operator: Eq,
            Operand2: 1,
            Operand2Type: 2,
        }]
	Updates: map[b:hello c:bye]
	Inserts: []
	Fields: []
}
```

### Example: UPDATE with multiple SETs and multiple conditions works

```
query, err := sqlparser.Parse(`UPDATE 'a' SET b = 'hello', c = 'bye' WHERE a = '1' AND b = '789'`)

query.Query {
	Type: Update
	TableName: a
	Conditions: [
        {
            Operand1: a,
            Operand1Type: 1,
            Operator: Eq,
            Operand2: 1,
            Operand2Type: 2,
        }
        {
            Operand1: b,
            Operand1Type: 1,
            Operator: Eq,
            Operand2: 789,
            Operand2Type: 2,
        }]
	Updates: map[b:hello c:bye]
	Inserts: []
	Fields: []
}
```

### Example: DELETE with WHERE works

```
query, err := sqlparser.Parse(`DELETE FROM 'a' WHERE b = '1'`)

query.Query {
	Type: Delete
	TableName: a
	Conditions: [
        {
            Operand1: b,
            Operand1Type: 1,
            Operator: Eq,
            Operand2: 1,
            Operand2Type: 2,
        }]
	Updates: map[]
	Inserts: []
	Fields: []
}
```

### Example: INSERT works

```
query, err := sqlparser.Parse(`INSERT INTO 'a' (b) VALUES ('1')`)

query.Query {
	Type: Insert
	TableName: a
	Conditions: []
	Updates: map[]
	Inserts: [[1]]
	Fields: [b]
}
```

### Example: INSERT with multiple fields works

```
query, err := sqlparser.Parse(`INSERT INTO 'a' (b,c,    d) VALUES ('1','2' ,  '3' )`)

query.Query {
	Type: Insert
	TableName: a
	Conditions: []
	Updates: map[]
	Inserts: [[1 2 3]]
	Fields: [b c d]
}
```

### Example: INSERT with multiple fields and multiple values works

```
query, err := sqlparser.Parse(`INSERT INTO 'a' (b,c,    d) VALUES ('1','2' ,  '3' ),('4','5' ,'6' )`)

query.Query {
	Type: Insert
	TableName: a
	Conditions: []
	Updates: map[]
	Inserts: [[1 2 3] [4 5 6]]
	Fields: [b c d]
}
```



### Example: empty query fails

```
query, err := sqlparser.Parse(``)

query type cannot be empty
```

### Example: SELECT without FROM fails

```
query, err := sqlparser.Parse(`SELECT`)

table name cannot be empty
```

### Example: SELECT without fields fails

```
query, err := sqlparser.Parse(`SELECT FROM 'a'`)

at SELECT: expected field to SELECT
```

### Example: SELECT with comma and empty field fails

```
query, err := sqlparser.Parse(`SELECT b, FROM 'a'`)

at SELECT: expected field to SELECT
```

### Example: SELECT with incomplete alias fails

```
query, err := sqlparser.Parse(`SELECT a AS`)

at AS: expected alias for a
```

### Example: SELECT with empty WHERE fails

```
query, err := sqlparser.Parse(`SELECT a, c, d FROM 'b' WHERE`)

at WHERE: empty WHERE clause
```

### Example: SELECT with WHERE with only operand fails

```
query, err := sqlparser.Parse(`SELECT a, c, d FROM 'b' WHERE a`)

at WHERE: condition without operator
```

### Example: Empty UPDATE fails

```
query, err := sqlparser.Parse(`UPDATE`)

table name cannot be empty
```

### Example: Incomplete UPDATE with table name fails

```
query, err := sqlparser.Parse(`UPDATE 'a'`)

at WHERE: WHERE clause is mandatory for UPDATE & DELETE
```

### Example: Incomplete UPDATE with table name and SET fails

```
query, err := sqlparser.Parse(`UPDATE 'a' SET`)

at WHERE: WHERE clause is mandatory for UPDATE & DELETE
```

### Example: Incomplete UPDATE with table name, SET with a field but no value and WHERE fails

```
query, err := sqlparser.Parse(`UPDATE 'a' SET b WHERE`)

at UPDATE: expected '='
```

### Example: Incomplete UPDATE with table name, SET with a field and = but no value and WHERE fails

```
query, err := sqlparser.Parse(`UPDATE 'a' SET b = WHERE`)

at UPDATE: expected quoted value
```

### Example: Incomplete UPDATE due to no WHERE clause fails

```
query, err := sqlparser.Parse(`UPDATE 'a' SET b = 'hello' WHERE`)

at WHERE: empty WHERE clause
```

### Example: Incomplete UPDATE due incomplete WHERE clause fails

```
query, err := sqlparser.Parse(`UPDATE 'a' SET b = 'hello' WHERE a`)

at WHERE: condition without operator
```

### Example: Empty DELETE fails

```
query, err := sqlparser.Parse(`DELETE FROM`)

table name cannot be empty
```

### Example: DELETE without WHERE fails

```
query, err := sqlparser.Parse(`DELETE FROM 'a'`)

at WHERE: WHERE clause is mandatory for UPDATE & DELETE
```

### Example: DELETE with empty WHERE fails

```
query, err := sqlparser.Parse(`DELETE FROM 'a' WHERE`)

at WHERE: empty WHERE clause
```

### Example: DELETE with WHERE with field but no operator fails

```
query, err := sqlparser.Parse(`DELETE FROM 'a' WHERE b`)

at WHERE: condition without operator
```

### Example: Empty INSERT fails

```
query, err := sqlparser.Parse(`INSERT INTO`)

table name cannot be empty
```

### Example: INSERT with no rows to insert fails

```
query, err := sqlparser.Parse(`INSERT INTO 'a'`)

at INSERT INTO: need at least one row to insert
```

### Example: INSERT with incomplete value section fails

```
query, err := sqlparser.Parse(`INSERT INTO 'a' (`)

at INSERT INTO: need at least one row to insert
```

### Example: INSERT with incomplete value section fails #2

```
query, err := sqlparser.Parse(`INSERT INTO 'a' (b`)

at INSERT INTO: need at least one row to insert
```

### Example: INSERT with incomplete value section fails #3

```
query, err := sqlparser.Parse(`INSERT INTO 'a' (b)`)

at INSERT INTO: need at least one row to insert
```

### Example: INSERT with incomplete value section fails #4

```
query, err := sqlparser.Parse(`INSERT INTO 'a' (b) VALUES`)

at INSERT INTO: need at least one row to insert
```

### Example: INSERT with incomplete row fails

```
query, err := sqlparser.Parse(`INSERT INTO 'a' (b) VALUES (`)

at INSERT INTO: value count doesn't match field count
```

### Example: INSERT * fails

```
query, err := sqlparser.Parse(`INSERT INTO 'a' (*) VALUES ('1')`)

at INSERT INTO: expected at least one field to insert
```

