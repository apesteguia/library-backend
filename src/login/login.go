package login

import (
	"fmt"
	"net/http"
	"time"

	"crypto/rand"
	"encoding/base64"

	"libreriaGO/src/db"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

func SendJson(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Hello, World!",
	})
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

func Login(c echo.Context) error {
	var requestBody struct {
		Name string `json:"name"`
		Pass string `json:"pass"`
	}
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
	}

	user, err := db.FindUser(requestBody.Name, requestBody.Pass)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid credentials"})
	}

	if user.Name != requestBody.Name || user.Pass != requestBody.Pass {
		return nil
	}

	// Create JWT token using the HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   user.Name, // Set the subject (user name) as the token's subject
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign the token using the SecretKey
	signedToken, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error creating token")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": signedToken,
	})
}

func Register(c echo.Context) error {
	var requestBody struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Pass  string `json:"pass"`
	}
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
	}

	_, err := db.FindUser(requestBody.Name, requestBody.Pass)
	if err == nil {
		return c.JSON(http.StatusConflict, echo.Map{"error": "User already exists"})
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &db.User{
		Name: requestBody.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	})

	signedToken, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error creating token")
	}

	// Add user to the database
	fmt.Println("HOLAJK", requestBody.Name)
	err = db.AddUser(requestBody.Name, requestBody.Email, requestBody.Pass)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error registering user"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "User registered successfully",
		"token":   signedToken,
	})
}
