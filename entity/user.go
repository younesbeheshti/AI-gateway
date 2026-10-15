package entity

import "time"

// You can add other fields like: password_hashed, role, status, is_email_verified. but I keep it simple
type User struct {
	ID           uint64
	Email        string
	ApiKey       string
	TokenBalance int64
	CreatedAt    time.Time
}
