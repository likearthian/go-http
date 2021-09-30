package query

import (
	"fmt"
	"strings"
)

type Operator int
const (
	OP_NEQ Operator = iota
	OP_EQ
	OP_GT
	OP_GTE
	OP_LT
	OP_LTE
	OP_AND
	OP_OR
)

var nullBytes = []byte("null")

func (o *Operator) UnmarshalText(text []byte) error {
	str := string(text)
	if strings.EqualFold(str, "neq") {
		*o = OP_NEQ
	} else if strings.EqualFold(str, "eq") {
		*o = OP_EQ
	} else if strings.EqualFold(str, "gt") {
		*o = OP_GT
	} else if strings.EqualFold(str, "gte") {
		*o = OP_GTE
	} else if strings.EqualFold(str, "lt") {
		*o = OP_LT
	} else if strings.EqualFold(str, "lte") {
		*o = OP_LTE
	} else if strings.EqualFold(str, "and") {
		*o = OP_AND
	} else if strings.EqualFold(str, "or") {
		*o = OP_OR
	} else {
		return fmt.Errorf("cannot unmarshal the value %q into operator type", str)
	}

	return nil
}