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

func (u *Users) FindOrCreateUser(identityKey string) (*wdk.TableUser, error) {
	user := &models.User{}
	err := u.db.First(&user, "identity_key = ?", identityKey).Error
	if err == nil {
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
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to find or create user: %w", err)
	}

	// by the time here we haven't found a user in a database and we need to create a new one
	user.IdentityKey = identityKey
	err = u.createUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Add default change basket for new user.
	outputBasket := &wdk.TableOutputBasket{
		UserID:                  user.UserID,
		Name:                    "default",
		NumberOfDesiredUTXOs:    32,
		MinimumDesiredUTXOValue: 1000,
		IsDeleted:               false,
	}

	err = u.outputBaskets.Create(outputBasket)
	if err != nil {
		return nil, fmt.Errorf("failed to create output basket for new user: %w", err)
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

func (u *Users) createUser(user *models.User) error {
	settings, err := u.settings.ReadSettings()
	if err != nil {
		return fmt.Errorf("failed to read settings: %w", err)
	}

	user.ActiveStorage = settings.StorageIdentityKey

	return u.db.Create(user).Error
}
