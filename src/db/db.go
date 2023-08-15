package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Email              string `json:"email"`
	Pass               string `json:"pass"`
	JWTToken           string `json:"jwt"`
	ProfileDescription string `json:"desc"`
	ProfilePictureURL  string `json:"pic"`
	jwt.StandardClaims
}

const DB = "database.db"

func AddUser(username string, email string, pass string) error {
	db, err := sql.Open("sqlite3", DB)
	if err != nil {
		return err
	}
	defer db.Close()

	boo := UserExistsByUsername(username)
	if boo == true {
		fmt.Println("User already exists:", username)
		return errors.New("user already exists")
	}

	insertSQL := `
        INSERT INTO users (username, email, password)
        VALUES (?, ?, ?)
    `
	_, err = db.Exec(insertSQL, username, email, pass)
	if err != nil {
		fmt.Println("Error in exec:", err, username)
		return err
	}

	return nil
}

func UserExistsByUsername(username string) bool {
	db, err := sql.Open("sqlite3", DB)
	if err != nil {
		return false
	}
	defer db.Close()

	query := "SELECT id FROM users WHERE username = ?"
	row := db.QueryRow(query, username)

	var userID int
	err = row.Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false // User does not exist
		}
		return false
	}

	return true
}
func FindUser(username string, pass string) (User, error) {
	var user User

	db, err := sql.Open("sqlite3", DB)
	if err != nil {
		fmt.Println("No abre")
		return user, err
	}
	defer db.Close()

	query := "SELECT username, password FROM users WHERE username = ?"
	row := db.QueryRow(query, username)

	var storedPass string
	err = row.Scan(&user.Name, &storedPass)
	if err != nil {
		fmt.Println("User not found 1")
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("User not found 2")
			return user, errors.New("user not found")
		}
		return user, err
	}

	if err != nil {
		fmt.Println("Invalid password")
		return user, errors.New("invalid password")
	}

	return user, nil
}
