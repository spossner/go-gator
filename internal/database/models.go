// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"database/sql"

	"github.com/google/uuid"
)

type Feed struct {
	ID            uuid.UUID
	Name          string
	Url           string
	UserID        uuid.UUID
	LastFetchedAt sql.NullTime
	CreatedAt     sql.NullTime
	UpdatedAt     sql.NullTime
}

type FeedFollow struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	FeedID    uuid.UUID
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
}

type Post struct {
	ID          uuid.UUID
	FeedID      uuid.UUID
	Title       string
	Url         string
	Description sql.NullString
	PublishedAt sql.NullTime
	CreatedAt   sql.NullTime
	UpdatedAt   sql.NullTime
}

type User struct {
	ID        uuid.UUID
	Name      string
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
}
