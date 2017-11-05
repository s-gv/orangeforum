// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package views

import (
	"net/http"
	"github.com/s-gv/orangeforum/models"
	"github.com/s-gv/orangeforum/models/db"
	"github.com/s-gv/orangeforum/templates"
	"time"
	"database/sql"
	"github.com/s-gv/orangeforum/utils"
	"strconv"
)

var CommentIndexHandler = UA(func(w http.ResponseWriter, r *http.Request, sess Session) {
	commentID := r.FormValue("id")
	var groupID, topicID, topicName, groupName, ownerID, ownerName, content, imgSrc string
	var cDate int64
	var isDeleted bool

	if db.QueryRow(`SELECT userid, topicid, content, image, is_deleted, created_date FROM comments WHERE id=?;`, commentID).Scan(
		&ownerID, &topicID, &content, &imgSrc, &isDeleted, &cDate) != nil {
		ErrNotFoundHandler(w, r)
		return
	}
	db.QueryRow(`SELECT groupid, title FROM topics WHERE id=?;`, topicID).Scan(&groupID, &topicName)
	db.QueryRow(`SELECT username FROM users WHERE id=?;`, ownerID).Scan(&ownerName)
	db.QueryRow(`SELECT name FROM groups WHERE id=?;`, groupID).Scan(&groupName)

	var tmp string
	db.QueryRow(`SELECT name FROM groups WHERE id=?;`, groupID).Scan(&groupName)
	isMod := sess.UserID.Valid && db.QueryRow(`SELECT id FROM mods WHERE groupid=? AND userid=?;`, groupID, sess.UserID).Scan(&tmp) == nil
	isAdmin := sess.UserID.Valid && db.QueryRow(`SELECT id FROM admins WHERE groupid=? AND userid=?;`, groupID, sess.UserID).Scan(&tmp) == nil
	isSuperAdmin := false
	if sess.UserID.Valid {
		db.QueryRow(`SELECT is_superadmin FROM users WHERE id=?`, sess.UserID).Scan(&isSuperAdmin)
	}
	isOwner := sess.UserID.Valid && db.QueryRow(`SELECT userid FROM comments WHERE id=?;`, commentID).Scan(&tmp) == nil

	templates.Render(w, "commentindex.html", map[string]interface{}{
		"Common": readCommonData(r, sess),
		"ID": commentID,
		"TopicID": topicID,
		"TopicName": topicName,
		"GroupName": groupName,
		"OwnerName": ownerName,
		"Content": formatComment(content),
		"ImgSrc": imgSrc,
		"IsMod": isMod,
		"IsAdmin": isAdmin,
		"IsSuperAdmin": isSuperAdmin,
		"IsOwner": isOwner,
		"IsDeleted": isDeleted,
		"CreatedDate": timeAgoFromNow(time.Unix(cDate, 0)),
	})
})

