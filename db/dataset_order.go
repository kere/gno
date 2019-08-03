package db

import "github.com/kere/gno/libs/util"

// NewInt64Indexs new
func NewInt64Indexs(dataset *DataSet, field string) util.Int64Indexs {
	fieldI := dataset.FieldI(field)
	if fieldI < 0 {
		panic(ErrNoField)
	}
	col := dataset.Columns[fieldI]

	l := dataset.Len()
	vals := make([][2]int64, l)
	for i := 0; i < l; i++ {
		vals[i] = [2]int64{col[i].(int64), int64(i)}
	}

	return vals
}
