// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package views

import (
	"net/http"
	"github.com/s-gv/orangeforum/templates"
	"log"
	"github.com/s-gv/orangeforum/models"
	"github.com/s-gv/orangeforum/utils"
	"strings"
	"errors"
	"html/template"
	"strconv"
	"github.com/s-gv/orangeforum/models/db"
	"sort"
	"time"
	"net/url"
	"database/sql"
	"os"
	"io"
	"path/filepath"
	"regexp"
	"runtime/debug"
)

var linkRe *regexp.Regexp

func ErrServerHandler(w http.ResponseWriter, r *http.Request) {
	if r := recover(); r != nil {
		log.Printf("[INFO] Recovered from panic: %s\n[INFO] Debug stack: %s\n", r, debug.Stack())
		http.Error(w, "Internal server error. This event has been logged.", http.StatusInternalServerError)
	}
}

func ErrNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func ErrForbiddenHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "403 Forbidden", http.StatusForbidden)
}

func UA(handler func(w http.ResponseWriter, r *http.Request, sess models.Session)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer ErrServerHandler(w, r)
		sess := models.OpenSession(w, r)
		if r.Method == "POST" && r.PostFormValue("csrf") != sess.CSRFToken {
			ErrForbiddenHandler(w, r)
			return
		}
		//log.Printf("[INFO] Request: %s\n", r.URL)
		handler(w, r, sess)
	}
}

func A(handler func(w http.ResponseWriter, r *http.Request, sess models.Session)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer ErrServerHandler(w, r)
		sess := models.OpenSession(w, r)
		if r.Method == "POST" && r.PostFormValue("csrf") != sess.CSRFToken {
			ErrForbiddenHandler(w, r)
			return
		}
		if !sess.UserID.Valid {
			redirectURL := r.URL.Path
			if r.URL.RawQuery != "" {
				redirectURL += "?"+r.URL.RawQuery
			}
			http.Redirect(w, r, "/login?next="+url.QueryEscape(redirectURL), http.StatusSeeOther)
			return
		}
		//log.Printf("[INFO] Request: %s\n", r.URL)
		handler(w, r, sess)
	}
}

var IndexHandler = UA(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
	if r.URL.Path != "/" {
		ErrNotFoundHandler(w, r)
		return
	}

	type Group struct {
		Name string
		Desc string
		IsSticky int
	}
	groups := []Group{}
	rows := db.Query(`SELECT name, description, is_sticky FROM groups WHERE is_closed=0 ORDER BY is_sticky DESC, RANDOM() LIMIT 25;`)
	for rows.Next() {
		groups = append(groups, Group{})
		g := &groups[len(groups)-1]
		rows.Scan(&g.Name, &g.Desc, &g.IsSticky)
	}
	sort.Slice(groups, func(i, j int) bool {return groups[i].Name < groups[j].Name})
	sort.Slice(groups, func(i, j int) bool {return groups[i].IsSticky > groups[j].IsSticky})

	type Topic struct {
		ID string
		Title string
		GroupName string
		OwnerName string
		CreatedDate string
		NumComments int
	}
	topics := []Topic{}
	trows := db.Query(`SELECT topics.id, topics.title, topics.num_comments, topics.created_date, groups.name, users.username FROM topics INNER JOIN groups ON topics.groupid=groups.id INNER JOIN users ON topics.userid=users.id ORDER BY topics.created_date DESC LIMIT 20;`)
	for trows.Next() {
		topics = append(topics, Topic{})
		t := &topics[len(topics)-1]
		var cDate int64
		trows.Scan(&t.ID, &t.Title, &t.NumComments, &cDate, &t.GroupName, &t.OwnerName)
		t.CreatedDate = timeAgoFromNow(time.Unix(cDate, 0))
	}
	templates.Render(w, "index.html", map[string]interface{}{
		"Common": models.ReadCommonData(r, sess),
		"GroupCreationDisabled": models.Config(models.GroupCreationDisabled) == "1",
		"HeaderMsg": models.Config(models.HeaderMsg),
		"Groups": groups,
		"Topics": topics,
	})
})

func init() {
	linkRe = regexp.MustCompile("https?://[^\\s]+[A-Za-z0-9/\\&\\+\\?#,_-]")
}

func timeAgoFromNow(t time.Time) string {
	diff := time.Now().Sub(t)
	if diff.Hours() > 24 {
		return strconv.Itoa(int(diff.Hours()/24)) + " days ago"
	} else if diff.Hours() >= 2 {
		return strconv.Itoa(int(diff.Hours())) + " hours ago"
	} else {
		return strconv.Itoa(int(diff.Minutes())) + " minutes ago"
	}
	return ""
}

