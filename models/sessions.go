package models

import (
	"errors"
	"net/http"
	"database/sql"
	"time"
	"github.com/s-gv/orangeforum/models/db"
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

func (sess *Session) SetUser(user User) {
	println(user.ID)
	sess.UserID = sql.NullInt64{int64(user.ID), true}
	db.Exec(`UPDATE sessions SET userid=? WHERE sessionid=?;`, sess.UserID, sess.SessionID)
}
/*
func (sess *Session) User() (models.User, error) {

	if sess.UserID.Valid {
		println("User in session valid")
		if u, err := models.ReadUserByID(int(sess.UserID.Int64)); err == nil {
			return u, nil
		}
	} else {
		println("Userid not valid in session")
	}
	return models.User{}, errors.New("Invalid user")

}
*/