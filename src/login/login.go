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

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get(echo.HeaderAuthorization)
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
		}

		token, err := jwt.ParseWithClaims(tokenString, &db.User{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
		}

		// Attach the user information to the request context
		user := token.Claims.(*db.User)
		c.Set("user", user)

		return next(c)
	}
}

func ProtectedRoute(c echo.Context) error {
	user := c.Get("user").(*db.User) // Get user info from the context
	// Your protected route logic here
	return c.JSON(http.StatusOK, echo.Map{"message": "Welcome, " + user.Name})
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
		fmt.Println("User not found 10", requestBody.Name, requestBody.Pass)
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid credentials"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &db.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		// Populate other fields as needed...
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	})

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