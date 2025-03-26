package wdk

import (
	"time"
)

// TableUser is a struct that defines the user from the DB
type TableUser struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    int       `json:"userId"`
	// IdentityKey is a pubKeyHex uniquely identifying user.
	// Typically, 66 hex digits.
	IdentityKey string `json:"identityKey"`
	// ActiveStorage is the storageIdentityKey value of the active wallet storage.
	ActiveStorage string `json:"activeStorage"`
}

// FindOrInsertUserResponse is a struct that holds information about the user and if it's new
type FindOrInsertUserResponse struct {
	User  TableUser `json:"user"`
	IsNew bool      `json:"isNew"`
}
