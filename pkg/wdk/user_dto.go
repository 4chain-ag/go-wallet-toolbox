package wdk

import "time"

type UserDTO struct {
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	ID            int       `json:"userId"`
	IdentityKey   string    `json:"identityKey"`
	ActiveStorage string    `json:"activeStorage"`
}
