package repo

import (
	"errors"
	"fmt"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/database/models"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
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
		User: wdk.User{
			UserID:        user.UserID,
			IdentityKey:   user.IdentityKey,
			ActiveStorage: user.ActiveStorage,
			CreatedAt:     user.CreatedAt,
			UpdatedAt:     user.UpdatedAt,
		},
		IsNew: false,
	}, nil
}

func (u *Users) CreateUser(identityKey string) (*wdk.TableUser, error) {
	user := &models.User{
		OutputBaskets: []*models.OutputBaskets{{
			Name:                    "default",
			NumberOfDesiredUTXOs:    32,
			MinimumDesiredUTXOValue: 1000,
		}},
	}
	user.IdentityKey = identityKey

	settings, err := u.settings.ReadSettings()
	if err != nil {
		return nil, fmt.Errorf("failed to read settings: %w", err)
	}

	user.ActiveStorage = settings.StorageIdentityKey

	err = u.db.Create(user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &wdk.TableUser{
		User: wdk.User{
			UserID:        user.UserID,
			IdentityKey:   user.IdentityKey,
			ActiveStorage: user.ActiveStorage,
			CreatedAt:     user.CreatedAt,
			UpdatedAt:     user.UpdatedAt,
		},
		IsNew: true,
	}, nil
}
