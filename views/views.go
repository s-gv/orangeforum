package views

import (
)
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
)

func ErrServerHandler(w http.ResponseWriter, r *http.Request) {
	if r := recover(); r != nil {
		log.Printf("[INFO] Recovered from panic: %s", r)
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
			http.Redirect(w, r, "/login?next="+redirectURL, http.StatusSeeOther)
			return
		}
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
	rows := db.Query(`SELECT name, desc, is_sticky FROM groups WHERE is_closed=0 ORDER BY is_sticky DESC, RANDOM() LIMIT 25;`)
	for rows.Next() {
		groups = append(groups, Group{})
		g := &groups[len(groups)-1]
		rows.Scan(&g.Name, &g.Desc, &g.IsSticky)
	}
	sort.Slice(groups, func(i, j int) bool {return groups[i].Name < groups[j].Name})
	sort.Slice(groups, func(i, j int) bool {return groups[i].IsSticky > groups[j].IsSticky})
	templates.Render(w, "index.html", map[string]interface{}{
		"Common": models.ReadCommonData(sess),
		"GroupCreationDisabled": models.Config(models.GroupCreationDisabled) == "1",
		"Groups": groups,
	})
})

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

var GroupHandler = UA(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
	name := r.FormValue("name")
	var groupID string
	if db.QueryRow(`SELECT id FROM groups WHERE name=?;`, name).Scan(&groupID) != nil {
		ErrNotFoundHandler(w, r)
		return
	}

	numTopicsPerPage := 30
	lastTopicDate, err := strconv.ParseInt(r.FormValue("ltd"), 10, 64)
	if err != nil {
		lastTopicDate = 0
	}

	type Topic struct {
		ID int
		Title string
		Owner string
		NumComments int
		CreatedDate string
		cDateUnix int64
	}
	var topics []Topic
	var rows *db.Rows
	if lastTopicDate == 0 {
		rows = db.Query(`SELECT topics.id, topics.title, topics.num_comments, topics.created_date, users.username FROM topics INNER JOIN users ON topics.userid = users.id AND topics.groupid=? ORDER BY topics.is_sticky DESC, topics.created_date DESC LIMIT ?;`, groupID, numTopicsPerPage)
	} else {
		rows = db.Query(`SELECT topics.id, topics.title, topics.num_comments, topics.created_date, users.username FROM topics INNER JOIN users ON topics.userid = users.id AND topics.groupid=? AND topics.is_sticky=0 AND topics.created_date < ? ORDER BY topics.created_date DESC LIMIT ?;`, groupID, lastTopicDate, numTopicsPerPage)
	}
	for rows.Next() {
		topics = append(topics, Topic{})
		t := &topics[len(topics)-1]
		rows.Scan(&t.ID, &t.Title, &t.NumComments, &t.cDateUnix, &t.Owner)
		diff := time.Now().Sub(time.Unix(t.cDateUnix, 0))
		if diff.Hours() > 24 {
			t.CreatedDate = strconv.Itoa(int(diff.Hours()/24)) + " ago"
		} else if diff.Hours() >= 2 {
			t.CreatedDate = strconv.Itoa(int(diff.Hours())) + " hours ago"
		} else {
			t.CreatedDate = strconv.Itoa(int(diff.Minutes())) + " minutes ago"
		}
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
		"Common": models.ReadCommonData(sess),
		"GroupName": name,
		"GroupID": groupID,
		"Topics": topics,
		"IsMod": isMod,
		"IsAdmin": isAdmin,
		"IsSuperAdmin": isSuperAdmin,
		"LastTopicDate": lastTopicDate,
	})
})

