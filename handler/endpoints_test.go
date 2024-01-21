package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/SawitProRecruitment/UserService/config"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/tools"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func MockUser() *repository.User {
	return &repository.User{
		ID:          1,
		GUID:        uuid.New(),
		FullName:    "SawitPro Mania",
		PhoneNumber: "+62345678901",
		Password:    "IloveVirginCo2Nut123$",
		RecordTimeStamp: repository.RecordTimeStamp{
			CreatedAt:      time.Now(),
			LastModifiedAt: time.Now(),
			DeletedAt:      nil,
		},
	}
}

func TestRegisterUser(t *testing.T) {
	t.Run("when success register user", func(t *testing.T) {
		e := echo.New()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		requestBody := map[string]interface{}{
			"full_name":    "SawitPro Mania",
			"password":     "IloveVirginCo2Nut123$",
			"phone_number": "+62345678901",
		}
		requestJSON, err := json.Marshal(requestBody)
		if err != nil {
			t.Error(err)
		}

		mockRepo := repository.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(requestJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		s := &Server{
			Repository: mockRepo,
		}

		err = s.RegisterUser(ctx)

		assert.NoError(t, err, "should be no error")
		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}

func TestRegisterUser_Error(t *testing.T) {
	t.Run("when error binding request body", func(t *testing.T) {
		e := echo.New()
		reqParam := testRequestEndpointParam{
			e:          e,
			httpMethod: http.MethodPost,
			token:      "",
			url:        "/register",
			body:       []byte("invalid json"),
		}
		ctx, rec := TestRequestEndpoint(reqParam)

		s := &Server{}

		_ = s.RegisterUser(ctx)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when error validate request body", func(t *testing.T) {
		e := echo.New()
		reqParam := testRequestEndpointParam{
			e:          e,
			httpMethod: http.MethodPost,
			token:      "",
			url:        "/register",
			body:       []byte(`{"full_name": "Sa", "password": "password", "phone_number": "+44345678901"}`),
		}
		ctx, rec := TestRequestEndpoint(reqParam)

		s := &Server{}

		_ = s.RegisterUser(ctx)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	// t.Run("when error to hash password", func(t *testing.T) {
	// 	ctx, rec := TestRequestEndpoint(http.MethodPost, "/register", []byte(`{"full_name": "SawitPro Mania", "password": "Password_that_causes_error_123!test_lorem_ipsum panjang", "phone_number": "+62345678901"}`))

	// 	s := &Server{}

	// 	_ = s.RegisterUser(ctx)
	// 	assert.Equal(t, http.StatusBadRequest, rec.Code)
	// })

	t.Run("when error to create user", func(t *testing.T) {
		e := echo.New()
		reqParam := testRequestEndpointParam{
			e:          e,
			httpMethod: http.MethodPost,
			token:      "",
			url:        "/register",
			body:       []byte(`{"full_name": "SawitPro Mania", "password": "IloveVirginCo2Nut123$", "phone_number": "+62345678901"}`),
		}
		ctx, rec := TestRequestEndpoint(reqParam)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := repository.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))

		s := &Server{
			Repository: mockRepo,
		}

		_ = s.RegisterUser(ctx)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when error to create user due to duplicate phone number", func(t *testing.T) {
		e := echo.New()
		reqParam := testRequestEndpointParam{
			e:          e,
			httpMethod: http.MethodPost,
			token:      "",
			url:        "/register",
			body:       []byte(`{"full_name": "SawitPro Mania", "password": "IloveVirginCo2Nut123$", "phone_number": "+62345678901"}`),
		}
		ctx, rec := TestRequestEndpoint(reqParam)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := repository.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(&tools.Err{Code: http.StatusBadRequest, Message: "duplicate phone number"})

		s := &Server{
			Repository: mockRepo,
		}

		_ = s.RegisterUser(ctx)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestLoginUser(t *testing.T) {
	t.Run("when success login user", func(t *testing.T) {
		e := echo.New()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		requestBody := `{
			"password":     "IloveVirginCo2Nut123$",
			"phone_number": "+62345678901"
		}`
		hashed, _ := tools.HashPassword("IloveVirginCo2Nut123$")
		mockOutput := repository.LoginUserOutput{
			GUID:     uuid.New(),
			FullName: "joz gandoz",
			Password: hashed,
		}
		mockRepo := repository.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().GetUserLoginByPhoneNumber(gomock.Any(), "+62345678901").Return(mockOutput, nil)

		reqParam := testRequestEndpointParam{
			e:          e,
			httpMethod: http.MethodPost,
			token:      "",
			url:        "/login",
			body:       []byte(requestBody),
		}
		ctx, rec := TestRequestEndpoint(reqParam)

		mockConfig := &config.Config{
			RSAPrivateKey:           tools.MockRSAPrivateKey(),
			JWTTokenLifetimeInHours: 8,
		}

		s := &Server{
			Repository: mockRepo,
			Config:     *mockConfig,
		}

		err := s.LoginUser(ctx)

		assert.NoError(t, err, "should be no error")
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestLoginUser_Error(t *testing.T) {
	t.Run("when error binding request body", func(t *testing.T) {
		e := echo.New()
		reqParam := testRequestEndpointParam{
			e:          e,
			httpMethod: http.MethodPost,
			token:      "",
			url:        "/login",
			body:       []byte("invalid json"),
		}
		ctx, rec := TestRequestEndpoint(reqParam)

		s := &Server{}

		_ = s.LoginUser(ctx)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when error validate request body", func(t *testing.T) {
		e := echo.New()
		reqParam := testRequestEndpointParam{
			e:          e,
			httpMethod: http.MethodPost,
			token:      "",
			url:        "/login",
			body:       []byte(`{"password": "password", "phone_number": "+44345678901"}`),
		}
		ctx, rec := TestRequestEndpoint(reqParam)

		s := &Server{}

		_ = s.LoginUser(ctx)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when error due to phone number is not found", func(t *testing.T) {
		e := echo.New()
		reqParam := testRequestEndpointParam{
			e:          e,
			httpMethod: http.MethodPost,
			token:      "",
			url:        "/login",
			body:       []byte(`{"password": "IloveVirginCo2Nut123$", "phone_number": "+62345678901"}`),
		}
		ctx, rec := TestRequestEndpoint(reqParam)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().GetUserLoginByPhoneNumber(gomock.Any(), "+62345678901").Return(repository.LoginUserOutput{}, sql.ErrNoRows)
		s := &Server{
			Repository: mockRepo,
		}

		_ = s.LoginUser(ctx)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when error to get user by phone number", func(t *testing.T) {
		e := echo.New()
		reqParam := testRequestEndpointParam{
			e:          e,
			httpMethod: http.MethodPost,
			token:      "",
			url:        "/login",
			body:       []byte(`{"password": "IloveVirginCo2Nut123$", "phone_number": "+62345678901"}`),
		}
		ctx, rec := TestRequestEndpoint(reqParam)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := repository.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().GetUserLoginByPhoneNumber(gomock.Any(), "+62345678901").Return(repository.LoginUserOutput{}, fmt.Errorf("error db"))
		s := &Server{
			Repository: mockRepo,
		}

		_ = s.LoginUser(ctx)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when error due to password is wrong", func(t *testing.T) {
		e := echo.New()
		reqParam := testRequestEndpointParam{
			e:          e,
			httpMethod: http.MethodPost,
			token:      "",
			url:        "/login",
			body:       []byte(`{"password": "IloveVirginCo2Nut123$", "phone_number": "+62345678901"}`),
		}
		ctx, rec := TestRequestEndpoint(reqParam)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		hashed, _ := tools.HashPassword("IloveCocoNut123$")
		mockOutput := repository.LoginUserOutput{
			GUID:     uuid.New(),
			FullName: "joz gandoz",
			Password: hashed,
		}
		mockRepo := repository.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().GetUserLoginByPhoneNumber(gomock.Any(), "+62345678901").Return(mockOutput, nil)
		s := &Server{
			Repository: mockRepo,
		}

		_ = s.LoginUser(ctx)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when error generating jwt token", func(t *testing.T) {
		e := echo.New()
		reqParam := testRequestEndpointParam{
			e:          e,
			httpMethod: http.MethodPost,
			token:      "",
			url:        "/login",
			body:       []byte(`{"password": "IloveVirginCo2Nut123$", "phone_number": "+62345678901"}`),
		}
		ctx, rec := TestRequestEndpoint(reqParam)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		hashed, _ := tools.HashPassword("IloveVirginCo2Nut123$")
		mockOutput := repository.LoginUserOutput{
			GUID:     uuid.New(),
			FullName: "joz gandoz",
			Password: hashed,
		}
		mockRepo := repository.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().GetUserLoginByPhoneNumber(gomock.Any(), "+62345678901").Return(mockOutput, nil)
		s := &Server{
			Repository: mockRepo,
		}

		_ = s.LoginUser(ctx)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestGetUserProfile(t *testing.T) {
	mockConfig := &config.Config{
		RSAPrivateKey:           tools.MockRSAPrivateKey(),
		JWTTokenLifetimeInHours: 60,
		RSAPublicKey:            tools.MockRSAPublicKey(),
	}

	t.Run("when success get user data", func(t *testing.T) {
		e := echo.New()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUser := MockUser()
		tokenParams := tools.GenerateJWTTokenParams{
			FullName: mockUser.FullName,
			GUID:     mockUser.GUID,
		}
		token, _, err := tools.GenerateJWTToken(tokenParams, 60, tools.MockRSAPrivateKey())
		assert.NoError(t, err)

		mockRepo := repository.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().GetUserByGUID(gomock.Any(), mockUser.GUID).Return(mockUser, nil)
		s := &Server{
			Repository: mockRepo,
			Config:     *mockConfig,
		}

		reqParam := testRequestEndpointParam{
			e:          e,
			httpMethod: http.MethodGet,
			token:      token,
			url:        "/users",
			body:       nil,
		}

		ctx, rec := TestRequestEndpoint(reqParam)
		ctx.Set("UserGUID", mockUser.GUID.String())

		_ = s.GetUserProfile(ctx)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestGetUserProfile_Error(t *testing.T) {
	mockConfig := &config.Config{
		RSAPrivateKey:           tools.MockRSAPrivateKey(),
		JWTTokenLifetimeInHours: 60,
		RSAPublicKey:            tools.MockRSAPublicKey(),
	}
	mockUser := MockUser()
	tokenParams := tools.GenerateJWTTokenParams{
		FullName: mockUser.FullName,
		GUID:     mockUser.GUID,
	}
	token, _, err := tools.GenerateJWTToken(tokenParams, 60, tools.MockRSAPrivateKey())
	assert.NoError(t, err)

	t.Run("when error to get user data", func(t *testing.T) {
		e := echo.New()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reqParam := testRequestEndpointParam{
			e:          e,
			httpMethod: http.MethodGet,
			token:      token,
			url:        "/users",
			body:       nil,
		}

		ctx, rec := TestRequestEndpoint(reqParam)
		ctx.Set("UserGUID", mockUser.GUID.String())

		mockRepo := repository.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().GetUserByGUID(gomock.Any(), mockUser.GUID).Return(nil, fmt.Errorf("error db"))
		s := &Server{
			Repository: mockRepo,
			Config:     *mockConfig,
		}
		_ = s.GetUserProfile(ctx)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when error due to user is not found", func(t *testing.T) {
		e := echo.New()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reqParam := testRequestEndpointParam{
			e:          e,
			httpMethod: http.MethodGet,
			token:      token,
			url:        "/users",
			body:       nil,
		}

		ctx, rec := TestRequestEndpoint(reqParam)
		ctx.Set("UserGUID", mockUser.GUID.String())

		mockRepo := repository.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().GetUserByGUID(gomock.Any(), mockUser.GUID).Return(nil, sql.ErrNoRows)
		s := &Server{
			Repository: mockRepo,
			Config:     *mockConfig,
		}
		_ = s.GetUserProfile(ctx)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestUpdateUser(t *testing.T) {
	mockConfig := &config.Config{
		RSAPrivateKey:           tools.MockRSAPrivateKey(),
		JWTTokenLifetimeInHours: 60,
		RSAPublicKey:            tools.MockRSAPublicKey(),
	}
	mockUser := MockUser()
	tokenParams := tools.GenerateJWTTokenParams{
		FullName: mockUser.FullName,
		GUID:     mockUser.GUID,
	}
	token, _, err := tools.GenerateJWTToken(tokenParams, 60, tools.MockRSAPrivateKey())
	assert.NoError(t, err)

	t.Run("when success update user data", func(t *testing.T) {
		e := echo.New()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		requestBody := `{
			"full_name":   "test update user",
			"phone_number": "+62345678901"
		}`
		mockRepo := repository.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().GetUserByGUID(gomock.Any(), mockUser.GUID).Return(mockUser, nil)

		mockUser.FullName = "test update user"
		mockUser.PhoneNumber = "+62345678901"

		mockRepo.EXPECT().UpdateUser(gomock.Any(), mockUser).Return(nil)

		reqParam := testRequestEndpointParam{
			e:          e,
			httpMethod: http.MethodPut,
			token:      token,
			url:        "/users",
			body:       []byte(requestBody),
		}
		ctx, rec := TestRequestEndpoint(reqParam)
		ctx.Set("UserGUID", mockUser.GUID.String())

		s := &Server{
			Repository: mockRepo,
			Config:     *mockConfig,
		}

		_ = s.UpdateUser(ctx)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestUpdateUser_Error(t *testing.T) {
	successReqBody := `{
		"full_name":   "test update user",
		"phone_number": "+62345678901"
	}`
	mockConfig := &config.Config{
		RSAPrivateKey:           tools.MockRSAPrivateKey(),
		JWTTokenLifetimeInHours: 60,
		RSAPublicKey:            tools.MockRSAPublicKey(),
	}
	mockUser := MockUser()
	tokenParams := tools.GenerateJWTTokenParams{
		FullName: mockUser.FullName,
		GUID:     mockUser.GUID,
	}
	token, _, err := tools.GenerateJWTToken(tokenParams, 60, tools.MockRSAPrivateKey())
	assert.NoError(t, err)

	t.Run("when error binding request body", func(t *testing.T) {
		e := echo.New()
		reqParam := testRequestEndpointParam{
			e:          e,
			httpMethod: http.MethodPut,
			token:      token,
			url:        "/users",
			body:       []byte("invalid json"),
		}
		ctx, rec := TestRequestEndpoint(reqParam)
		ctx.Set("UserGUID", mockUser.GUID.String())

		s := &Server{}

		_ = s.UpdateUser(ctx)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when error validate request body", func(t *testing.T) {
		e := echo.New()
		reqParam := testRequestEndpointParam{
			e:          e,
			httpMethod: http.MethodPut,
			token:      token,
			url:        "/users",
			body:       []byte(`{"password": "password", "phone_number": "+44345678901"}`),
		}
		ctx, rec := TestRequestEndpoint(reqParam)
		ctx.Set("UserGUID", mockUser.GUID.String())

		s := &Server{}

		_ = s.UpdateUser(ctx)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when error to get user data", func(t *testing.T) {
		e := echo.New()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reqParam := testRequestEndpointParam{
			e:          e,
			httpMethod: http.MethodPut,
			token:      token,
			url:        "/users",
			body:       []byte(successReqBody),
		}

		ctx, rec := TestRequestEndpoint(reqParam)
		ctx.Set("UserGUID", mockUser.GUID.String())

		mockRepo := repository.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().GetUserByGUID(gomock.Any(), mockUser.GUID).Return(nil, fmt.Errorf("error db"))
		s := &Server{
			Repository: mockRepo,
			Config:     *mockConfig,
		}
		_ = s.UpdateUser(ctx)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when error due to user is not found", func(t *testing.T) {
		e := echo.New()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reqParam := testRequestEndpointParam{
			e:          e,
			httpMethod: http.MethodPut,
			token:      token,
			url:        "/users",
			body:       []byte(successReqBody),
		}

		ctx, rec := TestRequestEndpoint(reqParam)
		ctx.Set("UserGUID", mockUser.GUID.String())

		mockRepo := repository.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().GetUserByGUID(gomock.Any(), mockUser.GUID).Return(nil, sql.ErrNoRows)
		s := &Server{
			Repository: mockRepo,
			Config:     *mockConfig,
		}
		_ = s.UpdateUser(ctx)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
	t.Run("when error to update user", func(t *testing.T) {
		e := echo.New()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reqParam := testRequestEndpointParam{
			e:          e,
			httpMethod: http.MethodPut,
			token:      token,
			url:        "/users",
			body:       []byte(successReqBody),
		}

		ctx, rec := TestRequestEndpoint(reqParam)
		ctx.Set("UserGUID", mockUser.GUID.String())

		mockRepo := repository.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().GetUserByGUID(gomock.Any(), mockUser.GUID).Return(mockUser, nil)

		mockUser.FullName = "test update user"
		mockUser.PhoneNumber = "+62345678901"

		mockRepo.EXPECT().UpdateUser(gomock.Any(), mockUser).Return(fmt.Errorf("error db"))

		s := &Server{
			Repository: mockRepo,
			Config:     *mockConfig,
		}
		_ = s.UpdateUser(ctx)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when error to update due to phone number is duplicate", func(t *testing.T) {
		e := echo.New()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reqParam := testRequestEndpointParam{
			e:          e,
			httpMethod: http.MethodPut,
			token:      token,
			url:        "/users",
			body:       []byte(successReqBody),
		}

		ctx, rec := TestRequestEndpoint(reqParam)
		ctx.Set("UserGUID", mockUser.GUID.String())

		mockRepo := repository.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().GetUserByGUID(gomock.Any(), mockUser.GUID).Return(mockUser, nil)

		mockUser.FullName = "test update user"
		mockUser.PhoneNumber = "+62345678901"

		mockRepo.EXPECT().UpdateUser(gomock.Any(), mockUser).Return(&tools.Err{Code: http.StatusConflict, Message: "phone number is already exists"})

		s := &Server{
			Repository: mockRepo,
			Config:     *mockConfig,
		}
		_ = s.UpdateUser(ctx)

		assert.Equal(t, http.StatusConflict, rec.Code)
	})
}
