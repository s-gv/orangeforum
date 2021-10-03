// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import (
	"errors"
)

const CurrentDBVersion = 2

func Migrate() error {
	iver := GetDBVersion()
	if iver > CurrentDBVersion {
		return errors.New("Database schema version newer than this binary. Get the latest version of OrangeForum.")
	}
	if iver == CurrentDBVersion {
		return errors.New("Database schema is already up-to-date.")
	}
	if iver < 1 {
		migrate001(DB)
	}
	if iver < 2 {
		migrate002(DB)
	}
	return nil
}

func IsMigrationNeeded() error {
	iver := GetDBVersion()
	if iver < CurrentDBVersion {
		return errors.New("Database migration needed.")
	}
	if iver > CurrentDBVersion {
		return errors.New("Database schema version newer than this binary. Get the latest version of OrangeForum.")
	}
	return nil
}
