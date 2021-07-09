package query

import (
	"fmt"
	"strconv"
	"strings"
)

type QueryValue struct {
	Value string
	Operator Operator
	NextChain *QueryValue
	NextOperator Operator
}

func (q *QueryValue) Int() (int64, error) {
	return strconv.ParseInt(q.Value, 10, 64)
}

func (q *QueryValue) Float() (float64, error) {
	return strconv.ParseFloat(q.Value, 64)
}

//

func ParseQueryValue(str string) (*QueryValue, error) {
	strArr := strings.Split(str, "|")
	if len(strArr) == 1 {
		return &QueryValue{Value: str, Operator: OP_EQ}, nil
	}

	if len(strArr) == 2 {
		val := strArr[1]
		var op Operator
		if err := op.UnmarshalText([]byte(strArr[0])); err != nil {
			return nil, err
		}

		return &QueryValue{Value: val, Operator: op}, nil
	}

	if len(strArr) == 3 {
		return nil, fmt.Errorf("invalid query string: %s", str)
	}

	val := strArr[1]
	var op Operator
	var nextOp Operator
	if err := op.UnmarshalText([]byte(strArr[0])); err != nil {
		return nil, err
	}

	if err := nextOp.UnmarshalText([]byte(strArr[2])); err != nil {
		return nil, err
	}

	q := &QueryValue{
		Value:       val,
		Operator:     op,
		NextChain:    nil,
		NextOperator: nextOp,
	}

	newArr := strArr[3:]
	next, err := ParseQueryValue(strings.Join(newArr, "|"))
	if err != nil {
		return nil, err
	}

	q.NextChain = next
	return q, nil
}