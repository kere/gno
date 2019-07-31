package db

// Columns s
type Columns []DataColumn

// NewColumns n: column len, l : data length
func NewColumns(n, l int) Columns {
	cols := make([]DataColumn, n)
	if l > 0 {
		for i := 0; i < n; i++ {
			cols[i] = make([]interface{}, l)
		}
	}
	return cols
}

// SetMapRow set data
func (c Columns) SetMapRow(i int, mapRow MapRow, fields []string) {
	n := len(fields)
	for k := 0; k < n; k++ {
		c[k][i] = mapRow[fields[k]]
	}
}

// AddMapRow add row
func (c Columns) AddMapRow(mapRow MapRow, fields []string) {
	n := len(fields)
	for k := 0; k < n; k++ {
		v, isok := mapRow[fields[k]]
		if isok {
			c[k] = append(c[k], v)
		}
	}
}
