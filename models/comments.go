// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import (
	"database/sql"
	"html/template"
	"regexp"
	"strings"
	"time"

	"github.com/Depado/bfchroma"
	"github.com/alecthomas/chroma/styles"
	"github.com/golang/glog"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
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
	CommentID    int          `db:"comment_id"`
	TopicID      int          `db:"topic_id"`
	UserID       int          `db:"user_id"`
	Content      string       `db:"content"`
	IsSticky     bool         `db:"is_sticky"`
	ArchivedAt   sql.NullTime `db:"archived_at"`
	CreatedAt    time.Time    `db:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at"`
	DisplayName  string       `db:"display_name"`
	IsSuperAdmin bool         `db:"is_superadmin"`
	IsSuperMod   bool         `db:"is_supermod"`
}

var UGCPolicy = bluemonday.UGCPolicy()
var ChromaRenderer = bfchroma.NewRenderer(bfchroma.ChromaStyle(styles.GitHub))

func init() {
	UGCPolicy.AllowAttrs("style").Matching(regexp.MustCompile("^color:#[a-zA-Z0-9]+;background-color:#[a-zA-Z0-9]+$")).OnElements("pre")
	UGCPolicy.AllowAttrs("style").Matching(regexp.MustCompile("^color:#[a-zA-Z0-9]+(;font-weight:[A-Za-z0-9]+)?$")).OnElements("span")
}

func (c *CommentWithUser) CreatedAtStr() string {
	return RelTimeNowStr(c.CreatedAt)
}

func (c *CommentWithUser) UserIconColorStr() string {
	return UserIconColors[c.UserID%len(UserIconColors)]
}

func (c CommentWithUser) ContentRenderMarkdown() template.HTML {
	content := strings.ReplaceAll(c.Content, "\r\n", "\n")
	unsafe := blackfriday.Run([]byte(content), blackfriday.WithRenderer(ChromaRenderer))
	html := UGCPolicy.SanitizeBytes(unsafe)
	return template.HTML(string(html))
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
	_, err4 := DB.Exec("UPDATE users SET num_comments = (num_comments + 1) WHERE user_id = $1;", userID)
	if err4 != nil {
		glog.Errorf("Error updating comment count: %s\n", err2.Error())
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
		comments.*, users.display_name, users.is_superadmin, users.is_supermod
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

func DeleteCommentByID(commentID int, userID int, topicID int) {
	_, err := DB.Exec("DELETE FROM comments WHERE comment_id = $1;", commentID)
	if err != nil {
		glog.Errorf("Error deleting comment: %s\n", err.Error())
	}
	_, err2 := DB.Exec("UPDATE topics SET num_comments = (num_comments - 1) WHERE topic_id = $1;", topicID)
	if err2 != nil {
		glog.Errorf("Error updating comment count: %s\n", err2.Error())
	}
	_, err4 := DB.Exec("UPDATE users SET num_comments = (num_comments - 1) WHERE user_id = $1;", userID)
	if err4 != nil {
		glog.Errorf("Error updating comment count: %s\n", err2.Error())
	}
}
