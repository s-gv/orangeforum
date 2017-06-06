package models

import (
	"time"
	"github.com/s-gv/orangeforum/models/db"
	"math/rand"
	"errors"
)

const (
	VoteUp = 1
	VoteDown = 2
	VoteFlag = 3
)

const ModelVersion = 1

const (
	Secret string = "secret"
	ForumName string = "title"
	HeaderMsg string = "header_msg"
	SignupNeedsApproval string = "signup_needs_approval"
	PublicViewDisabled string = "public_view_disabled"
	SignupDisabled string = "signup_disabled"
	ImageUploadEnabled string = "image_upload_enabled"
	FileUploadEnabled string = "file_upload_enabled"
	AllowGroupSubscription string = "allow_group_subscription"
	AllowTopicSubscription string = "allow_topic_subscription"
	AutoSubscribeToMyTopic string = "auto_subscribe_to_my_topic"
)

var ErrIncorrectPasswd = errors.New("Incorrect username/password.")
var ErrUserNotFound = errors.New("Username not found.")

type User struct {
	ID int
	Username string
	PasswdHash string
	Email string
	About string
	Karma int
	IsBanned bool
	IsWarned bool
	IsSuperAdmin bool
	IsSuperMod bool
	IsApproved bool
	CreatedDate time.Time
	UpdatedDate time.Time
}

type Group struct {
	ID int
	Name string
	Desc string
	IsSticky string
	IsPrivate string
	IsClosed string
	HeaderMessage string
	CreatedDate time.Time
	UpdatedDate time.Time
}

type Topic struct {
	ID int
	Content string
	AuthorID int
	GroupID int
	IsDeleted bool
	IsSticky bool
	IsClosed bool
	NumComments int
	Upvotes int
	Downvotes int
	Flagvotes int
	CreatedDate time.Time
	UpdatedDate time.Time
}

type Comment struct {
	ID int
	Content string
	AuthorID int
	TopicID int
	ParentID int
	IsDeleted bool
	IsSticky bool
	Upvotes int
	Downvotes int
	Flagvotes int
	CreatedDate time.Time
	UpdatedDate time.Time
}

type Mod struct {
	ID int
	UserID int
	GroupID int
	CreatedDate time.Time
}

type Admin struct {
	ID int
	UserID int
	GroupID int
	CreatedDate time.Time
}

type TopicVote struct {
	ID int
	UserID int
	TopicID int
	VoteType int
	CreatedDate time.Time
}

type CommentVote struct {
	ID int
	UserID int
	CommentID int
	VoteType int
	CreatedDate time.Time
}

type TopicSubscription struct {
	ID int
	UserID int
	TopicID int
	CreatedDate time.Time
}

type GroupSubscription struct {
	ID int
	UserID int
	GroupID int
	CreatedDate time.Time
}

type ExtraNote struct {
	ID int
	Name string
	Content string
	URL string
	CreatedDate time.Time
	UpdatedDate time.Time
}

/*
func CreateUser(userName string, passwd string, email string) {
	if passwdHash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost); err == nil {
		db.CreateUser(userName, hex.EncodeToString(passwdHash), email)
	} else {
		log.Printf("[ERROR] %s.\n", err)
	}
}

func ProbeUser(userName string) bool {
	return db.ProbeUser(userName)
}

func Authenticate(userName string, passwd string) (User, error) {
	u, err := ReadUserByName(userName)
	if err != nil {
		return User{}, err
	}
	passwdHash, err := hex.DecodeString(u.PasswdHash)
	if err != nil {
		log.Fatalf("[ERROR] Error decoding password has from hex. %s\n", err)
	}
	if err := bcrypt.CompareHashAndPassword(passwdHash, []byte(passwd)); err != nil {
		return User{}, ErrIncorrectPasswd
	}
	return u, nil
}

func ReadUserByName(userName string) (User, error) {
	if row, err := db.QueryRow("ReadUserByName", `SELECT * FROM users WHERE username=?;`, userName); err == nil {
		u := User{}
		var cDate int64
		var uDate int64
		if err := row.Scan(&u.ID, &u.Username, &u.PasswdHash, &u.Email, &u.About, &u.Karma,
				&u.IsBanned, &u.IsWarned, &u.IsSuperAdmin, &u.IsSuperMod, &u.IsApproved,
				&cDate, &uDate); err == nil {
			u.CreatedDate = time.Unix(cDate, 0)
			u.UpdatedDate = time.Unix(uDate, 0)
			return u, nil
		} else {
			log.Println(err)
		}
	}
	return User{}, ErrUserNotFound
}

func ReadUserByID(userID int) (User, error) {
	if row, err := db.QueryRow("ReadUserByName", `SELECT * FROM users WHERE id=?;`, userID); err == nil {
		u := User{}
		var cDate int64
		var uDate int64
		if err := row.Scan(&u.ID, &u.Username, &u.PasswdHash, &u.Email, &u.About, &u.Karma,
			&u.IsBanned, &u.IsWarned, &u.IsSuperAdmin, &u.IsSuperMod, &u.IsApproved,
			&cDate, &uDate); err == nil {
			u.CreatedDate = time.Unix(cDate, 0)
			u.UpdatedDate = time.Unix(uDate, 0)
			return u, nil
		} else {
			log.Println(err)
		}
	}
	return User{}, ErrUserNotFound
}


func CreateUser(userName string, passwdHash string, email string) {
	now := int64(time.Now().Unix())
	exec("CreateUser", `INSERT INTO
			users(username, passwdhash, email, created_date, updated_date) values(?, ?, ?, ?, ?);`, userName, passwdHash, email, now, now)
}


func ProbeUser(userName string) bool {
	if row, err := QueryRow("ProbeUser", `SELECT username FROM users WHERE username=?;`, userName); err == nil {
		var tmp string
		if row.Scan(&tmp) == nil {
			return true
		}
	}
	return false
}
*/


func RandSeq(n int) string {
	var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func WriteConfig(key string, val string) {
	db.Exec(`INSERT OR REPLACE INTO configs(key, val) values(?, ?);`, key, val)
}


func Config(key string) string {
	row := db.QueryRow(`SELECT val FROM configs WHERE key=?;`, key)
	var val string
	if err := row.Scan(&val); err == nil {
		return val
	}
	return "0"
}

func Migrate() {
	db.CreateTables()

	WriteConfig("version", "1");
	WriteConfig(HeaderMsg, "")
	WriteConfig(ForumName, "OrangeForum")
	WriteConfig(Secret, RandSeq(32))
	WriteConfig(SignupNeedsApproval, "0")
	WriteConfig(PublicViewDisabled, "0")
	WriteConfig(SignupDisabled, "0")
	WriteConfig(FileUploadEnabled, "0")
	WriteConfig(ImageUploadEnabled, "0")
	WriteConfig(AllowGroupSubscription, "0")
	WriteConfig(AllowTopicSubscription, "0")
	WriteConfig(AutoSubscribeToMyTopic, "0")
}

func IsMigrationNeeded() bool {
	dbver := db.DBVersion()
	return dbver != ModelVersion

}