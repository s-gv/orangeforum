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

type CommentWithUser struct {
	CommentID   int          `db:"comment_id"`
	TopicID     int          `db:"topic_id"`
	UserID      int          `db:"user_id"`
	Content     string       `db:"content"`
	IsSticky    bool         `db:"is_sticky"`
	ArchivedAt  sql.NullTime `db:"archived_at"`
	CreatedAt   time.Time    `db:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at"`
	DisplayName string       `db:"display_name"`
}

func (c *CommentWithUser) CreatedAtStr() string {
	return RelTimeNowStr(c.CreatedAt)
}

func (c *CommentWithUser) UserIconColorStr() string {
	return UserIconColors[c.UserID%len(UserIconColors)]
}

func CreateComment(topicID int, userID int, content string, isSticky bool) int {
	var commentID int
	err := DB.QueryRow("INSERT INTO comments(topic_id, user_id, content, is_sticky) VALUES($1, $2, $3, $4) RETURNING comment_id;",
		topicID, userID, content, isSticky).Scan(&commentID)
	if err != nil {
		glog.Errorf("Error creating comment: %s\n", err.Error())
		return -1
	}
	_, err2 := DB.Exec("UPDATE topics SET num_comments = (num_comments + 1) WHERE topic_id = $1;", topicID)
	if err2 != nil {
		glog.Errorf("Error updating comment count: %s\n", err2.Error())
	}
	_, err3 := DB.Exec("UPDATE topics SET activity_at = current_timestamp WHERE topic_id = $1;", topicID)
	if err3 != nil {
		glog.Errorf("Error updating activity time: %s\n", err3.Error())
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

func GetCommentsByTopicID(topicID int) []CommentWithUser {
	var comments []CommentWithUser
	err := DB.Select(&comments, `
	SELECT 
		comments.*, users.display_name
	FROM 
		comments INNER JOIN users ON comments.user_id = users.user_id
	WHERE
		comments.topic_id = $1 ORDER BY comments.is_sticky DESC, comments.created_at;`,
		topicID)
	if err != nil {
		if err != sql.ErrNoRows {
			glog.Errorf("Error reading comments: %s\n", err.Error())
		}
	}
	return comments
}

func UpdateCommentByID(commentID int, content string, isSticky bool) {
	_, err := DB.Exec("UPDATE comments SET content = $2, is_sticky = $3 WHERE comment_id = $1;", commentID, content, isSticky)
	if err != nil {
		glog.Errorf("Error updating comment: %s\n", err.Error())
	}
}

func DeleteCommentByID(commentID int, topicID int) {
	_, err := DB.Exec("DELETE FROM comments WHERE comment_id = $1;", commentID)
	if err != nil {
		glog.Errorf("Error deleting comment: %s\n", err.Error())
	}
	_, err2 := DB.Exec("UPDATE topics SET num_comments = (num_comments - 1) WHERE topic_id = $1;", topicID)
	if err2 != nil {
		glog.Errorf("Error updating comment count: %s\n", err2.Error())
	}
}
