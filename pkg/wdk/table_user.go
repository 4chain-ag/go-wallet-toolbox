package wdk

import (
	"time"
)

// User is a struct that defines the user from the DB
type User struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    int       `json:"userId"`
	// IdentityKey is a pubKeyHex uniquely identifying user.
	// Typically, 66 hex digits.
	IdentityKey string `json:"identityKey"`
	// ActiveStorage is the storageIdentityKey value of the active wallet storage.
	ActiveStorage string `json:"activeStorage"`
}

// TableUser is a struct that holds information about the user and if it's new
type TableUser struct {
	User  User `json:"user"`
	IsNew bool `json:"isNew"`
}
