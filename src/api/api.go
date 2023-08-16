package api

import (
	"fmt"
	"net/http"

	"crypto/rand"
	"encoding/base64"

	"libreriaGO/src/db"

	"github.com/labstack/echo/v4"
)

func SendUserData(c echo.Context) error {
	username := c.Param("username")
	user, err := db.UserData(username)
	if err != nil {
		fmt.Println("Error en SendUserData")
	}
	fmt.Println(user)
	response := echo.Map{
		"message": "User data retrieved successfully",
		"user":    user, // Assuming user is a struct containing user data
	}

	return c.JSON(http.StatusOK, response)
}

func generateSecretKey(keyLength int) (string, error) {
	key := make([]byte, keyLength)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(key), nil
}

var SecretKey string

func init() {
	key, err := generateSecretKey(32)
	if err != nil {
		panic("Error generating secret key")
	}
	SecretKey = key
}
