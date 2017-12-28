// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import (
	"github.com/s-gv/orangeforum/models/db"
	"time"
)

func CreateGroupMod(userName string, groupID string) {
	if uid, err := ReadUserIDByName(userName); err == nil {
		db.Exec(`INSERT INTO mods(userid, groupid, created_date) VALUES(?, ?, ?);`, uid, groupID, time.Now().Unix())
	}
}

func CreateGroupAdmin(userName string, groupID string) {
	if uid, err := ReadUserIDByName(userName); err == nil {
		db.Exec(`INSERT INTO admins(userid, groupid, created_date) VALUES(?, ?, ?);`, uid, groupID, time.Now().Unix())
	}
}

func ReadMods(groupID string) []string {
	rows := db.Query(`SELECT users.username FROM users INNER JOIN mods ON users.id=mods.userid WHERE mods.groupid=?;`, groupID)
	var mods []string
	for rows.Next() {
		var mod string
		rows.Scan(&mod)
		mods = append(mods, mod)
	}
	return mods
}

func ReadAdmins(groupID string) []string {
	rows := db.Query(`SELECT users.username FROM users INNER JOIN admins ON users.id=admins.userid WHERE admins.groupid=?;`, groupID)
	var admins []string
	for rows.Next() {
		var admin string
		rows.Scan(&admin)
		admins = append(admins, admin)
	}
	return admins
}

func IsUserGroupAdmin(userID string, groupID string) bool {
	r := db.QueryRow(`SELECT id FROM admins WHERE userid=? AND groupid=?`, userID, groupID)
	var tmp string
	if err := r.Scan(&tmp); err == nil {
		return true
	}
	return false
}

func ReadGroupIDByName(name string) string {
	r := db.QueryRow(`SELECT id FROM groups WHERE name=?;`, name)
	var id string
	if err := r.Scan(&id); err == nil {
		return id
	}
	return ""
}
