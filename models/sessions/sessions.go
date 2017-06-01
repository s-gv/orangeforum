package sessions

import (
	"errors"
	"net/http"
	"log"
	"database/sql"
	"time"
	"github.com/s-gv/orangeforum/models/db"
	"github.com/s-gv/orangeforum/models"
)

type Session struct {
	SessionID string
	UserID sql.NullInt64
	CSRFToken string
	Msg string
	Data string
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

func Open(w http.ResponseWriter, r *http.Request) Session {
	cookie, err := r.Cookie("sessionid")
	if err == nil {
		sess := Session{}
		sess.SessionID = cookie.Value
		if err := db.ReadSession(sess.SessionID, &sess.UserID, &sess.CSRFToken, &sess.Msg, &sess.Data, &sess.CreatedDate, &sess.UpdatedDate); err == nil {
			if sess.UpdatedDate.After(time.Now().Add(-maxSessionLife)) {
				if sess.UpdatedDate.Before(time.Now().Add(-maxSessionLifeBeforeUpdate)) {
					db.UpdateSessionDate(sess.SessionID, time.Now())
				}
				return sess
			} else {
				log.Printf("[INFO] Session %s and last update date %s has expired.\n", sess.SessionID, sess.UpdatedDate)
			}
		} else {
			log.Printf("[INFO] Session %s not found. %s\n", sess.SessionID, err)
		}
	}

	sess := Session{db.RandSeq(32), sql.NullInt64{}, db.RandSeq(32), "", "", time.Now(), time.Now()}
	db.CreateSession(sess.SessionID, sess.UserID, sess.CSRFToken, sess.Msg, sess.Data, sess.CreatedDate, sess.UpdatedDate)
	http.SetCookie(w, &http.Cookie{Name: "sessionid", Value: sess.SessionID, HttpOnly: true})
	http.SetCookie(w, &http.Cookie{Name: "csrftoken", Value: sess.CSRFToken})

	db.DeleteSessions(time.Now().Add(-maxSessionLife))

	return sess
}

func (sess *Session) SetFlashMsg(msg string) {
	db.UpdateSessionFlashMsg(sess.SessionID, msg)
}

func (sess *Session) FlashMsg() string {
	db.UpdateSessionFlashMsg(sess.SessionID, "")
	msg := sess.Msg
	sess.Msg = ""
	return msg
}

func (sess *Session) SetUser(user models.User, valid bool) {
	println(user.ID)
	sess.UserID = sql.NullInt64{int64(user.ID), valid}
	db.UpdateSessionUserID(sess.SessionID, sess.UserID)
}

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