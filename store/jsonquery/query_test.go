package jsonquery

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepare(t *testing.T) {
	examples := [][]string{
		{"", ""},
		{"{}", "{}"},
		{"{{ array }}DATA{{ end_array }}", "(SELECT COALESCE(array_to_json(array_agg(row_to_json(array_row))),'[]'::json) FROM (DATA) array_row)"},
		{"{{array    }}DATA{{    end_array}}", "(SELECT COALESCE(array_to_json(array_agg(row_to_json(array_row))),'[]'::json) FROM (DATA) array_row)"},
		{"{{ object }}DATA {{ end_object }}", "(SELECT COALESCE(row_to_json(object_row),'{}'::json) FROM (DATA ) object_row)"},
	}

	for _, ex := range examples {
		assert.Equal(t, ex[1], Prepare(ex[0]))
	}
}
