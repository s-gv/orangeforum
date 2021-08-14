// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import (
	"database/sql"
	"time"
)

type Comment struct {
	CommentID  int          `db:"user_id"`
	TopicID    int          `db:"topic_id"`
	UserID     int          `db:"user_id"`
	Content    string       `db:"content"`
	IsSticky   bool         `db:"is_sticky"`
	ArchivedAt sql.NullTime `db:"archived_at"`
	CreatedAt  time.Time    `db:"created_at"`
	UpdatedAt  time.Time    `db:"updated_at"`
}
