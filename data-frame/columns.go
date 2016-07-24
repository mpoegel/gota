package df

import (
	"errors"
	"reflect"
)

// column represents a column inside a DataFrame
type column struct {
	cells    Cells
	colType  string
	colName  string
	numChars int
	empty    Cell
}

type columns []column

// newCol is the constructor for a new Column with the given colName and elements
func newCol(colName string, elements Cells) (*column, error) {
	col := column{
		colName: colName,
	}
	col, err := col.append(elements...)
	if err != nil {
		return nil, err
	}

	return &col, nil
}

func (col *column) ParseColumn(t string) error {
	switch t {
	case "string":
		newcells := Strings(col.cells)
		newcol, err := newCol(col.colName, newcells)
		if err != nil {
			return err
		}
		*col = *newcol
	case "int":
		newcells := Ints(col.cells)
		newcol, err := newCol(col.colName, newcells)
		if err != nil {
			return err
		}
		*col = *newcol
	case "float":
		newcells := Floats(col.cells)
		newcol, err := newCol(col.colName, newcells)
		if err != nil {
			return err
		}
		*col = *newcol
	case "bool":
		newcells := Bools(col.cells)
		newcol, err := newCol(col.colName, newcells)
		if err != nil {
			return err
		}
		*col = *newcol
	default:
		return errors.New("Can't parse the given type")
	}

	return nil
}

func (col *column) recountNumChars() {
	numChars := len(col.colName)
	for _, cell := range col.cells {
		cellStr := cell.String()
		if len(cellStr) > numChars {
			numChars = len(cellStr)
		}
	}

	col.numChars = numChars
}

// Append will add a value or values to a column
func (col column) append(values ...Cell) (column, error) {
	if len(values) == 0 {
		col.recountNumChars()
		return col, nil
	}

	col.empty = values[0].NA()
	for _, v := range values {
		t := reflect.TypeOf(v).String()
		if col.colType == "" {
			col.colType = t
		} else {
			if t != col.colType {
				return col, errors.New("Can't have elements of different type on the same column")
			}
		}

		col.cells = append(col.cells, v)
	}

	col.recountNumChars()

	return col, nil
}

func (col column) copy() column {
	cs := make(Cells, 0, len(col.cells))
	for _, v := range col.cells {
		cs = append(cs, v.Copy())
	}
	newcol := column{
		cells:    cs,
		colType:  col.colType,
		colName:  col.colName,
		numChars: col.numChars,
		empty:    col.empty.Copy(),
	}
	return newcol
}

func (col column) HasNA() bool {
	for _, v := range col.cells {
		if v.IsNA() {
			return true
		}
	}
	return false
}

func (col column) NA() []bool {
	naArray := make([]bool, len(col.cells))
	for k, v := range col.cells {
		if v.IsNA() {
			naArray[k] = true
		} else {
			naArray[k] = false
		}
	}
	return naArray
}

func (col column) AsString() []string {
	arr := make([]string, len(col.cells))
	for k, v := range col.cells {
		arr[k] = v.String()
	}
	return arr
}

func (col column) AsInt() []int {
	arr := make([]int, len(col.cells))
	for k, v := range col.cells {
		x, _ := v.Int()
		arr[k] = *x
	}
	return arr
}

func (col column) AsFloat() []float64 {
	arr := make([]float64, len(col.cells))
	for k, v := range col.cells {
		x, _ := v.Float()
		arr[k] = *x
	}
	return arr
}

func (col column) AsBool() []bool {
	arr := make([]bool, len(col.cells))
	for k, v := range col.cells {
		x, _ := v.Bool()
		arr[k] = *x
	}
	return arr
}
