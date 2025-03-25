package wdk

import (
	"time"
)

// TableOutputBasket is a struct that holds the output baskets details
type TableOutputBasket struct {
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
	BasketID                int       `json:"basketId"`
	UserID                  int       `json:"userId"`
	Name                    string    `json:"name"`
	NumberOfDesiredUTXOs    int       `json:"numberOfDesiredUTXOs"`
	MinimumDesiredUTXOValue int       `json:"minimumDesiredUTXOValue"`
	IsDeleted               bool      `json:"isDeleted"`
}
