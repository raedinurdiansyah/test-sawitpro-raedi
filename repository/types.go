// This file contains types that are used in the repository layer.
package repository

import (
	"time"

	"github.com/google/uuid"
)

type RecordTimeStamp struct {
	CreatedAt      time.Time
	LastModifiedAt time.Time
	DeletedAt      *time.Time
}

type User struct {
	ID          int
	GUID        uuid.UUID
	FullName    string
	PhoneNumber string
	Password    string

	RecordTimeStamp
}

type CreateUserInput struct {
	PhoneNumber string
	FullName    string
	Password    string
}

type LoginUserInput struct {
	PhoneNumber string
	Password    string
}

type LoginUserOutput struct {
	GUID     uuid.UUID
	FullName string
	Password string
}

type GetUserByGUIDOutput struct {
	ID          int
	CreatedAt   time.Time
	FullName    string
	PhoneNumber string
}
