package server

import "time"

// TableSettings is a struct that holds the settings of the whole DB
// from-kt: I suppose, better name would be StorageSettings, but I wanted to keep the original name
type TableSettings struct {
	StorageIdentityKey string    `json:"storageIdentityKey"`
	StorageName        string    `json:"storageName"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	Chain              string    `json:"chain"`
	DbType             string    `json:"dbtype"`
	MaxOutputScript    int       `json:"maxOutputScript"`
}