func validateName(name string) error {
	if len(name) == 0 {
		return errors.New("Name cannot be blank.")
	}
	hasSpecial := false
	for _, ch := range name {
		if (ch < 'A' || ch > 'Z') && (ch < 'a' || ch > 'z') && ch != '_' && ch != '-' && (ch < '0' || ch > '9') {
			hasSpecial = true
		}
	}
	if hasSpecial {
		return errors.New("Name can contain only english alphabets, numbers, hyphens, and underscore.")
	}
	return nil
}

func formatComment(comment string) template.HTML {
	comment = strings.Replace(comment, "\r", "", -1)
	formatted := "<p>"
	preClosed := true
	for _, para := range strings.Split(comment, "\n") {
		if para == "" {
			if !preClosed {
				formatted = formatted + "</pre>"
			}
			formatted = formatted + "</p><p>"
		} else {
			if len(para) > 4 && para[:4] == "    " {
				if preClosed {
					formatted = formatted + "<pre>"
					preClosed = false
				}
				formatted = formatted + template.HTMLEscapeString(para[4:]) + "\n"
			} else {
				escapedPara := template.HTMLEscapeString(para)
				linkedPara := linkRe.ReplaceAllString(escapedPara, "<a href=\"$0\">$0</a>")
				formatted = formatted + linkedPara + " "
			}
		}
	}
	if !preClosed {
		formatted = formatted + "</pre>"
	}
	formatted = formatted + "</p>"

	return template.HTML(formatted)
}

func saveImage(r *http.Request) string {
	imageName := ""
	if dataDir := models.Config(models.DataDir); dataDir != "" {
		r.ParseMultipartForm(32*1024*1024)
		file, handler, err := r.FormFile("img")
		if err == nil {
			defer file.Close()
			if handler.Filename != "" {
				ext := strings.ToLower(filepath.Ext(handler.Filename))
				if ext == ".jpg" || ext == ".png" || ext == ".jpeg" {
					fileName := models.RandSeq(64) + ext
					f, err := os.OpenFile(dataDir+fileName, os.O_WRONLY|os.O_CREATE, 0666)
					if err == nil {
						defer f.Close()
						io.Copy(f, file)
						imageName = fileName
					} else {
						log.Panicf("[ERROR] Error writing opening file: %s\n", err)
					}
				}
			}
		}
	} else {
		log.Panicf("[ERROR] Unable to accept file upload. DataDir not configured.\n")
	}
	return imageName
}

var GroupEditHandler = A(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
	if models.Config(models.GroupCreationDisabled) == "1" {
		ErrForbiddenHandler(w, r)
		return
	}
	commonData := models.ReadCommonData(r, sess)

	userName := commonData.UserName

	groupID := r.FormValue("id")
	name := r.FormValue("name")
	desc := r.FormValue("desc")
	headerMsg := r.FormValue("header_msg")
	isSticky := r.FormValue("is_sticky") != ""
	isDeleted := false
	mods := strings.Split(r.FormValue("mods"), ",")
	for i, mod := range mods {
		mods[i] = strings.TrimSpace(mod)
	}
	admins := strings.Split(r.FormValue("admins"), ",")
	for i, admin := range admins {
		admins[i] = strings.TrimSpace(admin)
	}
	if len(admins) == 1 && admins[0] == "" {
		admins[0] = userName
	}
	action := r.FormValue("action")

	if groupID != "" {
		if !models.IsUserGroupAdmin(strconv.Itoa(int(sess.UserID.Int64)), groupID) && !commonData.IsSuperAdmin {
			ErrForbiddenHandler(w, r)
			return
		}
	}

	if r.Method == "POST" {
		if action == "Create" {
			if err := validateName(name); err != nil {
				sess.SetFlashMsg(err.Error())
				http.Redirect(w, r, "/groups/edit", http.StatusSeeOther)
				return
			}
			db.Exec(`INSERT INTO groups(name, description, header_msg, is_sticky, created_date, updated_date) VALUES(?, ?, ?, ?, ?, ?);`, name, desc, headerMsg, isSticky, time.Now().Unix(), time.Now().Unix())
			groupID := models.ReadGroupIDByName(name)
			for _, mod := range mods {
				if mod != "" {
					models.CreateMod(mod, groupID)
				}
			}
			for _, admin := range admins {
				if admin != "" {
					models.CreateAdmin(admin, groupID)
				}
			}
			http.Redirect(w, r, "/groups?name="+name, http.StatusSeeOther)
		} else if action == "Update" {
			if err := validateName(name); err != nil {
				sess.SetFlashMsg(err.Error())
				http.Redirect(w, r, "/groups/edit?id="+groupID, http.StatusSeeOther)
				return
			}
			isUserSuperAdmin := false
			db.QueryRow(`SELECT is_superadmin FROM users WHERE id=?;`, sess.UserID).Scan(&isUserSuperAdmin)
			if !isUserSuperAdmin {
				db.QueryRow(`SELECT is_sticky FROM groups WHERE id=?;`, groupID).Scan(&isSticky)
			}
			db.Exec(`UPDATE groups SET name=?, description=?, header_msg=?, is_sticky=?, updated_date=? WHERE id=?;`, name, desc, headerMsg, isSticky, time.Now().Unix(), groupID)
			models.DeleteAdmins(groupID)
			models.DeleteMods(groupID)
			for _, mod := range mods {
				if mod != "" {
					models.CreateMod(mod, groupID)
				}
			}
			for _, admin := range admins {
				if admin != "" {
					models.CreateAdmin(admin, groupID)
				}
			}
			http.Redirect(w, r, "/groups?name="+name, http.StatusSeeOther)
		} else if action == "Delete" {
			models.DeleteGroup(groupID)
			http.Redirect(w, r, "/groups/edit?id="+groupID, http.StatusSeeOther)
		} else if action == "Undelete" {
			models.UndeleteGroup(groupID)
			http.Redirect(w, r, "/groups/edit?id="+groupID, http.StatusSeeOther)
		}
		return
 	}

	if groupID != "" {
		// Open to edit
		db.QueryRow(`SELECT name, description, header_msg, is_sticky, is_closed FROM groups WHERE id=?;`, groupID).Scan(
			&name, &desc, &headerMsg, &isSticky, &isDeleted,
		)
		mods = models.ReadMods(groupID)
		admins = models.ReadAdmins(groupID)
	}

	templates.Render(w, "groupedit.html", map[string]interface{}{
		"Common": models.ReadCommonData(r, sess),
		"ID": groupID,
		"GroupName": name,
		"Desc": desc,
		"HeaderMsg": headerMsg,
		"IsSticky": isSticky,
		"IsDeleted": isDeleted,
		"Mods": strings.Join(mods, ", "),
		"Admins": strings.Join(admins, ", "),
	})
})

