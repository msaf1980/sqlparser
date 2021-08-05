package query

import "strings"

// Query represents a parsed query
type Query struct {
	Type       Type
	TableName  string
	Conditions []Condition
	Updates    map[string]Operand
	Inserts    [][]Operand
	Fields     []string // Used for SELECT (i.e. SELECTed field names) and INSERT (INSERTEDed field names)
	Aliases    []string // Used for SELECT (i.e. SELECTed field_name AS alias_name)
}

// Type is the type of SQL query, e.g. SELECT/UPDATE
type Type int

const (
	// UnknownType is the zero value for a Type
	UnknownType Type = iota
	// Select represents a SELECT query
	Select
	// Update represents an UPDATE query
	Update
	// Insert represents an INSERT query
	Insert
	// Delete represents a DELETE query
	Delete
)

// TypeString is a string slice with the names of all types in order
var TypeString = []string{
	"UnknownType",
	"Select",
	"Update",
	"Insert",
	"Delete",
}

// Operator is between operands in a condition
type Operator int

const (
	// UnknownOperator is the zero value for an Operator
	UnknownOperator Operator = iota
	// Eq -> "="
	Eq
	// Ne -> "!="
	Ne
	// Gt -> ">"
	Gt
	// Lt -> "<"
	Lt
	// Gte -> ">="
	Gte
	// Lte -> "<="
	Lte
	// IN (..)
	In
)

// OperatorString is a string slice with the names of all operators in order
var OperatorString = []string{
	"UnknownOperator",
	"Eq",
	"Ne",
	"Gt",
	"Lt",
	"Gte",
	"Lte",
}

type OperandType int

const (
	OpUnknown OperandType = iota
	OpField
	OpQuoted
	OpNumber
)

type Operand interface {
	Dump() string
}

type OperandString struct {
	value string
}

func NewOperandString(value string) *OperandString {
	return &OperandString{value}
}

func (o *OperandString) Dump() string {
	return o.value
}

type OperandNumber struct {
	value string
}

func NewOperandNumber(value string) *OperandNumber {
	return &OperandNumber{value}
}

func (o *OperandNumber) Dump() string {
	return o.value
}

type OperandField struct {
	value string
}

func NewOperandField(value string) *OperandField {
	if len(value) > 0 && value[0] == '\'' {
		return &OperandField{value[1 : len(value)-1]}
	}
	return &OperandField{value}
}

func (o *OperandField) Dump() string {
	return o.value
}

type OperandStrArray struct {
	values []string
}

func NewOperandStrArray(value string) *OperandStrArray {
	return &OperandStrArray{[]string{value}}
}

func (o *OperandStrArray) Dump() string {
	return strings.Join(o.values, ",")
}

type OperandNumArray struct {
	value []string
}

// Condition is a single boolean condition in a WHERE clause
type Condition struct {
	// Operand1 is the left hand side operand
	Operand1 Operand
	// Operator is e.g. "=", ">"
	Operator Operator
	// Operand1 is the right hand side operand
	Operand2 Operand
}
