package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"image/gif" // Importing for handling GIFs
)

type AppDatabase interface {
	Ping() error

	GetUserIDBySecurityKey(securityKey string) (int, error)

	UpdateUserPhoto(userID int, photo *gif.GIF) error

	UpdateUserName(userID int, name string) error

	GetConversationsByUserID(userID int) ([]int, error)

	GetMessageIDsByConversationID(convID int) ([]int, error)

	AddMessageToConversation(convID int, message Message) (int, error) // Returns the new message ID

	DeleteMessageFromConversation(convID int, messageID int) error

	AddCommentToMessage(convID int, messageID int, comment string) error

	RemoveCommentFromMessage(convID int, messageID int) error

	SetConversationName(convID int, userID int, name string) error

	SetConversationPhoto(convID int, userID int, photo *gif.GIF) error

	GetMessageByID(convID int, messageID int) (Message, error)

	GetConversationNameByUserID(userID int, convID int) (string, error)

	GetConversationPhotoByUserID(userID int, convID int) (*gif.GIF, error)

	AddUserToConversation(convID int, userID int) error
		
	RemoveUserFromConversation(convID int, userID int) error
}

type User struct {
	ID           int
	Username     string
	SecurityKey  string
	GifPhoto     *gif.GIF
}

type Message struct {
	Content      string
	GifPhoto     *gif.GIF
	SenderID     int
	Checkmark    string // Possible values: "seen", "delivered", "unseen"
	Timestamp    time.Time
}

type Conversation struct {
	ID           int
	Name         string
	GifPhoto     *gif.GIF
	Members      []int       // List of user IDs
	Messages     []Message
}

type appdbimpl struct {
	c *sql.DB
}

func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building an AppDatabase")
	}

	// Create Users table
	usersTableStmt := `CREATE TABLE IF NOT EXISTS Users (
		id INTEGER NOT NULL PRIMARY KEY,
		username TEXT NOT NULL,
		security_key TEXT NOT NULL,
		gif_photo TEXT
	);`
	if _, err := db.Exec(usersTableStmt); err != nil {
		return nil, fmt.Errorf("error creating Users table: %w", err)
	}

	// Create Conversations table
	conversationsTableStmt := `CREATE TABLE IF NOT EXISTS Conversations (
		id INTEGER NOT NULL PRIMARY KEY,
		name TEXT NOT NULL,
		gif_photo TEXT,
		member_list TEXT NOT NULL,
		message_list TEXT
	);`
	if _, err := db.Exec(conversationsTableStmt); err != nil {
		return nil, fmt.Errorf("error creating Conversations table: %w", err)
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

