package storage

import (
	"fmt"
	"log/slog"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/lox"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/utils/to"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/validate"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/actions"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/database"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/database/models"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/repo"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/samber/lo"
)

// Repository is an interface for the actual storage repository.
type Repository interface {
	Migrate() error

	ReadSettings() (*wdk.TableSettings, error)
	SaveSettings(settings *wdk.TableSettings) error

	FindUser(identityKey string) (*wdk.TableUser, error)
	CreateUser(user *models.User) (*wdk.TableUser, error)

	CreateCertificate(certificate *models.Certificate) (uint, error)
	DeleteCertificate(userID int, args wdk.RelinquishCertificateArgs) error
	ListAndCountCertificates(userID int, opts repo.ListCertificatesActionParams) ([]*models.Certificate, int64, error)
}

// Provider is a storage provider.
type Provider struct {
	Chain defs.BSVNetwork

	settings *wdk.TableSettings
	repo     Repository
	actions  *actions.Actions
}

// GORMProviderConfig is a configuration for GORM storage provider.
type GORMProviderConfig struct {
	DB       defs.Database
	Chain    defs.BSVNetwork
	FeeModel defs.FeeModel
}

// NewGORMProvider creates a new storage provider with GORM repository.
func NewGORMProvider(logger *slog.Logger, config GORMProviderConfig, opts ...ProviderOption) (*Provider, error) {
	if err := config.FeeModel.Validate(); err != nil {
		return nil, fmt.Errorf("invalid fee model: %w", err)
	}

	options := toOptions(opts)

	db, err := configureDatabase(logger, config.DB, options)
	if err != nil {
		return nil, err
	}

	return &Provider{
		Chain:   config.Chain,
		repo:    db.CreateRepositories(),
		actions: actions.New(logger, db.CreateFunder(config.FeeModel)),
	}, nil
}

func configureDatabase(logger *slog.Logger, dbConfig defs.Database, options *providerOptions) (*database.Database, error) {
	if options.gormDB != nil {
		return database.NewWithGorm(options.gormDB, logger), nil
	}

	db, err := database.NewDatabase(dbConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}
	return db, nil
}

// Migrate migrates the storage and saves the settings.
func (p *Provider) Migrate(storageName, storageIdentityKey string) (string, error) {
	err := p.repo.Migrate()
	if err != nil {
		return "", fmt.Errorf("failed to migrate: %w", err)
	}

	// TODO: what if p.Chain != Chain from DB?

	err = p.repo.SaveSettings(&wdk.TableSettings{
		StorageIdentityKey: storageIdentityKey,
		StorageName:        storageName,
		Chain:              p.Chain,
		MaxOutputScript:    DefaultMaxScriptLength,
	})
	if err != nil {
		return "", fmt.Errorf("failed to save settings: %w", err)
	}

	// NOTE: GORM automigrate does not support db versioning
	// from-kt: In TS version I can't find any usage of returned version
	version := "auto-migrated"

	return version, nil
}

// MakeAvailable reads the settings and makes them available.
func (p *Provider) MakeAvailable() (*wdk.TableSettings, error) {
	settings, err := p.repo.ReadSettings()
	if err != nil {
		return nil, fmt.Errorf("failed to read settings: %w", err)
	}

	p.settings = settings
	return settings, nil
}

// InsertCertificateAuth inserts certificate to the database for authenticated user
func (p *Provider) InsertCertificateAuth(auth wdk.AuthID, certificate *wdk.TableCertificateX) (uint, error) {
	if auth.UserID == nil || certificate.UserID != *auth.UserID {
		return 0, fmt.Errorf("access is denied due to an authorization error")
	}

	// TODO: validate arguments?

	certModel := &models.Certificate{
		Type:               string(certificate.Type),
		SerialNumber:       string(certificate.SerialNumber),
		Certifier:          string(certificate.Certifier),
		Subject:            string(certificate.Subject),
		RevocationOutpoint: string(certificate.RevocationOutpoint),
		Signature:          string(certificate.Signature),

		UserID:            *auth.UserID,
		CertificateFields: lo.Map(certificate.Fields, lox.MappingFn(tableCertificateXFieldsToModelFields(*auth.UserID))),
	}

	if certificate.Verifier != nil {
		certModel.Verifier = string(*certificate.Verifier)
	}

	id, err := p.repo.CreateCertificate(certModel)
	if err != nil {
		return 0, fmt.Errorf("failed to create certificate: %w", err)
	}

	return id, nil
}

// RelinquishCertificate will relinquish existing certificate
func (p *Provider) RelinquishCertificate(auth wdk.AuthID, args wdk.RelinquishCertificateArgs) error {
	if auth.UserID == nil {
		return fmt.Errorf("access is denied due to an authorization error")
	}
	// TODO: validate args

	err := p.repo.DeleteCertificate(*auth.UserID, args)
	if err != nil {
		return fmt.Errorf("failed to relinquish certificate: %w", err)
	}

	return nil
}

// ListCertificates will list certificates with provided args
func (p *Provider) ListCertificates(auth wdk.AuthID, args wdk.ListCertificatesArgs) (*wdk.ListCertificatesResult, error) {
	if auth.UserID == nil {
		return nil, fmt.Errorf("access is denied due to an authorization error")
	}
	// TODO: validate args

	// prepare arguments
	filterOptions := listCertificatesArgsToActionParams(args)

	certModels, totalCount, err := p.repo.ListAndCountCertificates(*auth.UserID, filterOptions)
	if err != nil {
		return nil, fmt.Errorf("error during listing certificates action: %w", err)
	}

	tc, err := to.UInt(totalCount)
	if err != nil {
		return nil, fmt.Errorf("error during parsing total count of certificates: %w", err)
	}

	result := &wdk.ListCertificatesResult{
		TotalCertificates: wdk.PositiveIntegerOrZero(tc),
		Certificates:      lo.Map(certModels, lox.MappingFn(certModelToResult)),
	}

	return result, nil
}

// FindOrInsertUser will find user by their identityKey or inserts a new one if not found
func (p *Provider) FindOrInsertUser(identityKey string) (*wdk.FindOrInsertUserResponse, error) {
	user, err := p.repo.FindUser(identityKey)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user != nil {
		return &wdk.FindOrInsertUserResponse{
			User:  *user,
			IsNew: false,
		}, nil
	}

	newUser := &models.User{
		OutputBaskets: []*models.OutputBasket{{
			Name:                    "default",
			NumberOfDesiredUTXOs:    32,
			MinimumDesiredUTXOValue: 1000,
		}},
	}
	newUser.IdentityKey = identityKey

	settings, err := p.repo.ReadSettings()
	if err != nil {
		return nil, fmt.Errorf("failed to read settings: %w", err)
	}

	newUser.ActiveStorage = settings.StorageIdentityKey

	user, err = p.repo.CreateUser(newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return &wdk.FindOrInsertUserResponse{
		User:  *user,
		IsNew: true,
	}, nil
}

// CreateAction Storage level processing for wallet `createAction`.
func (p *Provider) CreateAction(auth wdk.AuthID, args wdk.ValidCreateActionArgs) (*wdk.StorageCreateActionResult, error) {
	if auth.UserID == nil {
		return nil, fmt.Errorf("missing user ID")
	}
	if err := validate.ValidCreateActionArgs(&args); err != nil {
		return nil, fmt.Errorf("invalid createAction args: %w", err)
	}

	res, err := p.actions.Create(auth, actions.FromValidCreateActionArgs(&args))
	if err != nil {
		return nil, fmt.Errorf("failed to create action: %w", err)
	}
	return res, nil
}
