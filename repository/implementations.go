package repository

import (
	"context"

	"github.com/google/uuid"
)

func (r *Repository) CreateUser(ctx context.Context, user *User) (err error) {
	err = r.Db.QueryRowContext(
		ctx,
		"INSERT INTO users (full_name, phone_number, password) VALUES ($1, $2, $3) RETURNING *",
		user.FullName,
		user.PhoneNumber,
		user.Password,
	).Scan(
		&user.ID,
		&user.GUID,
		&user.FullName,
		&user.PhoneNumber,
		&user.Password,
		&user.CreatedAt,
		&user.LastModifiedAt,
		&user.DeletedAt,
	)

	return ConvertPGError(err)
}

func (r *Repository) GetUserLoginByPhoneNumber(ctx context.Context, phoneNumber string) (
	output LoginUserOutput,
	err error,
) {
	err = r.Db.QueryRowContext(
		ctx,
		"SELECT guid, full_name, password FROM users WHERE phone_number = $1 AND deleted_at IS NULL",
		phoneNumber,
	).Scan(&output.GUID, &output.FullName, &output.Password)
	if err != nil {
		return
	}
	return
}

func (r *Repository) GetUserByGUID(ctx context.Context, guid uuid.UUID) (user *User, err error) {
	user = new(User)

	err = r.Db.QueryRowContext(
		ctx,
		"SELECT * FROM users WHERE guid = $1 AND deleted_at IS NULL",
		guid,
	).Scan(
		&user.ID,
		&user.GUID,
		&user.FullName,
		&user.PhoneNumber,
		&user.Password,
		&user.CreatedAt,
		&user.LastModifiedAt,
		&user.DeletedAt,
	)
	if err != nil {
		return
	}
	return
}

func (r *Repository) UpdateUser(ctx context.Context, user *User) error {
	err := r.Db.QueryRowContext(
		ctx,
		"UPDATE users SET full_name = $1, phone_number = $2 , last_modified_at = NOW() WHERE id = $3 RETURNING *",
		user.FullName, user.PhoneNumber, user.ID,
	).Scan(
		&user.ID,
		&user.GUID,
		&user.FullName,
		&user.PhoneNumber,
		&user.Password,
		&user.CreatedAt,
		&user.LastModifiedAt,
		&user.DeletedAt,
	)
	if err != nil {
		return ConvertPGError(err)
	}
	return nil
}
