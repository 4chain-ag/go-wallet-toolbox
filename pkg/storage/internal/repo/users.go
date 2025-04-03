package repo

import (
	"errors"
	"fmt"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/database/models"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/go-softwarelab/common/pkg/seq"
	"gorm.io/gorm"
)

type Users struct {
	db            *gorm.DB
	settings      *Settings
	outputBaskets *OutputBaskets
}

func NewUsers(db *gorm.DB, settings *Settings, outputBaskets *OutputBaskets) *Users {
	return &Users{db: db, settings: settings, outputBaskets: outputBaskets}
}

func (u *Users) FindUser(identityKey string) (*wdk.TableUser, error) {
	user := &models.User{}
	err := u.db.First(&user, "identity_key = ?", identityKey).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find or create user: %w", err)
	}

	return &wdk.TableUser{
		UserID:        user.UserID,
		IdentityKey:   user.IdentityKey,
		ActiveStorage: user.ActiveStorage,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}, nil
}

func (u *Users) CreateUser(identityKey, activeStorage string, baskets ...wdk.BasketConfiguration) (*wdk.TableUser, error) {
	user := models.User{
		IdentityKey:   identityKey,
		ActiveStorage: activeStorage,
		OutputBaskets: seq.Collect(seq.Map(seq.Of(baskets...), func(basket wdk.BasketConfiguration) *models.OutputBasket {
			return &models.OutputBasket{
				Name:                    basket.Name,
				NumberOfDesiredUTXOs:    basket.NumberOfDesiredUTXOs,
				MinimumDesiredUTXOValue: basket.MinimumDesiredUTXOValue,
			}
		})),
	}
	err := u.db.Create(&user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &wdk.TableUser{
		UserID:        user.UserID,
		IdentityKey:   user.IdentityKey,
		ActiveStorage: user.ActiveStorage,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}, nil
}
