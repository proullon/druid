package driver

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/proullon/ramsql/engine/log"
	"reflect"
	"sync"
)

type field struct {
	Value reflect.Value
	Type  reflect.Type
}

// Rows implements the sql/driver Rows interface
type resultSet struct {
	columnNames []string
	rows        [][]field
	currentRow  int
	sync.Mutex
}

type Rows struct {
	resultSet resultSet
}

type queryResponse [][]interface{}

func parseResponse(body []byte) (r *Rows, err error) {
	var results queryResponse
	err = json.Unmarshal(body, &results)
	if err != nil {
		return &Rows{}, err
	}
	if len(results) == 0 {
		return &Rows{}, sql.ErrNoRows
	}

	var columnNames []string
	for _, val := range results[0] {
		columnNames = append(columnNames, val.(string))
	}
	var returnedRows [][]field
	for i := 1; i < len(results); i++ {
		var cols []field
		for _, val := range results[i] {
			cols = append(cols, field{Value: reflect.ValueOf(val), Type: reflect.TypeOf(val)})
		}
		returnedRows = append(returnedRows, cols)
	}

	resultSet := resultSet{
		columnNames: columnNames,
		rows:        returnedRows,
		currentRow:  0,
	}
	return &Rows{
		resultSet: resultSet,
	}, nil
}

// Columns returns the names of the columns. The number of
// columns of the result is inferred from the length of the
// slice.  If a particular column name isn't known, an empty
// string should be returned for that entry.
func (r *Rows) Columns() []string {
	return r.resultSet.columnNames
}

// Close closes the rows iterator.
func (r *Rows) Close() error {
	return nil
}

// Next is called to populate the next row of data into
// the provided slice. The provided slice will be the same
// size as the Columns() are wide.
func (r *Rows) Next(dest []driver.Value) (err error) {
	if !r.HasNextResultSet() {
		return errors.New("druid: no next data record")
	}

	data := r.resultSet.rows[r.resultSet.currentRow]
	if len(data) != len(dest) {
		return errors.New("druid: number of refs passed to scan does not match column count")
	}
	for i := range dest {
		if data[i].Type == nil {
			log.Warning("druid: data is nil", data[i])
			continue
		}
		switch data[i].Type.Name() {
		case "bool":
			dest[i] = data[i].Value.Interface().(bool)
		case "string":
			dest[i] = data[i].Value.Interface().(string)
		case "int":
			dest[i] = data[i].Value.Interface().(int)
		case "int64":
			dest[i] = data[i].Value.Interface().(int64)
		case "float64":
			dest[i] = data[i].Value.Interface().(float64)
		default:
			log.Warning("druid: can't scan type  [%s]", data[i].Type.Name())
		}
	}
	return r.NextResultSet()
}

// NextResultSet implements driver.RowsNextResultSet
func (r *Rows) NextResultSet() error {
	if !r.HasNextResultSet() {
		return errors.New("NextResult is empty")
	}
	r.resultSet.currentRow++
	return nil
}

// HasNextResultSet implements driver.RowsNextResultSet
func (r *Rows) HasNextResultSet() bool {
	return r.resultSet.currentRow != len(r.resultSet.rows)
}

func (r *Rows) setColumns(columns []string) {
	r.resultSet.columnNames = columns
}