var GroupEditHandler = A(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
	if models.Config(models.GroupCreationDisabled) == "1" {
		ErrForbiddenHandler(w, r)
		return
	}
	if !sess.IsUserValid() {
		http.Redirect(w, r, "/login?next="+r.URL.Path, http.StatusSeeOther)
		return
	}
	commonData := models.ReadCommonData(sess)

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
	isUserInAdminList := commonData.IsSuperAdmin
	for _, u := range admins {
		if u == userName {
			isUserInAdminList = true
			break
		}
	}
	if !isUserInAdminList {
		if admins[0] == "" {
			admins[0] = userName
		} else {
			admins = append(admins, userName)
		}
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
			db.Exec(`INSERT INTO groups(name, desc, header_msg, is_sticky) VALUES(?, ?, ?, ?);`, name, desc, headerMsg, isSticky)
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
			db.Exec(`UPDATE groups SET name=?, desc=?, header_msg=?, is_sticky=? WHERE id=?;`, name, desc, headerMsg, isSticky, groupID)
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
		db.QueryRow(`SELECT name, desc, header_msg, is_sticky, is_closed FROM groups WHERE id=?;`, groupID).Scan(
			&name, &desc, &headerMsg, &isSticky, &isDeleted,
		)
		mods = models.ReadMods(groupID)
		admins = models.ReadAdmins(groupID)
	}

	templates.Render(w, "groupnew.html", map[string]interface{}{
		"Common": models.ReadCommonData(sess),
		"ID": groupID,
		"Name": name,
		"Desc": desc,
		"HeaderMsg": headerMsg,
		"IsSticky": isSticky,
		"IsDeleted": isDeleted,
		"Mods": strings.Join(mods, ", "),
		"Admins": strings.Join(admins, ", "),
	})
})

var TopicCreateHandler = A(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
	groupID := r.FormValue("gid")
	isGroupClosed := 1
	db.QueryRow(`SELECT is_closed FROM groups WHERE id=?;`, groupID).Scan(&isGroupClosed)
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
		db.Exec(`INSERT INTO topics(title, content, userid, groupid, is_sticky, created_date, updated_date) VALUES(?, ?, ?, ?, ?, ?, ?);`,
			title, content, sess.UserID, groupID, isSticky, int(time.Now().Unix()), int(time.Now().Unix()))
		var groupName string
		db.QueryRow(`SELECT name FROM groups WHERE id=?`, groupID).Scan(&groupName)
		http.Redirect(w, r, "/groups?name="+groupName, http.StatusSeeOther)
		return
	}

	templates.Render(w, "topicedit.html", map[string]interface{}{
		"Common":    models.ReadCommonData(sess),
		"GroupID":   groupID,
		"TopicID":   "",
		"Title":     "",
		"Content":   "",
		"IsSticky": false,
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
	isDeleted := true

	if db.QueryRow(`SELECT groupid FROM topics WHERE id=?;`, topicID).Scan(&groupID) != nil {
		ErrNotFoundHandler(w, r)
		return
	}

	isGroupClosed := 1
	db.QueryRow(`SELECT is_closed FROM groups WHERE id=?;`, groupID).Scan(&isGroupClosed)
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
		} else if action == "Delete" {
			db.Exec(`UPDATE topics SET is_closed=1 WHERE id=?;`, topicID)
		} else if action == "Undelete" {
			db.Exec(`UPDATE topics SET is_closed=0 WHERE id=?;`, topicID)
		}
		http.Redirect(w, r, "/topics/edit?id="+topicID, http.StatusSeeOther)
		return
	}

	if db.QueryRow(`SELECT title, content, is_sticky, is_closed FROM topics WHERE id=?;`, topicID).Scan(&title, &content, &isSticky, &isDeleted) != nil {
		ErrNotFoundHandler(w, r)
		return
	}

	templates.Render(w, "topicedit.html", map[string]interface{}{
		"Common": models.ReadCommonData(sess),
		"GroupID": groupID,
		"TopicID": topicID,
		"Title": title,
		"Content": content,
		"IsSticky": isSticky,
		"IsDeleted": isDeleted,
		"IsMod": isMod,
		"IsAdmin": isAdmin,
		"IsSuperAdmin": isSuperAdmin,
	})
})