var GroupHandler = UA(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
	name := r.FormValue("name")
	var groupID, groupDesc, headerMsg string
	if db.QueryRow(`SELECT id, description, header_msg FROM groups WHERE name=?;`, name).Scan(&groupID, &groupDesc, &headerMsg) != nil {
		ErrNotFoundHandler(w, r)
		return
	}

	subToken := ""
	if sess.UserID.Valid {
		db.QueryRow(`SELECT token FROM groupsubscriptions WHERE groupid=? AND userid=?;`, groupID, sess.UserID).Scan(&subToken)
	}


	numTopicsPerPage := 30
	lastTopicDate, err := strconv.ParseInt(r.FormValue("ltd"), 10, 64)
	if err != nil {
		lastTopicDate = 0
	}

	type Topic struct {
		ID int
		Title string
		IsDeleted bool
		IsClosed bool
		Owner string
		NumComments int
		CreatedDate string
		cDateUnix int64
	}
	var topics []Topic
	var rows *db.Rows
	if lastTopicDate == 0 {
		rows = db.Query(`SELECT topics.id, topics.title, topics.is_deleted, topics.is_closed, topics.num_comments, topics.created_date, users.username FROM topics INNER JOIN users ON topics.userid = users.id AND topics.groupid=? ORDER BY topics.is_sticky DESC, topics.activity_date DESC LIMIT ?;`, groupID, numTopicsPerPage)
	} else {
		rows = db.Query(`SELECT topics.id, topics.title, topics.is_deleted, topics.is_closed, topics.num_comments, topics.created_date, users.username FROM topics INNER JOIN users ON topics.userid = users.id AND topics.groupid=? AND topics.is_sticky=0 AND topics.created_date < ? ORDER BY topics.activity_date DESC LIMIT ?;`, groupID, lastTopicDate, numTopicsPerPage)
	}
	for rows.Next() {
		t := Topic{}
		rows.Scan(&t.ID, &t.Title, &t.IsDeleted, &t.IsClosed, &t.NumComments, &t.cDateUnix, &t.Owner)
		t.CreatedDate = timeAgoFromNow(time.Unix(t.cDateUnix, 0))
		topics = append(topics, t)
	}

	isSuperAdmin := false
	isAdmin := false
	isMod := false
	if sess.IsUserValid() {
		db.QueryRow(`SELECT is_superadmin FROM users WHERE id=?;`, sess.UserID).Scan(&isSuperAdmin)
		var tmp string
		isAdmin = db.QueryRow(`SELECT id FROM admins WHERE groupid=? AND userid=?;`, groupID, sess.UserID).Scan(&tmp) == nil
		isMod = db.QueryRow(`SELECT id FROM mods WHERE groupid=? AND userid=?;`, groupID, sess.UserID).Scan(&tmp) == nil
	}

	if len(topics) >= numTopicsPerPage {
		lastTopicDate = topics[len(topics)-1].cDateUnix
	} else {
		lastTopicDate = 0
	}

	templates.Render(w, "groupindex.html", map[string]interface{}{
		"Common": models.ReadCommonData(r, sess),
		"GroupName": name,
		"GroupDesc": groupDesc,
		"GroupID": groupID,
		"HeaderMsg": headerMsg,
		"SubToken": subToken,
		"Topics": topics,
		"IsMod": isMod,
		"IsAdmin": isAdmin,
		"IsSuperAdmin": isSuperAdmin,
		"LastTopicDate": lastTopicDate,
	})
})

