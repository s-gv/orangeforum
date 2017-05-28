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
	createConfigTable()
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

