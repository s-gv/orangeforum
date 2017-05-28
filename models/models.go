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
	CreatedDate time.Time
	UpdatedDate time.Time
}

type SubForum struct {
	ID int
	Name string
	Desc string
	CreatedDate time.Time
	UpdatedDate time.Time
}

type Mod struct {
	ID int
	UserID int
	SubforumID int
	CreatedDate time.Time
}

type Topic struct {
	ID int
	Content string
	SubForumID int
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
	Upvotes int
	Downvotes int
	Flagvotes int
	CreatedDate time.Time
	UpdatedDate time.Time
}

type TopicVote struct {
	ID int
	AuthorID int
	TopicID int
	VoteType int
	CreatedDate time.Time
}

type CommentVote struct {
	ID int
	AuthorID int
	CommentID int
	VoteType int
	CreatedDate time.Time
}