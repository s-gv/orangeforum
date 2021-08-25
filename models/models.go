// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import (
	"github.com/golang/glog"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func CleanDB() {
	if _, err := DB.Exec(`DELETE FROM domains;`); err != nil {
		glog.Errorf("Error deleting domains: %s\n", err.Error())
	}
	if _, err := DB.Exec(`DELETE FROM users;`); err != nil {
		glog.Errorf("Error deleting users: %s\n", err.Error())
	}
	if _, err := DB.Exec(`DELETE FROM categories;`); err != nil {
		glog.Errorf("Error deleting categories: %s\n", err.Error())
	}
	if _, err := DB.Exec(`DELETE FROM topics;`); err != nil {
		glog.Errorf("Error deleting topics: %s\n", err.Error())
	}
	if _, err := DB.Exec(`DELETE FROM comments;`); err != nil {
		glog.Errorf("Error deleting comments: %s\n", err.Error())
	}
	if _, err := DB.Exec(`DELETE FROM category_subs;`); err != nil {
		glog.Errorf("Error deleting category_subs: %s\n", err.Error())
	}
	if _, err := DB.Exec(`DELETE FROM topic_subs;`); err != nil {
		glog.Errorf("Error deleting topic_subs: %s\n", err.Error())
	}
	if _, err := DB.Exec(`DELETE FROM notes;`); err != nil {
		glog.Errorf("Error deleting notes: %s\n", err.Error())
	}
}
