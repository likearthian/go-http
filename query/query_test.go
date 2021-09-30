package query

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestParseQueryValue(t *testing.T) {
	str := "gt|10|or|lt|25|and|gte|65"
	q, err := ParseQueryValue(str)
	if err != nil {
		t.Error(err.Error())
	}

	t.Logf(spew.Sdump(q))
}
