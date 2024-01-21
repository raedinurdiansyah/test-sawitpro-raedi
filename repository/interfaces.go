// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import (
	"context"

	"github.com/google/uuid"
)

type RepositoryInterface interface {
	CreateUser(ctx context.Context, user *User) (err error)
	GetUserLoginByPhoneNumber(ctx context.Context, phoneNumber string) (output LoginUserOutput, err error)
	GetUserByGUID(ctx context.Context, guid uuid.UUID) (user *User, err error)
	UpdateUser(ctx context.Context, user *User) error
}