var SignupHandler = UA(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
	redirectURL := r.FormValue("next")
	if redirectURL == "" {
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
		"Common": models.ReadCommonData(sess),
		"next": template.URL(redirectURL),
	})
})

var LoginHandler = UA(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
	redirectURL := r.FormValue("next")
	if redirectURL == "" {
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
		"Common": models.ReadCommonData(sess),
		"next": template.URL(redirectURL),
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
		"Common": models.ReadCommonData(sess),
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
		"Common": models.ReadCommonData(sess),
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
		"Common": models.ReadCommonData(sess),
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

var CreateGroupHandler = A(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
	if r.Method == "POST" {
		groupName := r.PostFormValue("name")
		if groupName == "" {
			sess.SetFlashMsg("Group name is empty.")
			http.Redirect(w, r, "/creategroup", http.StatusSeeOther)
			return
		}
		hasSpecial := false
		for _, ch := range groupName {
			if (ch < 'A' || ch > 'Z') && (ch < 'a' || ch > 'z') && ch != '-' && (ch < '0' || ch > '9') {
				hasSpecial = true
			}
		}
		if hasSpecial {
			sess.SetFlashMsg("Username can contain only english alphabets, numbers, and hyphen.")
			http.Redirect(w, r, "/creategroup", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/g/"+groupName, http.StatusSeeOther)
		return
	}

	templates.Render(w, "creategroup.html", map[string]interface{}{
		"Common": models.ReadCommonData(sess),
	})
})

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
		fileUploadEnabled := "0"
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
		if r.PostFormValue("file_upload_enabled") != "" {
			fileUploadEnabled = "1"
		}
		if r.PostFormValue("allow_group_subscription") != "" {
			allowGroupSubscription = "1"
		}
		if r.PostFormValue("allow_topic_subscription") != "" {
			allowTopicSubscription = "1"
		}
		if dataDir != "" {
			if dataDir[len(dataDir)-1] == '/' {
				dataDir = dataDir[:len(dataDir)-1]
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
			models.WriteConfig(models.FileUploadEnabled, fileUploadEnabled)
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
		"Common": models.ReadCommonData(sess),
		"Config": models.ConfigAllVals(),
		"ExtraNotes": models.ReadExtraNotes(),
		"NumUsers": models.NumUsers(),
		"NumGroups": models.NumGroups(),
		"NumTopics": models.NumTopics(),
		"NumComments": models.NumComments(),
	})
})

var UserProfileHandler = UA(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
	if r.Method == "POST" {
		userName, err := sess.UserName()
		if err == nil {
			email := r.FormValue("email")
			about := r.FormValue("about")
			models.UpdateUserProfile(userName, email, about)
		}
		http.Redirect(w, r, "/users?u="+userName, http.StatusSeeOther)
		return
	}

	userName := r.FormValue("u")

	if !models.ProbeUser(userName) {
		ErrNotFoundHandler(w, r)
		return
	}

	isSelf := false
	if u, err := sess.UserName(); err == nil {
		if u == userName {
			isSelf = true
		}
	}
	templates.Render(w, "profile.html", map[string]interface{}{
		"Common": models.ReadCommonData(sess),
		"UserName": userName,
		"About": models.ReadUserAbout(userName),
		"Email": models.ReadUserEmail(userName),
		"IsSelf": isSelf,
	})
})

var NoteHandler = UA(func(w http.ResponseWriter, r *http.Request, sess models.Session) {
	id := r.FormValue("id")
	if e, err := models.ReadExtraNote(id); err == nil {
		if e.URL == "" {
			templates.Render(w, "extranote.html", map[string]interface{}{
				"Common": models.ReadCommonData(sess),
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
		http.ServeFile(w, r, dataDir+"/favicon.ico")
		return
	}
	ErrNotFoundHandler(w, r)
}
