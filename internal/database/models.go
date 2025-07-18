// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Application struct {
	ID          uuid.UUID
	ApplicantID uuid.UUID
	JobID       uuid.UUID
	ResumeUrl   string
	Status      string
}

type Job struct {
	ID          uuid.UUID
	Title       string
	Description string
	Location    string
	Type        string
	Salary      string
	EmployerID  uuid.UUID
}

type RefreshToken struct {
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	ExpiresAt time.Time
	RevokedAt sql.NullTime
}

type User struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Name           string
	Email          string
	HashedPassword string
	Role           string
}
