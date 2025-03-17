package server

import "time"

type TableSettings struct {
	StorageIdentityKey string    `json:"storageIdentityKey"`
	StorageName        string    `json:"storageName"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	Chain              string    `json:"chain"`
	DbType             string    `json:"dbtype"`
	MaxOutputScript    int       `json:"maxOutputScript"`
}
