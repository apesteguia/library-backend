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
	Books              []Book `json:"books"`
	jwt.StandardClaims
}

type Book struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	ImageURL      string `json:"image_url"`
	NumberOfPages int    `json:"number_of_pages"`
	Description   string `json:"description"`
}

const DB = "database.db"

func UserData(username string) (User, error) {
	var user User
	db, err := sql.Open("sqlite3", DB)
	if err != nil {
		return user, err
	}
	defer db.Close()

	query := "SELECT id, email, password, profile_description, profile_picture_url FROM users WHERE username = ?"

	row := db.QueryRow(query, username)
	err = row.Scan(&user.ID, &user.Email, &user.Pass, &user.ProfileDescription, &user.ProfilePictureURL)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errors.New("user not found")
		}
		return user, err
	}

	// Query books associated with the user
	booksQuery := "SELECT b.id, b.name, b.image_url, b.number_of_pages, b.description FROM books b JOIN user_books ub ON b.id = ub.book_id JOIN users u ON u.id = ub.user_id WHERE u.username = ?"
	rows, err := db.Query(booksQuery, username)
	if err != nil {
		return user, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Name, &book.ImageURL, &book.NumberOfPages, &book.Description)
		if err != nil {
			return user, err
		}
		books = append(books, book)
	}

	user.Books = books
	user.Name = username

	return user, nil
}

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
		return user, err
	}
	defer db.Close()

	query := "SELECT username, password FROM users WHERE username = ?"
	row := db.QueryRow(query, username)

	err = row.Scan(&user.Name, &user.Pass)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errors.New("user not found")
		}
		return user, err
	}

	return user, nil
}
