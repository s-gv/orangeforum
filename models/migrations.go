package models

import (
	"log"
	"errors"
)

const ModelVersion = 1

var ErrDBVer = errors.New("DB version not up-to-date. Migration needed.")
var ErrDBMigrationNotNeeded = errors.New("DB version is up-to-date.")
var ErrDBVerAhead = errors.New("DB written by a newer version.")

func RunMigrationZero() {
	db.Exec(`CREATE TABLE config(key TEXT, val TEXT);`)
	db.Exec(`CREATE UNIQUE INDEX key_index on config(key);`)

	WriteConfig("version", "1")
}

func Migrate() error {
	dbver := DBVersion()
	if dbver > ModelVersion {
		return ErrDBVerAhead
	} else if dbver == ModelVersion {
		return ErrDBMigrationNotNeeded
	}

	for dbver < ModelVersion {
		switch dbver {
		case 0:
			RunMigrationZero()
		}
		newDBVer := DBVersion()
		if newDBVer != dbver + 1 {
			log.Fatal("Migration failed ", dbver)
		}
		dbver = newDBVer
	}
	return nil
}

