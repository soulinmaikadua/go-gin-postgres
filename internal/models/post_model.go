package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `db:"id" json:"id" validate:"required,uuid"`
	Title     string    `db:"title" json:"title"`
	Details   string    `db:"details" json:"details"`
	IsPublish bool      `db:"is_publish" json:"is_publish"`
	UserId    uuid.UUID `db:"user_id" json:"user" validate:"required"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type UserOnPost struct {
	ID    uuid.UUID `db:"id" json:"id" validate:"required,uuid"`
	Name  string    `db:"name" json:"name"`
	Email string    `db:"email" json:"email" validate:"required"`
}

type ResponsePost struct {
	ID        uuid.UUID  `db:"id" json:"id" validate:"required,uuid"`
	Title     string     `db:"title" json:"title"`
	Details   string     `db:"details" json:"details"`
	IsPublish bool       `db:"is_publish" json:"is_publish"`
	User      UserOnPost `json:"user"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
}
