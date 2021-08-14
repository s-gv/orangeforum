// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import (
	"database/sql"
	"time"
)

type Topic struct {
	TopicID     int          `db:"topic_id"`
	CategoryID  int          `db:"category_id"`
	UserID      int          `db:"user_id"`
	Title       string       `db:"title"`
	Content     string       `db:"content"`
	IsSticky    bool         `db:"is_sticky"`
	IsReadOnly  bool         `db:"is_readonly"`
	NumComments int          `db:"num_comments"`
	NumViews    int          `db:"num_views"`
	ActivityAt  time.Time    `db:"activity_at"`
	ArchivedAt  sql.NullTime `db:"archived_at"`
	CreatedAt   time.Time    `db:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at"`
}
