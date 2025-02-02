package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/vipinvkartha/lenslocked/rand"
)

const (
	// The minimum number of bytes to be used for each session token.
	MinBytesPerToken = 32
)

type Session struct {
	ID     int
	UserID int
	// Token is only set when creating a new session. When looking up a session
	// this will be left empty, as we only store the hash of a session token
	// in our database and we cannot reverse it into a raw token.
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
	// BytesPerToken is used to determine how many bytes to use when generating
	// each session token. If this value is not set or is less than the
	// MinBytesPerToken const it will be ignored and MinBytesPerToken will be
	// used.
	BytesPerToken int
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	// TODO: Hash the session token
	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: ss.hash(token),
	}
	row := ss.DB.QueryRow(`
		UPDATE sessions
		SET token_hash = $2
		WHERE user_id = $1
    RETURNING id;`, session.UserID, session.TokenHash)
	err = row.Scan(&session.ID)
	if err == sql.ErrNoRows {
		// If no session exists, we will get ErrNoRows. That means we need to
		// create a session object for that user.
		row = ss.DB.QueryRow(`
			INSERT INTO sessions (user_id, token_hash)
			VALUES ($1, $2)
			RETURNING id;`, session.UserID, session.TokenHash)
		// The error will be overwritten with either a new error, or nil
		err = row.Scan(&session.ID)
	}
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	// TODO: Store the session in our DB
	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	// 1. Hash the session token
	tokenHash := ss.hash(token)
	// 2. Query for the session with that hash
	var user User
	row := ss.DB.QueryRow(`
	SELECT user_id
	FROM sessions
	WHERE token_hash = $1;`, tokenHash)
	err := row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}
	// 3. Using the UserID from the session, we need to query for that user

	row = ss.DB.QueryRow(`
	SELECT email, password_hash
	FROM users WHERE id = $1;`, user.ID)
	err = row.Scan(&user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}
	// 4. Return the user
	return &user, nil
}

func (ss *SessionService) Delete(token string) error {
	tokenHash := ss.hash(token)
	_, err := ss.DB.Exec(`
		DELETE FROM sessions
		WHERE token_hash = $1;`, tokenHash)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	// base64 encode the data into a string
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
