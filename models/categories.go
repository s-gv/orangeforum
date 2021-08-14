// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import (
	"database/sql"
	"time"

	"github.com/golang/glog"
)

type Category struct {
	DomainID     int          `db:"domain_id"`
	CategoryID   int          `db:"category_id"`
	Name         string       `db:"name"`
	Description  string       `db:"description"`
	HeaderMsg    string       `db:"header_msg"`
	NumTopics    int          `db:"num_topics"`
	IsPrivate    bool         `db:"is_private"`
	IsReadOnly   bool         `db:"is_readonly"`
	IsRestricted bool         `db:"is_restricted"`
	ArchivedAt   sql.NullTime `db:"archived_at"`
	CreatedAt    time.Time    `db:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at"`
}

func CreateCategory(domainID int, name string, description string) {
	_, err := DB.Exec("INSERT INTO categories(domain_id, name, description) VALUES($1, $2, $3);",
		domainID, name, description,
	)
	if err != nil {
		glog.Errorf("Error creating category: %s\n", err.Error())
	}
}

func GetCategoriesByDomainID(domainID int) []Category {
	categories := []Category{}
	err := DB.Select(&categories, "SELECT * FROM categories WHERE domain_id = $1 ORDER BY name;", domainID)
	if err != nil {
		glog.Errorf("Error listing categories: %s\n", err.Error())
	}
	return categories
}

func UpdateCategoryByID(categoryID int, name string, description string, isPrivate bool, isReadOnly bool, isArchived bool) {
	var archivedAt interface{}
	if isArchived {
		archivedAt = time.Now()
	}
	_, err := DB.Exec(`
		UPDATE categories SET 
			name = $2,
			description = $3,
			is_private = $4,
			is_readonly = $5,
			archived_at = $6
		WHERE
			category_id = $1;
		`,
		categoryID, name, description, isPrivate, isReadOnly, archivedAt,
	)

	if err != nil {
		glog.Errorf("Error updating category: %s\n", err.Error())
	}
}
