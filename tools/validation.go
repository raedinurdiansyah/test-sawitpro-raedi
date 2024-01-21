package tools

import (
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type ExtraValidation struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func init() {
	validate = validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	err := validate.RegisterValidation("pwd", ValidatePasswordVal)
	if err != nil {
		panic(err)
	}

	err = validate.RegisterValidation("phone_number", ValidatePhoneNumberVal)
	if err != nil {
		panic(err)
	}
}

func ValidatePhoneNumberVal(fl validator.FieldLevel) bool {
	v := fl.Field()
	if v.Kind() == reflect.String {
		if v.String() == "" {
			return true
		}
	}

	if len(v.String()) < 9 || len(v.String()) > 14 {
		return false
	}

	if !strings.HasPrefix(v.String(), "+62") {
		return false
	}

	regex := regexp.MustCompile(`^\+62[0-9]+$`)
	return regex.MatchString(v.String())
}

func ValidatePasswordVal(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	hasCapital := false
	hasNumber := false
	hasSpecial := false

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasCapital = true
		} else if unicode.IsNumber(char) {
			hasNumber = true
		} else if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			hasSpecial = true
		}
	}

	return hasCapital && hasNumber && hasSpecial
}

func ValidateRequestPayload(s interface{}) error {
	err := validate.Struct(s)

	if err == nil {
		return nil
	}

	if _, ok := err.(*validator.InvalidValidationError); ok {
		return err
	}

	extras := make([]ExtraValidation, 0)

	for _, err := range err.(validator.ValidationErrors) {
		extras = append(extras, ExtraValidation{
			Field:   normalizeNamespace(err.Namespace()),
			Message: err.Tag(),
		},
		)
	}

	return &Err{
		Code:    http.StatusBadRequest,
		Message: "invalid request payload values",
		Extra:   extras,
	}
}

// normalizeNamespace will remove upper case name as struct name
func normalizeNamespace(s string) string {
	splitStr := strings.Split(s, ".")
	var out []string
	for _, seg := range splitStr {
		if unicode.IsUpper(rune(seg[0])) {
			continue
		}
		out = append(out, seg)
	}
	return strings.Join(out, ".")
}
