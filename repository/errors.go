package repository

import (
	"errors"
	"log"
	"net/http"

	"github.com/SawitProRecruitment/UserService/tools"
	"github.com/lib/pq"
)

// Constraint Key
const (
	UserPhoneNumberUniqueConstraint         = "users_unique_phone_number_key"
	UserPhoneNumberPasswordUniqueConstraint = "users_unique_phone_number_password_key"
)

// Err Message
const PasswordAlreadyExistsErrMessage = "password is already exists"

var ConstraintErrorMessageMap = map[string]string{
	UserPhoneNumberUniqueConstraint:         PasswordAlreadyExistsErrMessage,
	UserPhoneNumberPasswordUniqueConstraint: PasswordAlreadyExistsErrMessage,
}

func ConvertPGError(err error) error {
	if err == nil {
		return nil
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		constraintName := pqErr.Constraint
		detail := pqErr.Detail
		log.Println(detail)
		switch pqErr.Code {
		// integrity violation code
		case "23000", "23001", "23502", "23503", "23505", "23514", "23P01":
			errorMessage := ConstraintErrorMessageMap[constraintName]
			return &tools.Err{
				Code:    http.StatusConflict,
				Message: errorMessage,
			}
		}
	}

	return &tools.Err{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	}
}
