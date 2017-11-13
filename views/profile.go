// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package views

import (
	"net/http"
	"github.com/s-gv/orangeforum/models/db"
	"strconv"
	"time"
	"github.com/s-gv/orangeforum/templates"
	"html/template"
)

var UserProfileHandler = UA(func(w http.ResponseWriter, r *http.Request, sess Session) {
	userName := r.FormValue("u")
	var about, email string
	var isBanned bool
	var userID int64
	if db.QueryRow(`SELECT id, about, email, is_banned FROM users WHERE username=?;`, userName).Scan(&userID, &about, &email, &isBanned) != nil {
		ErrNotFoundHandler(w, r)
		return
	}

	templates.Render(w, "profile.html", map[string]interface{}{
		"Common": readCommonData(r, sess),
		"UserName": userName,
		"About": about,
		"Email": email,
		"IsSelf": sess.UserID.Valid && (userID == sess.UserID.Int64),
		"IsBanned": isBanned,
	})
})

var UserProfileUpdateHandler = A(func(w http.ResponseWriter, r *http.Request, sess Session) {
	userName := r.FormValue("u")
	var about, email string
	var isBanned bool
	var userID int64
	if db.QueryRow(`SELECT id, about, email, is_banned FROM users WHERE username=?;`, userName).Scan(&userID, &about, &email, &isBanned) != nil {
		ErrNotFoundHandler(w, r)
		return
	}

	if r.Method == "POST" {
		if !sess.UserID.Valid {
			ErrForbiddenHandler(w, r)
			return
		}
		action := r.PostFormValue("action")
		var isSuperAdmin bool
		db.QueryRow(`SELECT is_superadmin FROM users WHERE id=?;`, sess.UserID).Scan(&isSuperAdmin)
		if action == "Update" {
			if isSuperAdmin || userID == sess.UserID.Int64 {
				email := r.FormValue("email")
				about := r.FormValue("about")
				if len(email) > 64 {
					sess.SetFlashMsg("Email should have fewer than 64 characters.")
					http.Redirect(w, r, "/users?u="+userName, http.StatusSeeOther)
					return
				}
				if len(about) > 1024 {
					sess.SetFlashMsg("About should have fewer than 1024 characters.")
					http.Redirect(w, r, "/users?u="+userName, http.StatusSeeOther)
					return
				}
				db.Exec(`UPDATE users SET email=?, about=? WHERE id=?;`, email, about, userID)
			} else {
				ErrForbiddenHandler(w, r)
				return
			}
		} else if action == "Ban" {
			if isSuperAdmin {
				db.Exec(`UPDATE users SET is_banned=1 WHERE id=?;`, userID)
				db.Exec(`DELETE FROM sessions WHERE userid=?;`, userID)
			} else {
				ErrForbiddenHandler(w, r)
				return
			}
		} else if action == "Unban" {
			if isSuperAdmin {
				db.Exec(`UPDATE users SET is_banned=0 WHERE id=?;`, userID)
			} else {
				ErrForbiddenHandler(w, r)
				return
			}
		}
	}
	sess.SetFlashMsg("Update successful.")
	http.Redirect(w, r, "/users?u="+userName, http.StatusSeeOther)
})

var UserCommentsHandler = UA(func(w http.ResponseWriter, r *http.Request, sess Session) {
	ownerName := r.FormValue("u")
	lastCommentDate, err := strconv.ParseInt(r.FormValue("lcd"), 10, 64)

	if err != nil {
		lastCommentDate = 0
	}

	var ownerID string
	if db.QueryRow(`SELECT id FROM users WHERE username=?;`, ownerName).Scan(&ownerID) != nil {
		ErrNotFoundHandler(w, r)
		return
	}

	type Comment struct {
		ID string
		Content template.HTML
		TopicID string
		TopicName string
		CreatedDate string
		ImgSrc string
		IsDeleted bool
	}

	commentsPerPage := 50

	var comments []Comment
	var rows *db.Rows
	if lastCommentDate == 0 {
		rows = db.Query(`SELECT topics.title, comments.topicid, comments.id, comments.content, comments.image, comments.created_date, comments.is_deleted FROM comments INNER JOIN topics ON topics.id = comments.topicid AND comments.userid=? ORDER BY comments.created_date DESC LIMIT ?;`, ownerID, commentsPerPage)
	} else {
		rows = db.Query(`SELECT topics.title, comments.topicid, comments.id, comments.content, comments.image, comments.created_date, comments.is_deleted FROM comments INNER JOIN topics ON topics.id = comments.topicid AND comments.userid=? AND comments.created_date < ? ORDER BY comments.created_date DESC LIMIT ?;`, ownerID, lastCommentDate, commentsPerPage)
	}

	var cDate int64
	for rows.Next() {
		comments = append(comments, Comment{})
		c := &comments[len(comments)-1]

		var content string
		rows.Scan(&c.TopicName, &c.TopicID, &c.ID, &content, &c.ImgSrc, &cDate, &c.IsDeleted)
		c.CreatedDate = timeAgoFromNow(time.Unix(cDate, 0))
		c.Content = formatComment(content)
	}

	if len(comments) >= commentsPerPage {
		lastCommentDate = cDate
	} else {
		lastCommentDate = 0
	}

	templates.Render(w, "profilecomments.html", map[string]interface{}{
		"Common": readCommonData(r, sess),
		"OwnerName": ownerName,
		"Comments": comments,
		"LastCommentDate": lastCommentDate,
	})
})

