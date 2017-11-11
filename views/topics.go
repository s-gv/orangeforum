// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package views

import (
	"github.com/s-gv/orangeforum/models/db"
	"time"
	"github.com/s-gv/orangeforum/templates"
	"github.com/s-gv/orangeforum/models"
	"net/http"
	"github.com/s-gv/orangeforum/utils"
	"strconv"
	"html/template"
)

var numCommentsPerPage = 100

var TopicIndexHandler = UA(func(w http.ResponseWriter, r *http.Request, sess Session) {
	topicID := r.FormValue("id")
	page64, err := strconv.ParseInt(r.FormValue("p"), 10, 64)
	if err != nil {
		page64 = 0
	}
	page := int(page64)
	if page < 0 {
		page = 0
	}
	var title, content, groupID, groupName string
	var isDeleted, isClosed bool
	var ownerID, createdDate int64
	if db.QueryRow(`SELECT title, content, userid, groupid, is_deleted, is_closed, created_date FROM topics WHERE id=?;`, topicID).Scan(
		&title, &content, &ownerID, &groupID, &isDeleted, &isClosed, &createdDate) != nil {
		ErrNotFoundHandler(w, r)
		return
	}
	if isDeleted {
		ErrNotFoundHandler(w, r)
		return
	}
	var ownerName string
	db.QueryRow(`SELECT username FROM users WHERE id=?;`, ownerID).Scan(&ownerName)

	subToken := ""
	if sess.UserID.Valid {
		db.QueryRow(`SELECT token FROM topicsubscriptions WHERE topicid=? AND userid=?;`, topicID, sess.UserID).Scan(&subToken)
	}

	var lastPos int
	db.QueryRow(`SELECT pos FROM comments WHERE topicid=? ORDER BY pos DESC LIMIT 1;`, topicID).Scan(&lastPos)
	isLastPage := (lastPos < (page+1)*numCommentsPerPage)
	isFirstPage := (page == 0)
	numPages := 0
	if lastPos > 0 {
		numPages = lastPos / numCommentsPerPage
	}

	type Comment struct {
		ID string
		Content template.HTML
		ImgSrc string
		CreatedDate string
		UserName string
		IsOwner bool
		IsDeleted bool
	}

	var comments []Comment
	var cDate int64
	var rows *db.Rows
	if page == 0 {
		rows = db.Query(`SELECT users.id, users.username, comments.id, comments.content, comments.image, comments.is_deleted, comments.created_date FROM comments INNER JOIN users ON comments.userid=users.id AND comments.topicid=? AND comments.pos < ? ORDER BY comments.pos;`, topicID, numCommentsPerPage)
	} else {
		rows = db.Query(`SELECT users.id, users.username, comments.id, comments.content, comments.image, comments.is_deleted, comments.created_date FROM comments INNER JOIN users ON comments.userid=users.id AND comments.topicid=? AND comments.pos >= ? AND comments.pos < ? ORDER BY comments.pos;`, topicID, page*numCommentsPerPage, (page+1)*numCommentsPerPage)
	}
	for rows.Next() {
		comments = append(comments, Comment{})
		c := &comments[len(comments)-1]
		var ownerID int64
		var content string
		rows.Scan(&ownerID, &c.UserName, &c.ID, &content, &c.ImgSrc, &c.IsDeleted, &cDate)
		c.CreatedDate = timeAgoFromNow(time.Unix(cDate, 0))
		c.IsOwner = sess.UserID.Valid && (ownerID == sess.UserID.Int64)
		c.Content = formatComment(content)
	}

	var tmp string
	db.QueryRow(`SELECT name FROM groups WHERE id=?;`, groupID).Scan(&groupName)
	isMod := db.QueryRow(`SELECT id FROM mods WHERE groupid=? AND userid=?;`, groupID, sess.UserID).Scan(&tmp) == nil
	isAdmin := db.QueryRow(`SELECT id FROM admins WHERE groupid=? AND userid=?;`, groupID, sess.UserID).Scan(&tmp) == nil
	isSuperAdmin := false
	db.QueryRow(`SELECT is_superadmin FROM users WHERE id=?`, sess.UserID).Scan(&isSuperAdmin)
	isOwner := sess.UserID.Valid && ownerID == sess.UserID.Int64

	templates.Render(w, "topicindex.html", map[string]interface{}{
		"Common": readCommonData(r, sess),
		"GroupID": groupID,
		"TopicID": topicID,
		"GroupName": groupName,
		"TopicName": title,
		"OwnerName": ownerName,
		"CreatedDate": timeAgoFromNow(time.Unix(createdDate, 0)),
		"SubToken": subToken,
		"Title": title,
		"Content": formatComment(content),
		"IsClosed": isClosed,
		"IsOwner": isOwner,
		"IsMod": isMod,
		"IsAdmin": isAdmin,
		"IsSuperAdmin": isSuperAdmin,
		"IsImageUploadEnabled": models.Config(models.ImageUploadEnabled) != "0",
		"Comments": comments,
		"IsFirstPage": isFirstPage,
		"IsLastPage": isLastPage,
		"NextPage": page+1,
		"CurrentPage": page,
		"Pages": make([]int, numPages),
	})
})

