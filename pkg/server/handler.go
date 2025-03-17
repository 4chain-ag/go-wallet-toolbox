package server

import (
	"fmt"
	"time"
)

type Handler struct {
}

func (h *Handler) MakeAvailable() TableSettings {
	fmt.Println("MakeAvailable")
	return TableSettings{
		StorageIdentityKey: "028f2daab7808b79368d99eef1ebc2d35cdafe3932cafe3d83cf17837af034ec29",
		StorageName:        "test-go-jsonrpc",
		CreatedAt:          time.Now().Add(-time.Hour),
		UpdatedAt:          time.Now().Add(-time.Minute),
		Chain:              "test",
		DbType:             "MySQL",
		MaxOutputScript:    100 * 1024,
	}
}
