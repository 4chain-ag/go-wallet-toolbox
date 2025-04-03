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

func (u *UTXOs) FindAllUTXOs(ctx context.Context, userID int, page *paging.Page) ([]*models.UserUTXO, error) {
	var result []*models.UserUTXO
	err := u.db.WithContext(ctx).
		Scopes(scopes.UserID(userID), scopes.Paginate(page)).
		Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func NewUTXOs(db *gorm.DB) *UTXOs {
	return &UTXOs{
		db: db,
	}
}
