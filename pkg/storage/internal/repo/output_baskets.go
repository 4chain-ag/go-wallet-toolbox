package repo

import (
	"fmt"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/database/models"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"gorm.io/gorm"
)

type OutputBaskets struct {
	db *gorm.DB
}

func NewOutputBaskets(db *gorm.DB) *OutputBaskets {
	return &OutputBaskets{db: db}
}

func (u *OutputBaskets) Create(outputBasket *wdk.TableOutputBasket) error {
	outputBasketModel := &models.OutputBaskets{
		BasketID:                outputBasket.BasketID,
		UserID:                  outputBasket.UserID,
		MinimumDesiredUTXOValue: outputBasket.MinimumDesiredUTXOValue,
		NumberOfDesiredUTXOs:    outputBasket.NumberOfDesiredUTXOs,
		Name:                    outputBasket.Name,
	}

	err := u.db.Create(outputBasketModel).Error
	if err != nil {
		return fmt.Errorf("failed to create output basket model: %w", err)
	}
	return nil
}
