package models

import (
	"log"
)

const ModelVersion = 1

func RunMigrationZero() {
	createConfigTable()
	WriteConfig("version", "1")
}

func Migrate(driverName string, dataSourceName string) error {
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

