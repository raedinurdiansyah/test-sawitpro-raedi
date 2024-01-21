package handler

import (
	"bytes"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
)

type testRequestEndpointParam struct {
	e          *echo.Echo
	httpMethod string
	token      string
	url        string
	body       []byte
}

func TestRequestEndpoint(param testRequestEndpointParam) (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(param.httpMethod, param.url, bytes.NewBuffer(param.body))
	if param.token != "" {
		req.Header.Add(echo.HeaderAuthorization, "Bearer "+param.token)
	}
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	return param.e.NewContext(req, rec), rec
}
