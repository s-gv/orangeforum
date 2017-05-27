package models

import (
	"database/sql"
	"errors"
)

var ErrDBVer = errors.New("DB version not up-to-date. Migration needed.")
var ErrDBMigrationNotNeeded = errors.New("DB version is up-to-date.")
var ErrDBVerAhead = errors.New("DB written by a newer version.")

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

