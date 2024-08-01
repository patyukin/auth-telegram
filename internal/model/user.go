package model

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type UserRole string

const (
	UserRoleDefault UserRole = "user"
	UserRoleAdmin   UserRole = "admin"
)

type SignUpData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Telegram string `json:"telegram"`
}

type SignInData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SignInVerifyData struct {
	Code string `json:"code"`
}

type User struct {
	UUID         uuid.UUID      `json:"id,omitempty"`
	Login        string         `json:"login"`
	PasswordHash string         `json:"password_hash"`
	Telegram     string         `json:"telegram"`
	Name         sql.NullString `json:"name"`
	Surname      sql.NullString `json:"surname"`
	Role         UserRole       `json:"role"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    sql.NullTime   `json:"updated_at"`
}

type Token struct {
	Token     uuid.UUID    `json:"token"`
	UserUUID  uuid.UUID    `json:"user_id"`
	ExpiresAt time.Time    `json:"expires_at"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}
