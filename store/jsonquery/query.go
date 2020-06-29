package jsonquery

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/jinzhu/gorm"
)

const (
	sqlObject = `SELECT COALESCE(row_to_json(object_row),'{}'::json) FROM (%s) object_row`
	sqlArray  = `SELECT COALESCE(array_to_json(array_agg(row_to_json(array_row))),'[]'::json) FROM (%s) array_row`
)

var (
	reObjectStart = regexp.MustCompile(`(?m){{(\s+)?object(\s+)?}}`)
	reObjectEnd   = regexp.MustCompile(`(?m){{(\s+)?end_object(\s+)?}}`)
	reArrayStart  = regexp.MustCompile(`(?m){{(\s+)?array(\s+)?}}`)
	reArrayEnd    = regexp.MustCompile(`(?m){{(\s+)?end_array(\s+)?}}`)

	sqlObjectStart = "(" + strings.Split(sqlObject, "%s")[0]
	sqlObjectEnd   = strings.Split(sqlObject, "%s")[1] + ")"
	sqlArrayStart  = "(" + strings.Split(sqlArray, "%s")[0]
	sqlArrayEnd    = strings.Split(sqlArray, "%s")[1] + ")"

	defaultObject = []byte("{}")
	defaultArray  = []byte("[]")
)

// Prepare returns a query with expanded array/object blocks
// TODO: could use a refactor to utilize the replacer func
func Prepare(query string) string {
	query = reObjectStart.ReplaceAllString(query, sqlObjectStart)
	query = reObjectEnd.ReplaceAllString(query, sqlObjectEnd)
	query = reArrayStart.ReplaceAllString(query, sqlArrayStart)
	query = reArrayEnd.ReplaceAllString(query, sqlArrayEnd)

	return query
}

// Object returns an object result or a nil if nothing is found
func Object(db *gorm.DB, query string, values ...interface{}) ([]byte, error) {
	q := fmt.Sprintf(sqlObject, query)
	result, err := scanBytes(db.Raw(q, values...).Rows())
	if len(result) == 0 {
		result = nil
	}
	return result, err
}

// MustObject returns an object result or a default value if nothing is found
func MustObject(db *gorm.DB, query string, values ...interface{}) ([]byte, error) {
	result, err := Object(db, query, values...)
	if result == nil {
		result = defaultObject
	}
	return result, err
}

// Array returns an array result or a nil if nothing is found
func Array(db *gorm.DB, query string, values ...interface{}) ([]byte, error) {
	q := fmt.Sprintf(sqlArray, query)
	result, err := scanBytes(db.Raw(q, values...).Rows())
	if len(result) == 0 {
		result = nil
	}
	return result, err
}

// MustArray returns an array result or a default value if nothing is found
func MustArray(db *gorm.DB, query string, values ...interface{}) ([]byte, error) {
	result, err := Array(db, query, values...)
	if result == nil {
		result = defaultArray
	}
	return result, err
}

// scanBytes scans the first row and returns the raw content
func scanBytes(rows *sql.Rows, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []byte
	for rows.Next() {
		if err := rows.Scan(&data); err != nil {
			return nil, err
		}
	}

	return data, nil
}
