package handler

import (
	"net/http"
	"strings"

	"github.com/SawitProRecruitment/UserService/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(cfg config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing JWT token")
			}

			tokenString := strings.Split(authHeader, " ")[1]

			key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cfg.RSAPublicKey))
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid JWT token")
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return key, nil
			})

			if err != nil {
				return echo.NewHTTPError(http.StatusForbidden, "Invalid JWT token")
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid JWT claims")
			}

			c.Set("UserGUID", claims["user_guid"])
			c.Set("FullName", claims["full_name"])
			return next(c)
		}
	}
}
