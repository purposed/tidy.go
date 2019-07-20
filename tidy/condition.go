package tidy

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Field represents a file field.
type Field string

// Tracked file fields.
const (
	Extension Field = "extension"
	Name      Field = "name"
	Created   Field = "created"
	Modified  Field = "modified"
	FileType  Field = "type"
	Size      Field = "size"
	Age       Field = "age"
	Path      Field = "path"
)

// Fields aggregates all field types.
var Fields = []Field{Extension, Name, Created, Modified, FileType, Size, Age, Path}

// BoolOperator represents a boolean operator.
type BoolOperator string

// Available boolean operators.
const (
	And BoolOperator = "and"
	Or  BoolOperator = "or"
	Xor BoolOperator = "xor"
)

// BoolOperators aggregates all bool ops.
var BoolOperators = []BoolOperator{And, Or, Xor}

// Operator represents an operator in a condition.
type Operator string

// Available operators.
const (
	Eq          Operator = "="
	Neq         Operator = "!="
	Lt          Operator = "<"
	Gt          Operator = ">"
	Leq         Operator = "<="
	Geq         Operator = ">="
	StartsWith  Operator = "^="
	NStartsWith Operator = "!^="
	EndsWith    Operator = "$="
	NEndsWith   Operator = "!$="
	Contains    Operator = "?="
	NContains   Operator = "!?="
)

// Operators aggregates all operator types.
var Operators = []Operator{Eq, Neq, Lt, Gt, Leq, Geq, StartsWith, NStartsWith, EndsWith, NEndsWith, Contains, NContains}

// Condition represents an evaluable expression.
type Condition interface {
	Evaluate(file *File) (bool, error)
}

// BaseCondition represents a simple condition.
type BaseCondition struct {
	Field Field

	Operator Operator
	Value    string
}

// Evaluate evaluates the condition.
func (b *BaseCondition) Evaluate(file *File) (bool, error) {
	fieldValue := file.GetField(b.Field)

	switch b.Operator {
	case Eq:
		return fieldValue == b.Value, nil
	case Neq:
		return fieldValue != b.Value, nil
	case StartsWith:
		return strings.HasPrefix(fieldValue.(string), b.Value), nil
	case NStartsWith:
		return !strings.HasPrefix(fieldValue.(string), b.Value), nil
	case EndsWith:
		return strings.HasSuffix(fieldValue.(string), b.Value), nil
	case NEndsWith:
		return !strings.HasSuffix(fieldValue.(string), b.Value), nil
	case Contains:
		return strings.Contains(fieldValue.(string), b.Value), nil
	case NContains:
		return !strings.Contains(fieldValue.(string), b.Value), nil
	case Gt:
		if num, ok := fieldValue.(float64); ok {
			f, err := strconv.ParseFloat(b.Value, 32)
			if err != nil {
				return false, errors.New("value in config is not a float")
			}
			return num > f, nil
		} else if date, ok := fieldValue.(time.Duration); ok {
			var dur time.Duration
			var err error
			dur, err = time.ParseDuration(b.Value)
			if err != nil {
				if strings.HasPrefix(err.Error(), "time: unknown unit d") {
					// Workaround to get "d" unit.
					intermed := strings.ReplaceAll(b.Value, "d", "h")
					dur, err = time.ParseDuration(intermed)
					if err != nil {
						return false, nil
					}
					dur = dur * time.Duration(24)
				} else {
					return false, err
				}
			}
			return date > dur, nil
		}
		return false, errors.New("unsupported operator")
	default:
		return false, fmt.Errorf("unsupported operator: %s", b.Operator)
	}
}

// JoinedCondition represents a conjunction of conditions.
type JoinedCondition struct {
	LeftCond  Condition
	Op        BoolOperator
	RightCond Condition
}

// Evaluate evaluates the condition.
func (b *JoinedCondition) Evaluate(f *File) (bool, error) {
	l, err := b.LeftCond.Evaluate(f)
	if err != nil {
		return false, err
	}
	r, err := b.RightCond.Evaluate(f)
	if err != nil {
		return false, err
	}

	switch b.Op {
	case And:
		return l && r, nil
	case Or:
		return l || r, nil
	case Xor:
		return (l || r) && !(l && r), nil
	default:
		return false, fmt.Errorf("unknown operator: %s", b.Op)
	}
}

// ParseCondition parses a raw condition into an evaluable condition object.
func ParseCondition(raw string) (Condition, error) {
	p := newParser(Tokenize(raw))
	return p.Parse()
}