var TopicCreateHandler = A(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
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
		if len(title) < 1 || len(title) > 150 {
			sess.SetFlashMsg("Invalid number of characters in the title. Valid range: 1-150.")
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
		"Common":    models.ReadCommonData(r, sess),
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

var TopicUpdateHandler = A(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
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
		if len(title) < 1 || len(title) > 150 {
			sess.SetFlashMsg("Invalid number of characters in the title. Valid range: 1-150.")
			http.Redirect(w, r, "/topics/edit?id="+topicID, http.StatusSeeOther)
			return
		}

		if action == "Update" {
			db.Exec(`UPDATE topics SET title=?, content=?, is_sticky=?, updated_date=? WHERE id=?;`, title, content, isSticky, int(time.Now().Unix()), topicID)
		} else if action == "Close" {
			db.Exec(`UPDATE topics SET is_closed=1 WHERE id=?;`, topicID)
		} else if action == "Reopen" {
			db.Exec(`UPDATE topics SET is_closed=0 WHERE id=?;`, topicID)
		} else if action == "Delete" {
			db.Exec(`UPDATE topics SET is_deleted=1 WHERE id=?;`, topicID)
		} else if action == "Undelete" {
			db.Exec(`UPDATE topics SET is_deleted=0 WHERE id=?;`, topicID)
		}
		http.Redirect(w, r, "/topics/edit?id="+topicID, http.StatusSeeOther)
		return
	}

	if db.QueryRow(`SELECT title, content, is_sticky, is_deleted, is_closed FROM topics WHERE id=?;`, topicID).Scan(&title, &content, &isSticky, &isDeleted, &isClosed) != nil {
		ErrNotFoundHandler(w, r)
		return
	}

	templates.Render(w, "topicedit.html", map[string]interface{}{
		"Common": models.ReadCommonData(r, sess),
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

var TopicHandler = UA(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
	topicID := r.FormValue("id")
	lastCommentDate, err := strconv.ParseInt(r.FormValue("lcd"), 10, 64)
	if err != nil {
		lastCommentDate = 0
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

	type Comment struct {
		ID string
		Content template.HTML
		ImgSrc string
		CreatedDate string
		UserName string
		IsOwner bool
		IsDeleted bool
	}
	numCommentsPerPage := 100

	var comments []Comment
	var cDate int64
	var rows *db.Rows
	if lastCommentDate == 0 {
		rows = db.Query(`SELECT users.id, users.username, comments.id, comments.content, comments.image, comments.is_deleted, comments.created_date FROM comments INNER JOIN users ON comments.userid=users.id AND comments.topicid=? ORDER BY comments.is_sticky DESC, comments.created_date ASC LIMIT ?;`, topicID, numCommentsPerPage)
	} else {
		rows = db.Query(`SELECT users.id, users.username, comments.id, comments.content, comments.image, comments.is_deleted, comments.created_date FROM comments INNER JOIN users ON comments.userid=users.id AND comments.topicid=? AND comments.created_date > ? AND comments.is_sticky=0 ORDER BY comments.created_date ASC LIMIT ?;`, topicID, lastCommentDate, numCommentsPerPage)
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
	if len(comments) >= numCommentsPerPage {
		lastCommentDate = cDate
	} else {
		lastCommentDate = 0
	}

	var tmp string
	db.QueryRow(`SELECT name FROM groups WHERE id=?;`, groupID).Scan(&groupName)
	isMod := db.QueryRow(`SELECT id FROM mods WHERE groupid=? AND userid=?;`, groupID, sess.UserID).Scan(&tmp) == nil
	isAdmin := db.QueryRow(`SELECT id FROM admins WHERE groupid=? AND userid=?;`, groupID, sess.UserID).Scan(&tmp) == nil
	isSuperAdmin := false
	db.QueryRow(`SELECT is_superadmin FROM users WHERE id=?`, sess.UserID).Scan(&isSuperAdmin)
	isOwner := sess.UserID.Valid && ownerID == sess.UserID.Int64

	templates.Render(w, "topicindex.html", map[string]interface{}{
		"Common": models.ReadCommonData(r, sess),
		"GroupID": groupID,
		"TopicID": topicID,
		"GroupName": groupName,
		"TopicName": title,
		"OwnerName": ownerName,
		"CreatedDate": timeAgoFromNow(time.Unix(createdDate, 0)),
		"SubToken": subToken,
		"Title": title,
		"Content": content,
		"IsClosed": isClosed,
		"IsOwner": isOwner,
		"IsMod": isMod,
		"IsAdmin": isAdmin,
		"IsSuperAdmin": isSuperAdmin,
		"Comments": comments,
		"LastCommentDate": lastCommentDate,
	})
})

var CommentCreateHandler = A(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
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
			http.Redirect(w, r, "/comments/new?tid="+topicID, http.StatusSeeOther)
			return
		}
		if !isMod && !isAdmin && !isSuperAdmin {
			isSticky = false
		}
		imageName := ""
		if isImageUploadEnabled {
			imageName = saveImage(r)
		}
		db.Exec(`INSERT INTO comments(content, image, topicid, userid, parentid, is_sticky, created_date, updated_date) VALUES(?, ?, ?, ?, ?, ?, ?, ?);`,
			content, imageName, topicID, sess.UserID, sql.NullInt64{Valid:false}, isSticky, int64(time.Now().Unix()), int64(time.Now().Unix()))
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
		http.Redirect(w, r, "/topics?id="+topicID+"#comment-last", http.StatusSeeOther)
		return
	}

	templates.Render(w, "commentedit.html", map[string]interface{}{
		"Common": models.ReadCommonData(r, sess),
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

var CommentUpdateHandler = A(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
	commentID := r.FormValue("id")
	content := r.PostFormValue("content")
	isSticky := r.PostFormValue("is_sticky") != ""

	var groupID, topicID, groupName, topicName, parentComment, topicOwnerName, topicOwnerID string
	var topicCreatedDate int64
	if db.QueryRow(`SELECT topicid FROM comments WHERE id=?;`, commentID).Scan(&topicID) != nil {
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
				db.QueryRow(`SELECT is_sticky FROM comments WHERE id=?;`, commentID).Scan(&isSticky)
			}
			db.Exec(`UPDATE comments SET content=?, is_sticky=?, updated_date=? WHERE id=?;`, content, isSticky, int64(time.Now().Unix()), commentID)
			http.Redirect(w, r, "/topics?id="+topicID+"#comment-"+commentID, http.StatusSeeOther)
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
	db.QueryRow(`SELECT content, is_sticky, is_deleted FROM comments WHERE id=?;`, commentID).Scan(&content, &isSticky, &isDeleted)

	templates.Render(w, "commentedit.html", map[string]interface{}{
		"Common": models.ReadCommonData(r, sess),
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

var CommentHandler = UA(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
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
		"Common": models.ReadCommonData(r, sess),
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

var TopicSubscribeHandler = A(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
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
				sess.UserID, topicID, models.RandSeq(64), time.Now().Unix())
		}
	}
	http.Redirect(w, r, "/topics?id="+topicID, http.StatusSeeOther)
})

var TopicUnsubscribeHandler = UA(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
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

var GroupSubscribeHandler = A(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
	groupID := r.FormValue("id")
	if models.Config(models.AllowGroupSubscription) == "0" {
		ErrForbiddenHandler(w, r)
		return
	}
	var groupName string
	if db.QueryRow(`SELECT name FROM groups WHERE id=?;`, groupID).Scan(&groupName) != nil {
		ErrNotFoundHandler(w, r)
		return
	}
	if r.Method == "POST" {
		var tmp string
		if db.QueryRow(`SELECT id FROM groupsubscriptions WHERE userid=? AND groupid=?;`, sess.UserID, groupID).Scan(&tmp) != nil {
			db.Exec(`INSERT INTO groupsubscriptions(userid, groupid, token, created_date) VALUES(?, ?, ?, ?);`,
				sess.UserID, groupID, models.RandSeq(64), time.Now().Unix())
		}
	}
	http.Redirect(w, r, "/groups?name="+groupName, http.StatusSeeOther)
})

var GroupUnsubscribeHandler = UA(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
	token := r.FormValue("token")
	var groupID, groupName string
	if db.QueryRow(`SELECT groupid FROM groupsubscriptions WHERE token=?;`, token).Scan(&groupID) != nil {
		ErrNotFoundHandler(w, r)
		return
	}
	db.QueryRow(`SELECT name FROM groups WHERE id=?;`, groupID).Scan(&groupName)
	if r.Method == "POST" {
		db.Exec(`DELETE FROM groupsubscriptions WHERE token=?;`, token)
		if r.PostFormValue("noredirect") != "" {
			w.Write([]byte("Unsubscribed."))
		} else {
			http.Redirect(w, r, "/groups?name="+groupName, http.StatusSeeOther)
		}
		return
	}
	w.Write([]byte(`<!DOCTYPE html><html><head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1"></head>
	<body><form action="/groups/unsubscribe" method="POST">
	Unsubscribe from `+groupName+`?
	<input type="hidden" name="token" value=`+token+`>
	<input type="hidden" name="csrf" value="`+sess.CSRFToken+`">
	<input type="hidden" name="noredirect" value="1">
	<input type="submit" value="Unsubscribe">
	</form></body></html>`))
})

var SignupHandler = UA(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
	redirectURL, err := url.QueryUnescape(r.FormValue("next"))
	if redirectURL == "" || err != nil {
		redirectURL = "/"
	}
	if sess.IsUserValid() {
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		userName := r.PostFormValue("username")
		passwd := r.PostFormValue("passwd")
		passwdConfirm := r.PostFormValue("confirm")
		email := r.PostFormValue("email")
		if len(userName) == 0 {
			sess.SetFlashMsg("Username cannot be blank.")
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}
		hasSpecial := false
		for _, ch := range userName {
			if (ch < 'A' || ch > 'Z') && (ch < 'a' || ch > 'z') && ch != '_' && (ch < '0' || ch > '9') {
				hasSpecial = true
			}
		}
		if hasSpecial {
			sess.SetFlashMsg("Username can contain only alphabets, numbers, and underscore.")
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}
		if models.ProbeUser(userName) {
			sess.SetFlashMsg("Username already registered.")
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}
		if err := validatePasswd(passwd, passwdConfirm); err != nil {
			sess.SetFlashMsg(err.Error())
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}
		models.CreateUser(userName, passwd, email)
		sess.Authenticate(userName, passwd)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
	templates.Render(w, "signup.html", map[string]interface{}{
		"Common": models.ReadCommonData(r, sess),
		"next": template.URL(url.QueryEscape(redirectURL)),
	})
})

var LoginHandler = UA(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
	redirectURL, err := url.QueryUnescape(r.FormValue("next"))
	if redirectURL == "" || err != nil {
		redirectURL = "/"
	}
	if sess.IsUserValid() {
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		userName := r.PostFormValue("username")
		passwd := r.PostFormValue("passwd")
		if sess.Authenticate(userName, passwd) {
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		} else {
			sess.SetFlashMsg("Incorrect username/password")
			http.Redirect(w, r, "/login?next="+redirectURL, http.StatusSeeOther)
			return
		}
	}
	templates.Render(w, "login.html", map[string]interface{}{
		"Common": models.ReadCommonData(r, sess),
		"next": template.URL(url.QueryEscape(redirectURL)),
	})
})

var ChangePasswdHandler = UA(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
	userName, err := sess.UserName()
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method == "POST" {
		passwd := r.PostFormValue("passwd")
		newPasswd := r.PostFormValue("newpass")
		newPasswdConfirm := r.PostFormValue("confirm")
		if !sess.Authenticate(userName, passwd) {
			sess.SetFlashMsg("Current password incorrect.")
			http.Redirect(w, r, "/changepass", http.StatusSeeOther)
			return
		}
		if err := validatePasswd(newPasswd, newPasswdConfirm); err != nil {
			sess.SetFlashMsg(err.Error())
			http.Redirect(w, r, "/changepass", http.StatusSeeOther)
			return
		}
		if err := models.UpdateUserPasswd(userName, newPasswd); err != nil {
			log.Panicf("[ERROR] Error changing password: %s\n", err)
		}
		sess.SetFlashMsg("Password change successful.")
		http.Redirect(w, r, "/changepass", http.StatusSeeOther)
		return
	}
	templates.Render(w, "changepass.html", map[string]interface{}{
		"Common": models.ReadCommonData(r, sess),
	})
})

var ForgotPasswdHandler = UA(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
	if r.Method == "POST" {
		userName := r.PostFormValue("username")
		if userName == "" || !models.ProbeUser(userName) {
			sess.SetFlashMsg("Username doesn't exist.")
			http.Redirect(w, r, "/forgotpass", http.StatusSeeOther)
			return
		}
		email := models.ReadUserEmail(userName)
		if !strings.ContainsRune(email, '@') {
			sess.SetFlashMsg("E-mail address not set. Contact site admin to reset the password.")
			http.Redirect(w, r, "/forgotpass", http.StatusSeeOther)
			return
		}
		forumName := models.Config(models.ForumName)
		resetLink := "https://" + r.Host + "/resetpass?r=" + models.CreateResetToken(userName)
		sub := forumName + " Password Recovery"
		msg := "Someone (hopefully you) requested we reset your password at " + forumName + ".\r\n" +
			"If you want to change it, visit "+resetLink+"\r\n\r\nIf not, just ignore this message."
		utils.SendMail(email, sub, msg)
		sess.SetFlashMsg("Password reset link sent to your e-mail.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return

	}
	templates.Render(w, "forgotpass.html", map[string]interface{}{
		"Common": models.ReadCommonData(r, sess),
	})
})

func validatePasswd(passwd string, passwdConfirm string) error {
	if len(passwd) < 8 {
		return errors.New("Password should have at least 8 characters.")
	}
	if passwd != passwdConfirm {
		return errors.New("Passwords don't match.")
	}
	return nil
}

var ResetPasswdHandler = UA(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
	resetToken := r.FormValue("r")
	userName, err := models.ReadUserNameByToken(resetToken)
	if err != nil {
		ErrForbiddenHandler(w, r)
		return
	}
	if r.Method == "POST" {
		passwd := r.PostFormValue("passwd")
		passwdConfirm := r.PostFormValue("confirm")
		if err := validatePasswd(passwd, passwdConfirm); err != nil {
			sess.SetFlashMsg(err.Error())
			http.Redirect(w, r, "/resetpass?r="+resetToken, http.StatusSeeOther)
			return
		}
		models.UpdateUserPasswd(userName, passwd)
		sess.SetFlashMsg("Password change successful.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	templates.Render(w, "resetpass.html", map[string]interface{}{
		"ResetToken": resetToken,
		"Common": models.ReadCommonData(r, sess),
	})
})

func TestHandler(w http.ResponseWriter, r *http.Request) {
	defer ErrServerHandler(w, r)
	sess := models.OpenSession(w, r)
	sess.SetFlashMsg("hi there")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	defer ErrServerHandler(w, r)
	models.ClearSession(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

var AdminIndexHandler = A(func (w http.ResponseWriter, r *http.Request, sess models.Session) {
	if !sess.IsUserSuperAdmin() {
		ErrForbiddenHandler(w, r)
		return
	}

	linkID := r.PostFormValue("linkid")

	if r.Method == "POST" && linkID == "" {
		forumName := r.PostFormValue("forum_name")
		headerMsg := r.PostFormValue("header_msg")
		signupDisabled := "0"
		groupCreationDisabled := "0"
		imageUploadEnabled := "0"
		allowGroupSubscription := "0"
		allowTopicSubscription := "0"
		dataDir := r.PostFormValue("data_dir")
		bodyAppendage := r.PostFormValue("body_appendage")
		defaultFromEmail := r.PostFormValue("default_from_mail")
		smtpHost := r.PostFormValue("smtp_host")
		smtpPort := r.PostFormValue("smtp_port")
		smtpUser := r.PostFormValue("smtp_user")
		smtpPass := r.PostFormValue("smtp_pass")
		if r.PostFormValue("signup_disabled") != "" {
			signupDisabled = "1"
		}
		if r.PostFormValue("group_creation_disabled") != "" {
			groupCreationDisabled = "1"
		}
		if r.PostFormValue("image_upload_enabled") != "" {
			imageUploadEnabled = "1"
		}
		if r.PostFormValue("allow_group_subscription") != "" {
			allowGroupSubscription = "1"
		}
		if r.PostFormValue("allow_topic_subscription") != "" {
			allowTopicSubscription = "1"
		}
		if dataDir != "" {
			if dataDir[len(dataDir)-1] != '/' {
				dataDir = dataDir + "/"
			}
		}

		errMsg := ""
		if forumName == "" {
			errMsg = "Forum name is empty."
		}

		if errMsg == "" {
			models.WriteConfig(models.ForumName, forumName)
			models.WriteConfig(models.HeaderMsg, headerMsg)
			models.WriteConfig(models.SignupDisabled, signupDisabled)
			models.WriteConfig(models.GroupCreationDisabled, groupCreationDisabled)
			models.WriteConfig(models.ImageUploadEnabled, imageUploadEnabled)
			models.WriteConfig(models.AllowGroupSubscription, allowGroupSubscription)
			models.WriteConfig(models.AllowTopicSubscription, allowTopicSubscription)
			models.WriteConfig(models.DataDir, dataDir)
			models.WriteConfig(models.BodyAppendage, bodyAppendage)
			models.WriteConfig(models.DefaultFromMail, defaultFromEmail)
			models.WriteConfig(models.SMTPHost, smtpHost)
			models.WriteConfig(models.SMTPPort, smtpPort)
			models.WriteConfig(models.SMTPUser, smtpUser)
			models.WriteConfig(models.SMTPPass, smtpPass)
		} else {
			sess.SetFlashMsg(errMsg)
		}
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" && linkID != "" {
		name := r.PostFormValue("name")
		URL := r.PostFormValue("url")
		content := r.PostFormValue("content")
		if linkID == "new" {
			if name != "" && (URL != "" || content != "") {
				models.CreateExtraNote(name, URL, content)
			} else {
				sess.SetFlashMsg("Enter an external URL or type some content for the footer link.")
			}
		} else {
			if r.PostFormValue("submit") == "Delete" {
				models.DeleteExtraNote(linkID)
			} else {
				models.UpdateExtraNote(linkID, name, URL, content)
			}

		}
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	templates.Render(w, "adminindex.html", map[string]interface{}{
		"Common": models.ReadCommonData(r, sess),
		"Config": models.ConfigAllVals(),
		"ExtraNotes": models.ReadExtraNotes(),
		"NumUsers": models.NumUsers(),
		"NumGroups": models.NumGroups(),
		"NumTopics": models.NumTopics(),
		"NumComments": models.NumComments(),
	})
})

var UserProfileHandler = UA(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
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
		http.Redirect(w, r, "/users?u="+userName, http.StatusSeeOther)
		return
	}

	templates.Render(w, "profile.html", map[string]interface{}{
		"Common": models.ReadCommonData(r, sess),
		"UserName": userName,
		"About": about,
		"Email": email,
		"IsSelf": sess.UserID.Valid && (userID == sess.UserID.Int64),
		"IsBanned": isBanned,
	})
})

var UserCommentsHandler = UA(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
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
		Content string
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

		rows.Scan(&c.TopicName, &c.TopicID, &c.ID, &c.Content, &c.ImgSrc, &cDate, &c.IsDeleted)
		c.CreatedDate = timeAgoFromNow(time.Unix(cDate, 0))
	}

	if len(comments) >= commentsPerPage {
		lastCommentDate = cDate
	} else {
		lastCommentDate = 0
	}

	templates.Render(w, "profilecomments.html", map[string]interface{}{
		"Common": models.ReadCommonData(r, sess),
		"OwnerName": ownerName,
		"Comments": comments,
		"LastCommentDate": lastCommentDate,
	})
})

var UserTopicsHandler = UA(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
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
	}

	if len(topics) >= numTopicsPerPage {
		lastTopicDate = cDate
	} else {
		lastTopicDate = 0
	}

	templates.Render(w, "profiletopics.html", map[string]interface{}{
		"Common": models.ReadCommonData(r, sess),
		"OwnerName": ownerName,
		"Topics": topics,
		"LastTopicDate": lastTopicDate,
	})
})

