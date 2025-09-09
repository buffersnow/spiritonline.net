package db

import (
	"database/sql"
	"time"
)

type Account struct {
	UserID        int64     `db:"user_id"`
	Screenname    string    `db:"screenname"`
	Mail          string    `db:"mail"`
	Flags         int64     `db:"flags"`
	LastUpdatedOn time.Time `db:"last_updated_on"`

	User         User          `db:"-"`
	Profile      Profile       `db:"-"`
	AppPasswords []AppPassword `db:"-"`
	Contacts     []Contact     `db:"-"`
}

type User struct {
	UserID         int64          `db:"user_id"`
	AvatarHash     sql.NullString `db:"avatar_hash"`
	ContactListRev int64          `db:"contact_list_rev"`
	LastCLUpdate   sql.NullTime   `db:"last_cl_update"`
	LastLogin      sql.NullTime   `db:"last_login"`
	FirstLogin     sql.NullTime   `db:"first_login"`
	SignupDate     time.Time      `db:"signup_date"`
	LastUpdatedOn  time.Time      `db:"last_updated_on"`
}

type Profile struct {
	UserID         int64          `db:"user_id"`
	StatusMessage  sql.NullString `db:"status_message"`
	AwayMessage    sql.NullString `db:"away_message"`
	OfflineMessage sql.NullString `db:"offline_message"`
	LastUpdatedOn  time.Time      `db:"last_updated_on"`
}

type Contact struct {
	SenderID      int64          `db:"sender_id"`
	FriendID      int64          `db:"friend_id"`
	Reason        sql.NullString `db:"reason"`
	IsBlocked     bool           `db:"is_blocked"`
	AddedOn       time.Time      `db:"added_on"`
	LastUpdatedOn time.Time      `db:"last_updated_on"`
}

type AppPassword struct {
	UserID        int64          `db:"user_id"`
	Type          string         `db:"type"`
	Contents      sql.NullString `db:"contents"`
	LastUpdatedOn time.Time      `db:"last_updated_on"`
}
