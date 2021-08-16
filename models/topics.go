// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import (
	"database/sql"
	"time"

	"github.com/golang/glog"
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

func GetTopicsByCategoryID(categoryID int, before time.Time) []Topic {
	var topics []Topic
	err := DB.Select(&topics, `
		SELECT * 
		FROM topics 
		WHERE category_id = $1 AND activity_at < $2
		ORDER BY is_sticky DESC, activity_at DESC LIMIT 30;`,
		categoryID, before,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			glog.Errorf("Error reading topics: %s\n", err.Error())
		}
	}
	return topics
}

func CreateTopic(categoryID int, userID int, title string, content string, isSticky bool) int {
	var id int
	err := DB.QueryRow("INSERT INTO topics(category_id, user_id, title, content, is_sticky) VALUES($1, $2, $3, $4, $5) RETURNING topic_id;",
		categoryID, userID, title, content, isSticky).Scan(&id)
	if err != nil {
		glog.Errorf("Error inserting row: %s\n", err.Error())
		return -1
	}
	return id
}

func GetTopicByID(topicID int) *Topic {
	topic := Topic{}
	err := DB.Get(&topic, "SELECT * FROM topics WHERE topic_id = $1;", topicID)
	if err != nil && err != sql.ErrNoRows {
		glog.Errorf("Error reading topic: %s\n", err.Error())
	}
	return &topic
}

func UpdateTopicByID(topicID int, title string, content string, isSticky bool) {
	_, err := DB.Exec("UPDATE topics SET title = $2, content = $3, is_sticky = $4 WHERE topic_id = $1;", topicID, title, content, isSticky)
	if err != nil {
		glog.Errorf("Error updating topic: %s\n", err.Error())
	}
}

func DeleteTopicByID(topicID int) {
	_, err := DB.Exec("DELETE FROM topics WHERE topic_id = $1;", topicID)
	if err != nil {
		glog.Errorf("Error updating topic: %s\n", err.Error())
	}
}
