package models

import "time"

const (
	VoteUp = 1
	VoteDown = 2
	VoteFlag = 3
)

type User struct {
	ID int
	Username string
	PasswdHash string
	Email string
	About string
	Karma int
	IsBanned bool
	IsWarned bool
	IsAdmin bool
	IsSuperMod bool
	CreatedDate time.Time
	UpdatedDate time.Time
}

type Category struct {
	ID int
	Name string
	Desc string
	CreatedDate time.Time
	UpdatedDate time.Time
}

type Mod struct {
	ID int
	UserID int
	CategoryID int
	CreatedDate time.Time
}

type Topic struct {
	ID int
	Content string
	AuthorID int
	CategoryID int
	IsDeleted bool
	IsClosed bool
	IsSticky bool
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
	Upvotes int
	Downvotes int
	Flagvotes int
	CreatedDate time.Time
	UpdatedDate time.Time
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

type Session struct {
	ID int
	SessionID int
	UserID int
	Data string
	CreatedDate time.Time
	UpdateDate time.Time
}