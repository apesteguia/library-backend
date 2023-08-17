-- Create tables
CREATE TABLE friends (
    user_id INTEGER,
    friend_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (friend_id) REFERENCES users(id)
);

CREATE TABLE user_owned_books (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    book_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (book_id) REFERENCES books(id)
);

CREATE TABLE user_books (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    book_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (book_id) REFERENCES books(id)
);

CREATE TABLE books (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    image_url TEXT,
    number_of_pages INTEGER,
    description TEXT
);

CREATE TABLE sqlite_sequence (name, seq);

CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT,
    password TEXT,
    jwt_token TEXT,
    profile_description TEXT,
    profile_picture_url TEXT
);

-- Insert data into users table
INSERT INTO users (username, email, password, jwt_token, profile_description, profile_picture_url)
VALUES ('mikel', 'lizarrakiller@gmail.com', 'mikel000', NULL, 'This is a random description for user mikel.', NULL);

-- You can add more INSERT statements for other data as needed

