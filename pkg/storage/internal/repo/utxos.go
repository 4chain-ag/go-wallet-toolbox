package repo

import (
	"context"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/database/models"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/database/scopes"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/paging"
	"gorm.io/gorm"
)

type UTXOs struct {
	db *gorm.DB
}

func NewUTXOs(db *gorm.DB) *UTXOs {
	return &UTXOs{
		db: db,
	}
}

func (u *UTXOs) FindAllUTXOs(ctx context.Context, userID int, basketID int, page *paging.Page) ([]*models.UserUTXO, error) {
	var result []*models.UserUTXO
	err := u.db.WithContext(ctx).
		Scopes(scopes.UserID(userID), scopes.BasketID(basketID), scopes.Paginate(page)).
		Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *UTXOs) CountUTXOs(ctx context.Context, userID int, basket int) (int64, error) {
	count := int64(0)

	err := u.db.WithContext(ctx).
		Model(&models.UserUTXO{}).
		Scopes(scopes.UserID(userID), scopes.BasketID(basket)).
		Count(&count).Error

	return count, err
}
