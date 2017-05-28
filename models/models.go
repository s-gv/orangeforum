package models

import (
	"database/sql"
)


func Init(driverName string, dataSourceName string) error {
	mydb, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}
	db = mydb

	dbver := DBVersion()
	if dbver < ModelVersion {
		return ErrDBVer
	}
	return nil
}

