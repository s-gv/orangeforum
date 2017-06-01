package sessions

import (
	"errors"
	"net/http"
	"log"
	"database/sql"
	"time"
	"github.com/s-gv/orangeforum/models/db"
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
			return sess
		} else {
			log.Printf("[ERROR] Session error with session ID %s. %s\n", sess.SessionID, err)
		}
	}
	sess := Session{db.RandSeq(32), sql.NullInt64{}, db.RandSeq(32), "", "", time.Now(), time.Now()}
	db.CreateSession(sess.SessionID, sess.UserID, sess.CSRFToken, sess.Msg, sess.Data, sess.CreatedDate, sess.UpdatedDate)
	http.SetCookie(w, &http.Cookie{Name: "sessionid", Value: sess.SessionID, HttpOnly: true})
	http.SetCookie(w, &http.Cookie{Name: "csrftoken", Value: sess.CSRFToken})
	return sess
}

func (sess *Session) SetFlashMsg(msg string) {
	db.UpdateSessionFlashMsg(sess.SessionID, msg)
}

func (sess *Session) FlashMsg() string {
	db.UpdateSessionFlashMsg(sess.SessionID, "")
	return sess.Msg
}