var TopicCreateHandler = A(func(w http.ResponseWriter, r *http.Request, sess Session) {
	groupID := r.FormValue("gid")
	var groupName string
	isGroupClosed := 1
	db.QueryRow(`SELECT name, is_closed FROM groups WHERE id=?;`, groupID).Scan(&groupName, &isGroupClosed)
	if isGroupClosed == 1 {
		ErrForbiddenHandler(w, r)
		return
	}

	var tmp int
	isMod := db.QueryRow(`SELECT id FROM mods WHERE groupid=? AND userid=?;`, groupID, sess.UserID).Scan(&tmp) == nil
	isAdmin := db.QueryRow(`SELECT id FROM admins WHERE groupid=? AND userid=?;`, groupID, sess.UserID).Scan(&tmp) == nil
	isSuperAdmin := false
	db.QueryRow(`SELECT is_superadmin FROM users WHERE id=?`, sess.UserID).Scan(&isSuperAdmin)

	if r.Method == "POST" {
		title := r.PostFormValue("title")
		content := r.PostFormValue("content")
		isSticky := r.PostFormValue("is_sticky") != ""
		if len(title) < 8 || len(title) > 80 {
			sess.SetFlashMsg("Title should have 8-80 characters.")
			http.Redirect(w, r, "/topics/new?gid="+groupID, http.StatusSeeOther)
			return
		}
		if len(content) > 5000 {
			sess.SetFlashMsg("Content should have less than 5000 characters.")
			http.Redirect(w, r, "/topics/new?gid="+groupID, http.StatusSeeOther)
			return
		}
		db.Exec(`INSERT INTO topics(title, content, userid, groupid, is_sticky, created_date, updated_date, activity_date) VALUES(?, ?, ?, ?, ?, ?, ?, ?);`,
			title, content, sess.UserID, groupID, isSticky, int(time.Now().Unix()), int(time.Now().Unix()), int(time.Now().Unix()))

		if models.Config(models.AllowGroupSubscription) != "0" {
			groupURL := "http://" + r.Host + "/groups?name=" + groupName
			rows := db.Query(`SELECT users.email, groupsubscriptions.token FROM users INNER JOIN groupsubscriptions ON users.id=groupsubscriptions.userid AND groupsubscriptions.groupid=?;`, groupID)
			for rows.Next() {
				var email, token string
				rows.Scan(&email, &token)
				if email != "" {
					unSubURL := "http://" + r.Host + "/groups/unsubscribe?token=" + token
					utils.SendMail(email, `New topic in `+groupName,
						"A new topic titled \""+title+"\" has been posted to "+groupName+".\r\nSee topics posted to the group at "+groupURL+"\r\n\r\nIf you do not want these emails, unsubscribe by following this link: "+unSubURL)
				}
			}
		}
		http.Redirect(w, r, "/groups?name="+groupName, http.StatusSeeOther)
		return
	}

	templates.Render(w, "topicedit.html", map[string]interface{}{
		"Common":    readCommonData(r, sess),
		"GroupID":   groupID,
		"GroupName": groupName,
		"TopicID":   "",
		"Title":     "",
		"Content":   "",
		"IsSticky": false,
		"IsClosed": false,
		"IsDeleted": false,
		"IsMod": isMod,
		"IsAdmin": isAdmin,
		"IsSuperAdmin": isSuperAdmin,
	})
})

