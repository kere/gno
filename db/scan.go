package db

import (
	"database/sql"
)

// ScanToDataSet db
func ScanToDataSet(rows *sql.Rows) (DataSet, error) {
	var result DataSet
	cols, err := rows.Columns()
	if err != nil {
		return result, err
	}
	colsNum := len(cols)

	typs, err := rows.ColumnTypes()
	if err != nil {
		return result, err
	}

	fields := make([]string, colsNum)
	typItems := make([]ColType, colsNum)
	for i := 0; i < colsNum; i++ {
		typItems[i] = NewColType(typs[i])
		fields[i] = typs[i].Name()
	}

	result = DataSet{Types: typItems, Fields: fields, Columns: make([]DataColumn, colsNum)}

	var row, tem []interface{}

	for rows.Next() {
		row = make([]interface{}, colsNum)
		tem = make([]interface{}, colsNum)

		for i := 0; i < colsNum; i++ {
			tem[i] = &row[i]
		}

		if err = rows.Scan(tem...); err != nil {
			return result, err
		}

		result.AddDataRow(row)
	}

	return result, rows.Err()
}

// ScanToMapRows db
// parse db.RowData from sql.rows
func ScanToMapRows(rows *sql.Rows) (MapRows, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	colsNum := len(cols)
	mapRows := MapRows{}

	var row, tem []interface{}

	for rows.Next() {
		row = make([]interface{}, colsNum)
		tem = make([]interface{}, colsNum)

		for i := 0; i < colsNum; i++ {
			tem[i] = &row[i]
		}

		if err = rows.Scan(tem...); err != nil {
			return nil, err
		}

		rowData := make(map[string]interface{})

		for i := 0; i < colsNum; i++ {
			rowData[cols[i]] = row[i]
		}

		mapRows = append(mapRows, rowData)
	}

	return mapRows, rows.Err()
}