var UserGroupsHandler = A(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
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
		"Common": models.ReadCommonData(r, sess),
		"OwnerName": ownerName,
		"AdminInGroups": adminInGroups,
		"ModInGroups": modInGroups,
	})
})

var NoteHandler = UA(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
	id := r.FormValue("id")
	if e, err := models.ReadExtraNote(id); err == nil {
		if e.URL == "" {
			templates.Render(w, "extranote.html", map[string]interface{}{
				"Common": models.ReadCommonData(r, sess),
				"Name": e.Name,
				"UpdatedDate": e.UpdatedDate,
				"Content": template.HTML(e.Content),
			})
			return
		} else {
			http.Redirect(w, r, e.URL, http.StatusSeeOther)
			return
		}
	}
	ErrNotFoundHandler(w, r)
})

func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	defer ErrServerHandler(w, r)
	dataDir := models.Config(models.DataDir)
	if dataDir != "" {
		http.ServeFile(w, r, dataDir+"favicon.ico")
		return
	}
	ErrNotFoundHandler(w, r)
}

func StyleHandler(w http.ResponseWriter, r *http.Request) {
	defer ErrServerHandler(w, r)
	w.Header().Set("Content-Type", "text/css")
	w.Header().Set("Cache-Control", "max-age=31536000, public")
	templates.Render(w, "style.css", map[string]interface{}{})
}

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	defer ErrServerHandler(w, r)
	dataDir := models.Config(models.DataDir)
	if dataDir != "" {
		http.ServeFile(w, r, dataDir+r.FormValue("name"))
		return
	}
	ErrNotFoundHandler(w, r)
}
