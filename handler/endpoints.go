package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/tools"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (s *Server) RegisterUser(ctx echo.Context) error {
	var input generated.RegisterUserJSONRequestBody
	var errData *tools.Err

	err := ctx.Bind(&input)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: err.Error()})
	}

	err = tools.ValidateRequestPayload(input)
	if err != nil && errors.As(err, &errData) {
		return ctx.JSON(errData.Code, generated.ErrorWithExtraResponse{Message: errData.Message, Extra: &errData.Extra})
	}

	hashedPassword, err := tools.HashPassword(input.Password)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: err.Error()})
	}

	newUser := &repository.User{
		PhoneNumber: input.PhoneNumber,
		FullName:    input.FullName,
		Password:    hashedPassword,
	}

	err = s.Repository.CreateUser(ctx.Request().Context(), newUser)
	if err != nil {
		if errors.As(err, &errData) {
			return ctx.JSON(errData.Code, generated.ErrorResponse{Message: errData.Message})
		}
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	resp := generated.SuccessRegisterUserResponse{
		GUID:    newUser.GUID,
		Message: "User registered successfully",
	}
	return ctx.JSON(http.StatusCreated, resp)
}

func (s *Server) LoginUser(ctx echo.Context) error {
	var input generated.LoginUserJSONRequestBody
	var errData *tools.Err

	err := ctx.Bind(&input)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: err.Error()})
	}

	err = tools.ValidateRequestPayload(input)
	if err != nil && errors.As(err, &errData) {
		return ctx.JSON(errData.Code, generated.ErrorWithExtraResponse{Message: errData.Message, Extra: &errData.Extra})
	}

	output, err := s.Repository.GetUserLoginByPhoneNumber(ctx.Request().Context(), input.PhoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, generated.ErrorResponse{Message: "user is not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	isMatched := tools.IsValidPassword(output.Password, input.Password)
	if !isMatched {
		return ctx.JSON(http.StatusBadRequest, "invalid password")
	}

	tokenParam := tools.GenerateJWTTokenParams{
		FullName: output.FullName,
		GUID:     output.GUID,
	}

	token, expiredAt, err := tools.GenerateJWTToken(tokenParam, s.Config.JWTTokenLifetimeInHours, s.Config.RSAPrivateKey)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	resp := generated.SuccessLoginUserResponse{
		Token:     token,
		ExpiredAt: expiredAt.String(),
	}
	return ctx.JSON(http.StatusOK, resp)

}

func (s *Server) GetUserProfile(ctx echo.Context) error {
	guidString := ctx.Get("UserGUID").(string)

	user, err := s.Repository.GetUserByGUID(ctx.Request().Context(), uuid.MustParse(guidString))
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, generated.ErrorResponse{Message: "user is not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	resp := generated.SuccessGetUserProfileResponse{
		Guid:        user.GUID,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   &user.CreatedAt,
	}
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) UpdateUser(ctx echo.Context) error {
	rCtx := ctx.Request().Context()
	guid := uuid.MustParse(ctx.Get("UserGUID").(string))
	var errData *tools.Err

	var input generated.UpdateUserJSONRequestBody
	err := ctx.Bind(&input)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: err.Error()})
	}

	err = tools.ValidateRequestPayload(input)
	if err != nil && errors.As(err, &errData) {
		return ctx.JSON(errData.Code, generated.ErrorWithExtraResponse{Message: errData.Message, Extra: &errData.Extra})
	}

	user, err := s.Repository.GetUserByGUID(rCtx, guid)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, generated.ErrorResponse{Message: "user is not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	user.FullName = input.FullName
	user.PhoneNumber = input.PhoneNumber

	err = s.Repository.UpdateUser(rCtx, user)
	if err != nil {
		if errors.As(err, &errData) {
			return ctx.JSON(errData.Code, generated.ErrorResponse{Message: errData.Message})
		}
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, generated.DefaultUpdateResponse{Message: "user updated successfully"})
}
