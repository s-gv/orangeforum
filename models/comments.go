// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import (
	"database/sql"
	"time"

	"github.com/golang/glog"
)

type Comment struct {
	CommentID  int          `db:"comment_id"`
	TopicID    int          `db:"topic_id"`
	UserID     int          `db:"user_id"`
	Content    string       `db:"content"`
	IsSticky   bool         `db:"is_sticky"`
	ArchivedAt sql.NullTime `db:"archived_at"`
	CreatedAt  time.Time    `db:"created_at"`
	UpdatedAt  time.Time    `db:"updated_at"`
}

func CreateComment(topicID int, userID int, content string) int {
	var commentID int
	err := DB.QueryRow("INSERT INTO comments(topic_id, user_id, content) VALUES($1, $2, $3) RETURNING comment_id;",
		topicID, userID, content).Scan(&commentID)
	if err != nil {
		glog.Errorf("Error creating comment: %s\n", err.Error())
		return -1
	}
	return commentID
}

func GetCommentByID(commentID int) *Comment {
	comment := Comment{}
	err := DB.Get(&comment, "SELECT * FROM comments WHERE comment_id = $1;", commentID)
	if err != nil {
		if err != sql.ErrNoRows {
			glog.Errorf("Error reading comment: %s", err.Error())
		}
		return nil
	}
	return &comment
}

func GetCommentsByTopicID(topicID int) []Comment {
	var comments []Comment
	err := DB.Select(&comments, "SELECT * FROM comments WHERE topic_id = $1 ORDER BY created_at;", topicID)
	if err != nil {
		if err != sql.ErrNoRows {
			glog.Errorf("Error reading comments: %s\n", err.Error())
		}
	}
	return comments
}

func UpdateCommentByID(commentID int, content string) {
	_, err := DB.Exec("UPDATE comments SET content = $2 WHERE comment_id = $1;", commentID, content)
	if err != nil {
		glog.Errorf("Error updating comment: %s\n", err.Error())
	}
}

func DeleteCommentByID(commentID int) {
	_, err := DB.Exec("DELETE FROM comments WHERE comment_id = $1;", commentID)
	if err != nil {
		glog.Errorf("Error deleting comment: %s\n", err.Error())
	}
}
