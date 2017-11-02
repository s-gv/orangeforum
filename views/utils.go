package views

import (
	"net/http"
	"log"
	"runtime/debug"
	"github.com/s-gv/orangeforum/models"
	"net/url"
	"regexp"
	"strings"
	"path/filepath"
	"os"
	"io"
	"html/template"
	"time"
	"strconv"
	"errors"
)

var linkRe *regexp.Regexp

func init() {
	linkRe = regexp.MustCompile("https?://[^\\s]+[A-Za-z0-9/\\&\\+\\?#,_-]")
}

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

func validatePasswd(passwd string, passwdConfirm string) error {
	if len(passwd) < 8 {
		return errors.New("Password should have at least 8 characters.")
	}
	if passwd != passwdConfirm {
		return errors.New("Passwords don't match.")
	}
	return nil
}
