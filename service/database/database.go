package database

import (
	"database/sql"
	"errors"
	"fmt"
	"image/gif"
	"time"
)

type AppDatabase interface {
	Ping() error

	CheckIfUserExists(username string) (bool, error)

	AddNewUser(username string, securityKey string) (int, error)

	GetUserName(userID int) (string, error)

	GetUserKey(userID int) (string, error)

	GetUserID(username string) (int, error)
}

type User struct {
	ID          int
	Username    string
	SecurityKey string
	GifPhoto    *gif.GIF
}

type Message struct {
	ID        int
	Content   string
	GifPhoto  *gif.GIF
	SenderID  int
	Checkmark string
	Timestamp time.Time
}

type Conversation struct {
	ID       int
	Name     string
	GifPhoto *gif.GIF
	Members  []int
}

type appdbimpl struct {
	c *sql.DB
}

func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building an AppDatabase")
	}

	usersTableStmt := `CREATE TABLE IF NOT EXISTS Users (
				id INTEGER NOT NULL PRIMARY KEY,
				username TEXT NOT NULL,
				security_key TEXT NOT NULL,
				gif_photo BLOB
				);`
	if _, err := db.Exec(usersTableStmt); err != nil {
		return nil, fmt.Errorf("error creating Users table: %w", err)
	}

	conversationsTableStmt := `CREATE TABLE IF NOT EXISTS Conversations (
					id INTEGER NOT NULL PRIMARY KEY,
					name TEXT NOT NULL,
					gif_photo BLOB
					);`
	if _, err := db.Exec(conversationsTableStmt); err != nil {
		return nil, fmt.Errorf("error creating Conversations table: %w", err)
	}

	conversationMembersStmt := `CREATE TABLE IF NOT EXISTS ConversationMembers (
					 conv_id INTEGER NOT NULL,
					 user_id INTEGER NOT NULL,
					 PRIMARY KEY (conv_id, user_id),
					 FOREIGN KEY (conv_id) REFERENCES Conversations(id),
					 FOREIGN KEY (user_id) REFERENCES Users(id)
					 );`
	if _, err := db.Exec(conversationMembersStmt); err != nil {
		return nil, fmt.Errorf("error creating ConversationMembers table: %w", err)
	}

	messagesTableStmt := `CREATE TABLE IF NOT EXISTS Messages (
				   id INTEGER NOT NULL PRIMARY KEY,
				   conv_id INTEGER NOT NULL,
				   content TEXT,
				   gif_photo BLOB,
				   sender_id INTEGER NOT NULL,
				   checkmark TEXT NOT NULL,
				   timestamp DATETIME NOT NULL,
				   FOREIGN KEY (conv_id) REFERENCES Conversations(id),
				   FOREIGN KEY (sender_id) REFERENCES Users(id)
				   );`
	if _, err := db.Exec(messagesTableStmt); err != nil {
		return nil, fmt.Errorf("error creating Messages table: %w", err)
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
