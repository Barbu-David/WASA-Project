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
/*
	GetUserIDBySecurityKey(securityKey string) (int, error)

	UpdateUserName(userID int, name string) error
	
	UpdateUserPhoto(userID int, photo *gif.GIF) error

	GetConversationsByUserID(userID int) ([]int, error)

	GetMessageIDsByConversationID(convID int) ([]int, error)

	AddMessageToConversation(convID int, message Message) (int, error)

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
*/
	CheckIfUserExists(username string) (bool, error)

	AddNewUser(username string, securityKey string) (int, error) 

	GetUserName(userID int) (string, error)
	
	GetUserKey(userID int) (string, error)

	GetUserID(username string) (int, error)
}

type User struct {
	ID           int
		Username     string
		SecurityKey  string
		GifPhoto     *gif.GIF
}

type Message struct {
	ID           int
		Content      string
		GifPhoto     *gif.GIF
		SenderID     int
		Checkmark    string 
		Timestamp    time.Time
}

type Conversation struct {
	ID           int
		Name         string
		GifPhoto     *gif.GIF
		Members      []int   
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

