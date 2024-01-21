package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidValidation(t *testing.T) {
	funcZ := func() {}

	err := ValidateRequestPayload(funcZ)

	assert.Error(t, err)
}

func TestValidatePhoneNumberVal(t *testing.T) {
	type TestPayload struct {
		PhoneNumber string `json:"phone_number" validate:"phone_number"`
	}
	t.Run("pass empty string ", func(t *testing.T) {
		p := TestPayload{
			PhoneNumber: "",
		}
		err := ValidateRequestPayload(p)
		assert.NoError(t, err)
	})
	t.Run("must be at minimum 10 characters and maximum 13 characters", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			p := TestPayload{
				PhoneNumber: "+62345678901",
			}

			err := ValidateRequestPayload(p)
			assert.NoError(t, err)
		})
		t.Run("fail", func(t *testing.T) {
			p := TestPayload{
				PhoneNumber: "+623456789012345",
			}

			err := ValidateRequestPayload(p)
			assert.Error(t, err)
		})
	})
	t.Run("must start with the Indonesia country code “+62”", func(t *testing.T) {
		t.Run("fail", func(t *testing.T) {
			p := TestPayload{
				PhoneNumber: "+44345678901",
			}

			err := ValidateRequestPayload(p)
			assert.Error(t, err)
		})
	})
}

func TestValidatePasswordVal(t *testing.T) {
	type TestPayload struct {
		Password string `json:"password" validate:"pwd"`
	}
	t.Run("containing at least 1 capital characters AND 1 number AND 1 special (non alpha-numeric) characters", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			p := TestPayload{
				Password: "IloveVirginCo2Nut123$",
			}

			err := ValidateRequestPayload(p)
			assert.NoError(t, err)
		})
		t.Run("fail", func(t *testing.T) {
			p := TestPayload{
				Password: "password",
			}
			err := ValidateRequestPayload(p)
			assert.Error(t, err)

			p = TestPayload{
				Password: "",
			}
			err = ValidateRequestPayload(p)
			assert.Error(t, err)
		})
	})
}
