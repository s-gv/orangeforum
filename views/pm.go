// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package views

import (
	"github.com/s-gv/orangeforum/models/db"
	"github.com/s-gv/orangeforum/templates"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var messagesPerPage = 50

var PrivateMessageHandler = A(func(w http.ResponseWriter, r *http.Request, sess Session) {
	startDate := time.Now().Unix()
	lmd := r.FormValue("lmd")
	if lmd != "" {
		if d, err := strconv.Atoi(lmd); err == nil {
			startDate = int64(d)
		}
	}

	type Message struct {
		ID          string
		From        string
		To          string
		IsRead      bool
		CreatedDate string
		Content     template.HTML
	}

	var lastMessageDate int64
	var msgs []Message
	var cDate int64
	var content string
	var rows *db.Rows
	rows = db.Query(`SELECT messages.id, fromusers.username, tousers.username, messages.content, messages.is_read, messages.created_date
		FROM messages INNER JOIN users fromusers ON fromusers.id=messages.fromid INNER JOIN users tousers ON tousers.id=messages.toid
		WHERE messages.toid=? AND messages.created_date <= ? ORDER BY messages.created_date DESC LIMIT ?;`, sess.UserID, startDate, messagesPerPage+1)
	for rows.Next() {
		msg := Message{}
		rows.Scan(&msg.ID, &msg.From, &msg.To, &content, &msg.IsRead, &cDate)
		msg.CreatedDate = timeAgoFromNow(time.Unix(cDate, 0))
		msg.Content = formatComment(content)
		if len(msgs) < messagesPerPage {
			msgs = append(msgs, msg)
		} else {
			lastMessageDate = cDate
		}
	}

	to, cont := "", ""

	if pmid := r.FormValue("quote"); pmid != "" {
		db.QueryRow(`SELECT users.username, messages.content
			FROM messages INNER JOIN users ON messages.fromid=users.id
			WHERE messages.id=?;`, pmid).Scan(&to, &cont)
		cont = formatReply(to, cont)
	}

	if flag := r.FormValue("flag"); flag != "" {
		rows := db.Query(`SELECT users.username FROM mods
			INNER JOIN users ON users.id=mods.userid
			INNER JOIN topics ON topics.groupid=mods.groupid
			INNER JOIN comments ON comments.topicid=topics.id
			WHERE comments.id=?;`, flag)
		for rows.Next() {
			var mod string
			rows.Scan(&mod)
			if to != "" {
				to = to + ", "
			}
			to = to + mod
		}
		cont = "Flagging " + "http://" + r.Host + "/comments?id=" + flag
	}

	if lmd != "" && len(msgs) == 0 {
		http.Redirect(w, r, "/pm", http.StatusSeeOther)
		return
	}

	db.Exec(`UPDATE messages SET is_read=? WHERE toid=?;`, true, sess.UserID)

	templates.Render(w, "pm.html", map[string]interface{}{
		"Common":           readCommonData(r, sess),
		"Messages":         msgs,
		"LastMessageDate":  lastMessageDate,
		"FirstMessageDate": startDate,
		"To":               to,
		"Content":          cont,
	})
})

var PrivateMessageCreateHandler = A(func(w http.ResponseWriter, r *http.Request, sess Session) {
	if r.Method == "POST" {
		tousers := strings.TrimSpace(r.PostFormValue("to"))
		content := strings.TrimSpace(r.PostFormValue("content"))

		if tousers == "" {
			sess.SetFlashMsg("No users to send the message to.")
			http.Redirect(w, r, "/pm", http.StatusSeeOther)
			return
		}
		if content == "" {
			sess.SetFlashMsg("Content is empty.")
			http.Redirect(w, r, "/pm", http.StatusSeeOther)
			return
		}

		tousernames := strings.Split(tousers, ",")
		touserids := []string{}
		for _, tousername := range tousernames {
			username := strings.TrimSpace(tousername)
			var userid string
			if err := db.QueryRow(`SELECT id FROM users WHERE username=?;`, username).Scan(&userid); err == nil {
				touserids = append(touserids, userid)
			} else {
				sess.SetFlashMsg("Username not found: " + username)
				http.Redirect(w, r, "/pm#end", http.StatusSeeOther)
				return
			}
		}

		for _, userid := range touserids {
			db.Exec(`INSERT INTO messages(fromid, toid, content, created_date) VALUES(?, ?, ?, ?);`, sess.UserID, userid, content, int(time.Now().Unix()))
		}

		sess.SetFlashMsg("Message sent.")
		http.Redirect(w, r, "/pm#end", http.StatusSeeOther)
		return
	}
})

var PrivateMessageDeleteHandler = A(func(w http.ResponseWriter, r *http.Request, sess Session) {
	if r.Method == "POST" {
		id := r.PostFormValue("id")
		db.Exec(`DELETE FROM messages WHERE id=? AND toid=?;`, id, sess.UserID)
		http.Redirect(w, r, "/pm?lmd="+r.PostFormValue("lmd"), http.StatusSeeOther)
		return
	}
})
