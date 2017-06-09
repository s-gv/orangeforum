package models

import (
	"errors"
	"net/http"
	"database/sql"
	"time"
	"github.com/s-gv/orangeforum/models/db"
	"encoding/hex"
	"log"
	"golang.org/x/crypto/bcrypt"
)

type Session struct {
	SessionID string
	UserID sql.NullInt64
	CSRFToken string
	Msg string
	CreatedDate time.Time
	UpdatedDate time.Time
}

const maxSessionLife = 200*time.Hour
const maxSessionLifeBeforeUpdate = 100*time.Hour

var ErrAuthFail = errors.New("username / password incorrect")
var ErrNoFlashMsg = errors.New("No flash message")

func Authenticate() error {
	return nil
}

func OpenSession(w http.ResponseWriter, r *http.Request) Session {
	cookie, err := r.Cookie("sessionid")
	if err == nil {
		sessionId := cookie.Value
		query := db.QueryRow(`SELECT sessionid, userid, csrf, msg, created_date, updated_date FROM sessions WHERE sessionid=?;`, sessionId)
		sess := Session{}
		var cDate int64
		var uDate int64
		if err := db.ScanRow(query, &sess.SessionID, &sess.UserID, &sess.CSRFToken, &sess.Msg, &cDate, &uDate); err == nil {
			sess.CreatedDate = time.Unix(cDate, 0)
			sess.UpdatedDate = time.Unix(uDate, 0)
			if sess.UpdatedDate.After(time.Now().Add(-maxSessionLife)) {
				if sess.UpdatedDate.Before(time.Now().Add(-maxSessionLifeBeforeUpdate)) {
					nowDate := int64(time.Now().Unix())
					db.Exec(`UPDATE sessions SET updated_date=? WHERE sessionid=?;`, nowDate, sessionId)
				}
				return sess
			} else {
				//log.Printf("[INFO] Session %s and last update date %s has expired.\n", sess.SessionID, sess.UpdatedDate)
			}
		} else {
			//log.Printf("[INFO] Session %s not found. %s\n", sess.SessionID, err)
		}
	}

	sess := Session{RandSeq(32), sql.NullInt64{}, RandSeq(32), "", time.Now(), time.Now()}
	db.Exec(`INSERT INTO sessions(sessionid, userid, csrf, msg, created_date, updated_date) values(?, ?, ?, ?, ?, ?);`,
		sess.SessionID, sess.UserID, sess.CSRFToken, sess.Msg, int64(sess.CreatedDate.Unix()), int64(sess.UpdatedDate.Unix()))
	db.Exec(`DELETE FROM sessions WHERE updated_date < ?;`, int64(sess.UpdatedDate.Unix()))

	http.SetCookie(w, &http.Cookie{Name: "sessionid", Value: sess.SessionID, HttpOnly: true})
	http.SetCookie(w, &http.Cookie{Name: "csrftoken", Value: sess.CSRFToken})

	return sess
}

func (sess *Session) SetFlashMsg(msg string) {
	db.Exec(`UPDATE sessions SET msg=? WHERE sessionid=?;`, msg, sess.SessionID)
}

func (sess *Session) FlashMsg() string {
	msg := sess.Msg
	sess.Msg = ""
	db.Exec(`UPDATE sessions SET msg=? WHERE sessionid=?;`, "", sess.SessionID)
	return msg
}

func (sess *Session) Authenticate(userName string, passwd string) bool {
	r := db.QueryRow(`SELECT id, passwdhash FROM users WHERE username=?;`, userName)
	var passwdHashStr string
	var userID int
	if err := db.ScanRow(r, &userID, &passwdHashStr); err != nil {
		return false
	}
	passwdHash, err := hex.DecodeString(passwdHashStr)
	if err != nil {
		log.Fatalf("[ERROR] Error in converting password hash from hex to byte slice: %s\n", err)
	}
	if err := bcrypt.CompareHashAndPassword(passwdHash, []byte(passwd)); err != nil {
		return false
	}
	sess.UserID = sql.NullInt64{int64(userID), true}
	db.Exec(`UPDATE sessions SET userid=? WHERE sessionid=?;`, sess.UserID, sess.SessionID)
	return true
}

func (sess *Session) IsUserValid() bool {
	return sess.UserID.Valid
}

func (sess *Session) IsUserSuperAdmin() bool {
	if sess.IsUserValid() {
		row := db.QueryRow(`SELECT is_superadmin FROM users WHERE id=?;`, sess.UserID)
		IsSuperAdmin := false
		if err := db.ScanRow(row, &IsSuperAdmin); err == nil {
			return IsSuperAdmin
		}
	}
	return false
}

func (sess *Session) UserName() (string, error) {
	if sess.UserID.Valid {
		r := db.QueryRow(`SELECT username FROM users WHERE id=?;`, sess.UserID)
		var userName string
		if err := db.ScanRow(r, &userName); err == nil {
			return userName, nil
		}
	}
	return "", errors.New("Invalid user")
}

func ClearSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sessionid")
	if err == nil {
		sessionID := cookie.Value
		db.Exec(`DELETE FROM sessions WHERE sessionid=?;`, sessionID)
	}
	http.SetCookie(w, &http.Cookie{Name: "sessionid", Value: "", Expires: time.Now().Add(-300*time.Hour), HttpOnly: true})
	http.SetCookie(w, &http.Cookie{Name: "csrftoken", Value: "", Expires: time.Now().Add(-300*time.Hour)})
}