var CommentCreateHandler = A(func(w http.ResponseWriter, r *http.Request, sess Session) {
	topicID := r.FormValue("tid")
	content := r.PostFormValue("content")
	isSticky := r.PostFormValue("is_sticky") != ""
	isImageUploadEnabled := models.Config(models.ImageUploadEnabled) != "0"
	var groupID, groupName, topicName, parentComment, topicOwnerID, topicOwnerName string
	var topicCreatedDate int64

	if db.QueryRow(`SELECT userid, groupid, title, content, created_date FROM topics WHERE id=?;`, topicID).Scan(
		&topicOwnerID, &groupID, &topicName, &parentComment, &topicCreatedDate) != nil {
		ErrNotFoundHandler(w, r)
		return
	}
	isClosed := true
	db.QueryRow(`SELECT is_closed FROM groups WHERE id=?;`, groupID).Scan(&isClosed)

	if isClosed {
		ErrForbiddenHandler(w, r)
		return
	}
	db.QueryRow(`SELECT username FROM users WHERE id=?;`, topicOwnerID).Scan(&topicOwnerName)

	var tmp string
	db.QueryRow(`SELECT name FROM groups WHERE id=?;`, groupID).Scan(&groupName)
	isMod := db.QueryRow(`SELECT id FROM mods WHERE groupid=? AND userid=?;`, groupID, sess.UserID).Scan(&tmp) == nil
	isAdmin := db.QueryRow(`SELECT id FROM admins WHERE groupid=? AND userid=?;`, groupID, sess.UserID).Scan(&tmp) == nil
	isSuperAdmin := false
	db.QueryRow(`SELECT is_superadmin FROM users WHERE id=?`, sess.UserID).Scan(&isSuperAdmin)

	if r.Method == "POST" {
		if content == "" {
			http.Redirect(w, r, "/topics?id="+topicID+"#comment-last", http.StatusSeeOther)
			return
		}
		if !isMod && !isAdmin && !isSuperAdmin {
			isSticky = false
		}
		imageName := ""
		if isImageUploadEnabled {
			imageName = saveImage(r)
		}

		var lastPos int
		db.QueryRow(`SELECT pos FROM comments WHERE topicid=? ORDER BY pos DESC LIMIT 1;`, topicID).Scan(&lastPos)
		newPos := lastPos + 1
		if isSticky {
			newPos = -newPos
		}

		db.Exec(`INSERT INTO comments(content, image, topicid, userid, parentid, pos, created_date, updated_date) VALUES(?, ?, ?, ?, ?, ?, ?, ?);`,
			content, imageName, topicID, sess.UserID, sql.NullInt64{Valid:false}, newPos, int64(time.Now().Unix()), int64(time.Now().Unix()))
		db.Exec(`UPDATE topics SET num_comments=num_comments+1, activity_date=? WHERE id=?;`, int(time.Now().Unix()), topicID)
		if models.Config(models.AllowTopicSubscription) != "0" {
			var userName string
			db.QueryRow(`SELECT username FROM users WHERE id=?;`, sess.UserID).Scan(&userName)
			topicURL := "http://" + r.Host + "/topics?id=" + topicID
			rows := db.Query(`SELECT users.email, topicsubscriptions.token FROM users INNER JOIN topicsubscriptions ON users.id=topicsubscriptions.userid AND topicsubscriptions.topicid=?;`, topicID)
			for rows.Next() {
				var email, token string
				rows.Scan(&email, &token)
				if email != "" {
					unSubURL := "http://" + r.Host + "/topics/unsubscribe?token=" + token
					utils.SendMail(email, `New comment in "`+topicName+`"`,
						"A new comment has been posted by "+userName+" in \""+topicName+"\".\r\nSee the comment at "+topicURL+"\r\n\r\nIf you do not want these emails, unsubscribe by following this link: "+unSubURL)
				}
			}
		}
		page := newPos / numCommentsPerPage
		if page < 0 {
			page = 0
		}
		http.Redirect(w, r, "/topics?id="+topicID+"&p="+strconv.Itoa(page)+"#comment-last", http.StatusSeeOther)
		return
	}

	templates.Render(w, "commentedit.html", map[string]interface{}{
		"Common": readCommonData(r, sess),
		"TopicID": topicID,
		"TopicOwnerName": topicOwnerName,
		"TopicCreatedDate": timeAgoFromNow(time.Unix(topicCreatedDate, 0)),
		"CommentID": "",
		"TopicName": topicName,
		"GroupName": groupName,
		"ParentComment": parentComment,
		"Content": "",
		"IsSticky": false,
		"IsMod": isMod,
		"IsAdmin": isAdmin,
		"IsSuperAdmin": isSuperAdmin,
		"IsImageUploadEnabled": isImageUploadEnabled,
	})
})