var UserTopicsHandler = UA(func(w http.ResponseWriter, r *http.Request, sess Session) {
	ownerName := r.FormValue("u")
	var ownerID string
	if db.QueryRow(`SELECT id FROM users WHERE username=?;`, ownerName).Scan(&ownerID) != nil {
		ErrNotFoundHandler(w, r)
		return
	}
	lastTopicDate, err := strconv.ParseInt(r.FormValue("ltd"), 10, 64)
	if err != nil {
		lastTopicDate = 0
	}

	numTopicsPerPage := 50
	type Topic struct {
		ID string
		Title string
		IsClosed bool
		IsDeleted bool
		CreatedDate string
	}
	var topics []Topic
	var rows *db.Rows
	var cDate int64
	if lastTopicDate == 0 {
		rows = db.Query(`SELECT id, title, is_deleted, is_closed, created_date FROM topics WHERE userid=? ORDER BY created_date DESC LIMIT ?;`, ownerID, numTopicsPerPage)
	} else {
		rows = db.Query(`SELECT id, title, is_deleted, is_closed, created_date FROM topics WHERE userid=? AND created_date < ? ORDER BY created_date DESC LIMIT ?;`, ownerID, lastTopicDate, numTopicsPerPage)
	}
	for rows.Next() {
		topics = append(topics, Topic{})
		t := &topics[len(topics)-1]
		rows.Scan(&t.ID, &t.Title, &t.IsDeleted, &t.IsClosed, &cDate)
		t.CreatedDate = timeAgoFromNow(time.Unix(cDate, 0))
		t.Title = censor(t.Title)
	}

	if len(topics) >= numTopicsPerPage {
		lastTopicDate = cDate
	} else {
		lastTopicDate = 0
	}

	templates.Render(w, "profiletopics.html", map[string]interface{}{
		"Common": readCommonData(r, sess),
		"OwnerName": ownerName,
		"Topics": topics,
		"LastTopicDate": lastTopicDate,
	})
})

var UserGroupsHandler = A(func(w http.ResponseWriter, r *http.Request, sess Session) {
	ownerID := sess.UserID.Int64
	var ownerName string

	type Group struct {
		ID string
		Name string
		IsClosed bool
		CreatedDate string
	}
	var adminInGroups []Group
	rows := db.Query(`SELECT groups.id, groups.name, groups.is_closed, groups.created_date FROM groups INNER JOIN admins ON admins.groupid=groups.id AND admins.userid=?;`, ownerID)
	for rows.Next() {
		adminInGroups = append(adminInGroups, Group{})
		g := &adminInGroups[len(adminInGroups)-1]
		var cDate int64
		rows.Scan(&g.ID, &g.Name, &g.IsClosed, &cDate)
		g.CreatedDate = timeAgoFromNow(time.Unix(cDate, 0))
	}

	var modInGroups []Group
	rows = db.Query(`SELECT groups.id, groups.name, groups.is_closed, groups.created_date FROM groups INNER JOIN mods ON mods.groupid=groups.id AND mods.userid=?;`, ownerID)
	for rows.Next() {
		modInGroups = append(modInGroups, Group{})
		g := &adminInGroups[len(adminInGroups)-1]
		var cDate int64
		rows.Scan(&g.ID, &g.Name, &g.IsClosed, &cDate)
		g.CreatedDate = timeAgoFromNow(time.Unix(cDate, 0))
	}

	templates.Render(w, "profilegroups.html", map[string]interface{}{
		"Common": readCommonData(r, sess),
		"OwnerName": ownerName,
		"AdminInGroups": adminInGroups,
		"ModInGroups": modInGroups,
	})
})
