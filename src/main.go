package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

import "libreriaGO/src/login"

const PORT = ":8080"

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization},
	}))

	e.POST("/login", login.Login)
	e.POST("/register", login.Register)
	e.GET("/users", login.ProtectedRoute, login.AuthMiddleware)

	e.Logger.Fatal(e.Start(PORT))
}
