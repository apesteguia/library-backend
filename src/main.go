package main

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

import "libreriaGO/src/login"

const PORT = ":8080"

type jwtCustomClaims struct {
	Name string `json:"name"`
	Pass bool   `json:"pass"`
	jwt.RegisteredClaims
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	name := claims.Name
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization},
	}))

	e.POST("/login", login.Login)
	e.POST("/register", login.Register)

	r := e.Group("/profile")
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte(login.SecretKey),
	}
	r.Use(echojwt.WithConfig(config))
	r.GET("/:username", login.SendUserData)

	e.Logger.Fatal(e.Start(PORT))
}
