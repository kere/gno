package db

import "database/sql"

// ScanRows db
// parse db.RowData from sql.rows
func ScanRows(rows *sql.Rows) (DataSet, error) {
	defer rows.Close()

	cols, _ := rows.Columns()
	colsNum := len(cols)

	result := DataSet{}
	var err error
	var row, tem []interface{}
	var rowData DataRow

	for rows.Next() {
		row = make([]interface{}, colsNum)
		tem = make([]interface{}, colsNum)

		for i := range row {
			tem[i] = &row[i]
		}

		if err = rows.Scan(tem...); err != nil {
			return nil, err
		}

		rowData = DataRow{}
		for i, col := range cols {
			switch row[i].(type) {
			case []byte:
				// prefix = byte_
				if len(col) > 5 && col[:5] == ColumnBytePrefix {
					rowData[col] = row[i].([]byte)
				} else {
					rowData[col] = string(row[i].([]byte))
				}
			default:
				rowData[col] = row[i]
			}
		}

		result = append(result, rowData)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
