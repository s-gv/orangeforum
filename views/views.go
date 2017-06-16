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

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	defer ErrServerHandler(w, r)
	sess := models.OpenSession(w, r)
	if r.URL.Path != "/" {
		ErrNotFoundHandler(w, r)
		return
	}

	type Group struct {
		Name string
		Desc string
	}
	groups := []Group{}
	rows := db.Query(`SELECT name, desc FROM groups WHERE is_closed=0 ORDER BY is_sticky ASC, RANDOM() LIMIT 25;`)
	for rows.Next() {
		groups = append(groups, Group{})
		rows.Scan(&groups[len(groups)-1].Name, &groups[len(groups)-1].Desc)
	}
	sort.Slice(groups, func(i, j int) bool {return groups[i].Name < groups[j].Name})
	templates.Render(w, "index.html", map[string]interface{}{
		"Common": models.ReadCommonData(sess),
		"GroupCreationDisabled": models.Config(models.GroupCreationDisabled) == "1",
		"Groups": groups,
	})
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
		return errors.New("Name can contain only alphabets, numbers, and underscore.")
	}
	return nil
}

func GroupHandler(w http.ResponseWriter, r *http.Request) {
	defer ErrServerHandler(w, r)
	sess := models.OpenSession(w, r)
	if r.Method == "POST" && r.PostFormValue("csrf") != sess.CSRFToken {
		ErrForbiddenHandler(w, r)
		return
	}

	name := r.FormValue("name")
	id := models.ReadGroupIDByName(name)
	if id == "" || models.ReadGroupIsDeleted(id) {
		ErrNotFoundHandler(w, r)
		return
	}

	templates.Render(w, "groupindex.html", map[string]interface{}{
		"Common": models.ReadCommonData(sess),
		"Name": name,
		"ID": id,
		"IsAdmin": sess.IsUserValid() && models.IsUserGroupAdmin(strconv.Itoa(int(sess.UserID.Int64)), id),
	})
}

func GroupEditHandler(w http.ResponseWriter, r *http.Request) {
	defer ErrServerHandler(w, r)
	sess := models.OpenSession(w, r)
	if r.Method == "POST" && r.PostFormValue("csrf") != sess.CSRFToken {
		ErrForbiddenHandler(w, r)
		return
	}
	if models.Config(models.GroupCreationDisabled) == "1" {
		ErrForbiddenHandler(w, r)
		return
	}
	if !sess.IsUserValid() {
		http.Redirect(w, r, "/login?next=/groups/edit", http.StatusSeeOther)
		return
	}

	userName, _ := sess.UserName()

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
	isUserInAdminList := false
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
		if !models.IsUserGroupAdmin(strconv.Itoa(int(sess.UserID.Int64)), groupID) {
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
			db.Exec(`UPDATE groups SET name=?, desc=?, header_msg=?, is_sticky=? WHERE id=?`, name, desc, headerMsg, isSticky, groupID)
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
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	defer ErrServerHandler(w, r)
	sess := models.OpenSession(w, r)
	if r.Method == "POST" && r.PostFormValue("csrf") != sess.CSRFToken {
		ErrForbiddenHandler(w, r)
		return
	}

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
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	defer ErrServerHandler(w, r)
	sess := models.OpenSession(w, r)
	if r.Method == "POST" && r.PostFormValue("csrf") != sess.CSRFToken {
		ErrForbiddenHandler(w, r)
		return
	}

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
}

func ChangePasswdHandler(w http.ResponseWriter, r *http.Request) {
	defer ErrServerHandler(w, r)
	sess := models.OpenSession(w, r)
	if r.Method == "POST" && r.PostFormValue("csrf") != sess.CSRFToken {
		ErrForbiddenHandler(w, r)
		return
	}

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
}

func ForgotPasswdHandler(w http.ResponseWriter, r *http.Request) {
	defer ErrServerHandler(w, r)
	sess := models.OpenSession(w, r)
	if r.Method == "POST" && r.PostFormValue("csrf") != sess.CSRFToken {
		ErrForbiddenHandler(w, r)
		return
	}

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
}

func validatePasswd(passwd string, passwdConfirm string) error {
	if len(passwd) < 8 {
		return errors.New("Password should have at least 8 characters.")
	}
	if passwd != passwdConfirm {
		return errors.New("Passwords don't match.")
	}
	return nil
}

func ResetPasswdHandler(w http.ResponseWriter, r *http.Request) {
	defer ErrServerHandler(w, r)
	sess := models.OpenSession(w, r)
	if r.Method == "POST" && r.PostFormValue("csrf") != sess.CSRFToken {
		ErrForbiddenHandler(w, r)
		return
	}

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
}

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

func CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	defer ErrServerHandler(w, r)
	sess := models.OpenSession(w, r)
	if r.Method == "POST" && r.PostFormValue("csrf") != sess.CSRFToken {
		ErrForbiddenHandler(w, r)
		return
	}
	if !sess.IsUserValid() {
		http.Redirect(w, r, "/login?next=/creategroup", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		groupName := r.PostFormValue("name")
		//groupDesc := r.PostFormValue("desc")
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
}

func AdminIndexHandler(w http.ResponseWriter, r *http.Request) {
	defer ErrServerHandler(w, r)
	sess := models.OpenSession(w, r)
	if r.Method == "POST" && r.PostFormValue("csrf") != sess.CSRFToken {
		ErrForbiddenHandler(w, r)
		return
	}
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
}

func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	defer ErrServerHandler(w, r)
	sess := models.OpenSession(w, r)
	if r.Method == "POST" && r.PostFormValue("csrf") != sess.CSRFToken {
		ErrForbiddenHandler(w, r)
		return
	}

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
}

func NoteHandler(w http.ResponseWriter, r *http.Request) {
	defer ErrServerHandler(w, r)
	sess := models.OpenSession(w, r)
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
}

func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	defer ErrServerHandler(w, r)
	dataDir := models.Config(models.DataDir)
	if dataDir != "" {
		http.ServeFile(w, r, dataDir+"/favicon.ico")
		return
	}
	ErrNotFoundHandler(w, r)
}
