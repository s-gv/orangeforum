package models

import "time"

const (
	VoteUp = 1
	VoteDown = 2
	VoteFlag = 3
)

const ModelVersion = 1

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
	Secret string
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

type Session struct {
	ID int
	SessionID string
	UserID string
	Msg string
	Data string
	CreatedDate time.Time
	UpdatedDate time.Time
}