var TopicUpdateHandler = A(func(w http.ResponseWriter, r *http.Request, sess Session) {
	topicID := r.FormValue("id")
	groupID := ""
	title := r.PostFormValue("title")
	content := r.PostFormValue("content")
	action := r.PostFormValue("action")
	isSticky := r.PostFormValue("is_sticky") != ""
	isClosed := true
	isDeleted := true

	if db.QueryRow(`SELECT groupid FROM topics WHERE id=?;`, topicID).Scan(&groupID) != nil {
		ErrNotFoundHandler(w, r)
		return
	}

	isGroupClosed := 1
	var groupName string
	db.QueryRow(`SELECT name, is_closed FROM groups WHERE id=?;`, groupID).Scan(&groupName, &isGroupClosed)
	if isGroupClosed == 1 {
		ErrForbiddenHandler(w, r)
		return
	}

	var tmp int
	var uID int64
	db.QueryRow(`SELECT userid FROM topics WHERE id=?;`, topicID).Scan(&uID)

	isOwner := (uID == sess.UserID.Int64)
	isMod := db.QueryRow(`SELECT id FROM mods WHERE groupid=? AND userid=?;`, groupID, sess.UserID).Scan(&tmp) == nil
	isAdmin := db.QueryRow(`SELECT id FROM admins WHERE groupid=? AND userid=?;`, groupID, sess.UserID).Scan(&tmp) == nil
	isSuperAdmin := false
	db.QueryRow(`SELECT is_superadmin FROM users WHERE id=?`, sess.UserID).Scan(&isSuperAdmin)

	if !isMod && !isAdmin && !isSuperAdmin {
		db.QueryRow(`SELECT is_sticky FROM topics WHERE id=?;`, topicID).Scan(&isSticky)
		if !isOwner {
			ErrForbiddenHandler(w, r)
			return
		}
	}

	if r.Method == "POST" {
		if len(title) < 8 || len(title) > 80 {
			sess.SetFlashMsg("Title should have 8-80 characters.")
			http.Redirect(w, r, "/topics/edit?id="+topicID, http.StatusSeeOther)
			return
		}
		if len(content) > 5000 {
			sess.SetFlashMsg("Content should have less than 5000 characters.")
			http.Redirect(w, r, "/topics/edit?id="+topicID, http.StatusSeeOther)
			return
		}
		if action == "Update" {
			db.Exec(`UPDATE topics SET title=?, content=?, is_sticky=?, updated_date=? WHERE id=?;`, title, content, isSticky, int(time.Now().Unix()), topicID)
		} else if action == "Close" && (isMod || isAdmin || isSuperAdmin) {
			db.Exec(`UPDATE topics SET is_closed=1 WHERE id=?;`, topicID)
		} else if action == "Reopen" && (isMod || isAdmin || isSuperAdmin) {
			db.Exec(`UPDATE topics SET is_closed=0 WHERE id=?;`, topicID)
		} else if action == "Delete" {
			db.Exec(`UPDATE topics SET is_deleted=1 WHERE id=?;`, topicID)
			http.Redirect(w, r, "/topics/edit?id="+topicID, http.StatusSeeOther)
			return
		} else if action == "Undelete" {
			db.Exec(`UPDATE topics SET is_deleted=0 WHERE id=?;`, topicID)
		}
		http.Redirect(w, r, "/topics?id="+topicID, http.StatusSeeOther)
		return
	}

	if db.QueryRow(`SELECT title, content, is_sticky, is_deleted, is_closed FROM topics WHERE id=?;`, topicID).Scan(&title, &content, &isSticky, &isDeleted, &isClosed) != nil {
		ErrNotFoundHandler(w, r)
		return
	}

	templates.Render(w, "topicedit.html", map[string]interface{}{
		"Common": readCommonData(r, sess),
		"GroupID": groupID,
		"GroupName": groupName,
		"TopicID": topicID,
		"Title": title,
		"Content":      content,
		"IsSticky":     isSticky,
		"IsClosed":    isClosed,
		"IsDeleted": isDeleted,
		"IsMod":        isMod,
		"IsAdmin":      isAdmin,
		"IsSuperAdmin": isSuperAdmin,
	})
})

var TopicSubscribeHandler = A(func(w http.ResponseWriter, r *http.Request, sess Session) {
	topicID := r.FormValue("id")
	if models.Config(models.AllowTopicSubscription) == "0" {
		ErrForbiddenHandler(w, r)
		return
	}
	var tmp string
	if db.QueryRow(`SELECT id FROM topics WHERE id=?;`, topicID).Scan(&tmp) != nil {
		ErrNotFoundHandler(w, r)
		return
	}
	if r.Method == "POST" {
		var tmp string
		if db.QueryRow(`SELECT id FROM topicsubscriptions WHERE userid=? AND topicid=?;`, sess.UserID, topicID).Scan(&tmp) != nil {
			db.Exec(`INSERT INTO topicsubscriptions(userid, topicid, token, created_date) VALUES(?, ?, ?, ?);`,
				sess.UserID, topicID, randSeq(64), time.Now().Unix())
		}
	}
	http.Redirect(w, r, "/topics?id="+topicID, http.StatusSeeOther)
})

var TopicUnsubscribeHandler = UA(func(w http.ResponseWriter, r *http.Request, sess Session) {
	token := r.FormValue("token")
	var topicID, topicName string
	if db.QueryRow(`SELECT topicid FROM topicsubscriptions WHERE token=?;`, token).Scan(&topicID) != nil {
		ErrNotFoundHandler(w, r)
		return
	}
	db.QueryRow(`SELECT title FROM topics WHERE id=?;`, topicID).Scan(&topicName)
	if r.Method == "POST" {
		db.Exec(`DELETE FROM topicsubscriptions WHERE token=?;`, token)
		if r.PostFormValue("noredirect") != "" {
			w.Write([]byte("Unsubscribed."))
		} else {
			http.Redirect(w, r, "/topics?id="+topicID, http.StatusSeeOther)
		}
		return
	}
	w.Write([]byte(`<!DOCTYPE html><html><head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1"></head>
	<body><form action="/topics/unsubscribe" method="POST">
	Unsubscribe from `+ topicName +`?
	<input type="hidden" name="token" value="`+token+`">
	<input type="hidden" name="csrf" value="`+sess.CSRFToken+`">
	<input type="hidden" name="noredirect" value="1">
	<input type="submit" value="Unsubscribe">
	</form></body></html>`))
})