var CommentUpdateHandler = A(func(w http.ResponseWriter, r *http.Request, sess Session) {
	commentID := r.FormValue("id")
	content := r.PostFormValue("content")
	isSticky := r.PostFormValue("is_sticky") != ""

	var groupID, topicID, groupName, topicName, parentComment, topicOwnerName, topicOwnerID string
	var topicCreatedDate int64
	var pos int
	if db.QueryRow(`SELECT topicid, pos FROM comments WHERE id=?;`, commentID).Scan(&topicID, &pos) != nil {
		ErrNotFoundHandler(w, r)
		return
	}
	if db.QueryRow(`SELECT userid, groupid, title, content, created_date FROM topics WHERE id=?;`, topicID).Scan(
		&topicOwnerID, &groupID, &topicName, &parentComment, &topicCreatedDate) != nil {
		ErrNotFoundHandler(w, r)
		return
	}
	isClosed := true
	db.QueryRow(`SELECT is_closed FROM groups WHERE id=?;`, groupID).Scan(&isClosed)
	if !isClosed {
		db.QueryRow(`SELECT is_closed FROM topics WHERE id=?;`, topicID).Scan(&isClosed)
	}

	if isClosed {
		ErrForbiddenHandler(w, r)
		return
	}

	db.QueryRow(`SELECT username FROM users WHERE id=?;`, topicOwnerID).Scan(&topicOwnerName)

	var tmp string
	db.QueryRow(`SELECT name FROM groups WHERE id=?;`, groupID).Scan(&groupName)
	isMod := db.QueryRow(`SELECT id FROM mods WHERE groupid=? AND userid=?;`, groupID, sess.UserID).Scan(&tmp) == nil
	isAdmin := db.QueryRow(`SELECT id FROM admins WHERE groupid=? AND userid=?;`, groupID, sess.UserID).Scan(&tmp) == nil
	isSuperAdmin := false
	db.QueryRow(`SELECT is_superadmin FROM users WHERE id=?`, sess.UserID).Scan(&isSuperAdmin)
	isOwner := db.QueryRow(`SELECT userid FROM comments WHERE id=?;`, commentID).Scan(&tmp) == nil

	if !isOwner && !isMod && !isAdmin && !isSuperAdmin {
		ErrForbiddenHandler(w, r)
		return
	}

	if r.Method == "POST" {
		action := r.PostFormValue("action")
		if action == "Update" {
			if content == "" {
				http.Redirect(w, r, "/comments/edit?id="+commentID, http.StatusSeeOther)
				return
			}
			if !isMod && !isAdmin && !isSuperAdmin {
				isSticky = (pos < 0)
			}
			if isSticky {
				if pos > 0 {
					pos = -pos
				}
			} else {
				if pos < 0 {
					pos = -pos
				}
			}
			db.Exec(`UPDATE comments SET content=?, pos=?, updated_date=? WHERE id=?;`, content, pos, int64(time.Now().Unix()), commentID)
			page := pos / numCommentsPerPage
			if page < 0 {
				page = 0
			}
			http.Redirect(w, r, "/topics?id="+topicID+"&p="+strconv.Itoa(page)+"#comment-"+commentID, http.StatusSeeOther)
		}
		if action == "Delete" {
			db.Exec(`UPDATE comments SET is_deleted=1 WHERE id=?;`, commentID)
			http.Redirect(w, r, "/comments/edit?id="+commentID, http.StatusSeeOther)
		}
		if action == "Undelete" {
			db.Exec(`UPDATE comments SET is_deleted=0 WHERE id=?;`, commentID)
			http.Redirect(w, r, "/comments/edit?id="+commentID, http.StatusSeeOther)
		}
		return
	}
	isDeleted := false
	db.QueryRow(`SELECT content, is_deleted FROM comments WHERE id=?;`, commentID).Scan(&content, &isDeleted)
	isSticky = (pos < 0)

	templates.Render(w, "commentedit.html", map[string]interface{}{
		"Common": readCommonData(r, sess),
		"TopicID": topicID,
		"TopicOwnerName": topicOwnerName,
		"TopicCreatedDate": timeAgoFromNow(time.Unix(topicCreatedDate, 0)),
		"CommentID": commentID,
		"TopicName": topicName,
		"GroupName": groupName,
		"ParentComment": parentComment,
		"Content": content,
		"IsSticky": isSticky,
		"IsMod": isMod,
		"IsAdmin": isAdmin,
		"IsSuperAdmin": isSuperAdmin,
		"IsDeleted": isDeleted,
		"IsImageUploadEnabled": false,
	})
})

