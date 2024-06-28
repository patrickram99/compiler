package object

import (
	"fmt"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

const (
	INTEGER_OBJ = "INTEGER"
	FLOAT_OBJ   = "FLOAT"
	BOOL_OBJ    = "BOOL"
	NULL_OBJ    = "NULL"
	RETURN_OBJ  = "RETURN_VAL"
)

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Float struct {
	Value float64
}

func (i *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) Inspect() string  { return fmt.Sprintf("%f", f.Value) }

type Bool struct {
	Value bool
}

func (b *Bool) Type() ObjectType { return BOOL_OBJ }
func (b *Bool) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type Null struct {
}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type ReturnVal struct {
	Value Object
}

func (rv *ReturnVal) Type() ObjectType { return RETURN_OBJ }
func (rv *ReturnVal) Inspect() string  { return rv.Value.Inspect() }
