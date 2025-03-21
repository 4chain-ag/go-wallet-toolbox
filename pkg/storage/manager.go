package storage

import (
	"fmt"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"sync"
)

type managedStorage struct {
	settings *wdk.SettingsDTO
	user     *wdk.UserDTO
	provider wdk.WalletStorageProvider
}

func (ms *managedStorage) isAvailable() bool {
	return ms != nil && ms.settings != nil && ms.user != nil
}

type Manager struct {
	stores []*managedStorage
	authID wdk.AuthID

	active             *managedStorage
	backups            []*managedStorage
	conflictingActives []*managedStorage
	lockActive         sync.Locker
}

func NewManager(identityKey string, active wdk.WalletStorageProvider, backups ...wdk.WalletStorageProvider) *Manager {
	stores := make([]*managedStorage, len(backups)+1)
	for i, provider := range append([]wdk.WalletStorageProvider{active}, backups...) {
		stores[i] = &managedStorage{
			provider: provider,
		}
	}

	return &Manager{
		stores: stores,
		authID: wdk.AuthID{IdentityKey: identityKey},
	}
}

func (m *Manager) isStorageEnabled(managed *managedStorage) bool {
	return managed.isAvailable() && managed.settings.StorageIdentityKey == managed.user.ActiveStorage
}

func (m *Manager) IsAvailable() bool {
	return m.active.isAvailable()
}

func (m *Manager) MakeAvailable() (*wdk.SettingsDTO, error) {
	if m.active.isAvailable() {
		return m.active.settings, nil
	}

	if len(m.stores) == 0 {
		return nil, fmt.Errorf("no storage providers available")
	}

	m.lockActive.Lock()
	defer m.lockActive.Unlock()

	m.active = nil
	m.backups = make([]*managedStorage, 0)
	m.conflictingActives = make([]*managedStorage, 0)

	var err error
	var backups []*managedStorage
	for _, store := range m.stores {
		if !store.isAvailable() {
			store.settings, err = store.provider.MakeAvailable()
			if err != nil {
				return nil, fmt.Errorf("failed to make storage provider available: %w", err)
			}
			// TODO: Handle user/findOrInsertUser
			store.user = &wdk.UserDTO{} // fixme: for now, just a stub
		}

		if m.active == nil {
			// stores[0] becomes the default active store.
			// It may be replaced if it is not the user's "enabled" activeStorage and that store is found among the remainder (backups).
			m.active = store
		} else {
			userActive := store.user.ActiveStorage
			storageIdentity := store.settings.StorageIdentityKey

			if userActive == storageIdentity && m.isStorageEnabled(store) {
				backups = append(backups, m.active)
				m.active = store
			} else {
				backups = append(backups, store)
			}
		}
	}

	// review backups, partition out conflicting actives.
	activeIdentityKey := m.active.settings.StorageIdentityKey
	for _, backupStore := range backups {
		if backupStore.user.ActiveStorage != activeIdentityKey {
			m.conflictingActives = append(m.conflictingActives, backupStore)
		} else {
			m.backups = append(m.backups, backupStore)
		}
	}

	userID := m.active.user.ID
	m.authID.UserID = &userID

	isActiveEnabled := m.isStorageEnabled(m.active) && len(m.conflictingActives) == 0
	m.authID.IsActive = &isActiveEnabled

	return m.active.settings, nil
}
