// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import (
	"database/sql"
	"github.com/s-gv/orangeforum/models/db"
)

func NumUsers() int64 {
	r := db.QueryRow(`SELECT MAX(_ROWID_) FROM users LIMIT 1;`)
	var n sql.NullInt64
	if err := r.Scan(&n); err == nil {
		return n.Int64
	}
	return 0
}

func NumGroups() int64 {
	r := db.QueryRow(`SELECT MAX(_ROWID_) FROM groups LIMIT 1;`)
	var n sql.NullInt64
	if err := r.Scan(&n); err == nil {
		return n.Int64
	}
	return 0
}

func NumTopics() int64 {
	r := db.QueryRow(`SELECT MAX(_ROWID_) FROM topics LIMIT 1;`)
	var n sql.NullInt64
	if err := r.Scan(&n); err == nil {
		return n.Int64
	}
	return 0
}

func NumComments() int64 {
	r := db.QueryRow(`SELECT MAX(_ROWID_) FROM comments LIMIT 1;`)
	var n sql.NullInt64
	if err := r.Scan(&n); err == nil {
		return n.Int64
	}
	return 0
}
