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
	SetUserName(userID int, username string) error

	GetUserPhoto(userID int) (*gif.GIF, error)
	SetUserPhoto(userID int, photo *gif.GIF) error

	GetUserKey(userID int) (string, error)
	GetUserID(username string) (int, error)
	GetUserIDbyKey(security_key string) (int, error)

	GetMaxUserID() (int, error)

	NewConversation(name string, group bool) (int, error)
	NewConversationMember(user_id int, conv_id int) error
	IsMemberConversation(user_id int, conv_id int) (bool, error)

	DeleteConversationMember(user_id int, conv_id int) error

	GetConversationName(conv_id int) (string, error)
	SetConversationName(conv_id int, name string) error

	GetConversationPhoto(conv_id int) (*gif.GIF, error)
	SetConversationPhoto(conv_id int, photo *gif.GIF) error

	IsGroupConversation(conv_id int) (bool, error)

	GetUserConversations(userID int) ([]int, error)
	GetConversationUsers(conv_id int) ([]int, error)
	GetConversationMessages(conv_id int) ([]int, error)

	SendMessage(senderID int, convID int, textContent string, forwarded bool, timestamp time.Time) error
	GetMessage(m_id int) (int, string, bool, time.Time, error)
	DeleteMessage(m_id int) error

	SeeMessage(user_id int, m_id int) error
	ReceiveMessage(user_id int, m_id int) error

	GetMessageCommentList(m_id int) ([]int, []string, error)
	GetMessageLatest(conv_id int) (string, time.Time, bool, error)
	
	AddComment(sender_id int, m_id int, content string) error
	RemoveComment(sender_id int, m_id int) error

	IsSeenByAll(m_id int) (bool, error)
	IsDeliveredToAll(m_id int) (bool, error)

	// GetMessagePhoto(m_id int) (gif.GIF, error)
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
					is_group BOOL,
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
				   timestamp DATETIME NOT NULL,
				   forwarded BOOL,
				   is_photo BOOL, 
				   FOREIGN KEY (conv_id) REFERENCES Conversations(id),
				   FOREIGN KEY (sender_id) REFERENCES Users(id)
				   );`
	if _, err := db.Exec(messagesTableStmt); err != nil {
		return nil, fmt.Errorf("error creating Messages table: %w", err)
	}

	seenListStmt := `CREATE TABLE IF NOT EXISTS SeenList (
					 m_id INTEGER NOT NULL,
					 user_id INTEGER NOT NULL,
					 comment TEXT,
					 seen BOOL,
					 received BOOL,
					 id INTEGER NOT NULL PRIMARY KEY,
					 FOREIGN KEY (m_id) REFERENCES Messages(id),
					 FOREIGN KEY (user_id) REFERENCES Users(id)
					 );`
	if _, err := db.Exec(seenListStmt); err != nil {
		return nil, fmt.Errorf("error creating SeenList table: %w", err)
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
