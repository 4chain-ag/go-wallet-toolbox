package repo

import (
	"fmt"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/database/models"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"gorm.io/gorm"
)

type Certificates struct {
	db *gorm.DB
}

type ListCertificatesOptions struct {
	SerialNumber       *wdk.Base64String
	Subject            *wdk.PubKeyHex
	RevocationOutpoint *wdk.OutpointString
	Signature          *wdk.HexString
	Certifiers         []wdk.PubKeyHex
	Types              []wdk.Base64String
	Limit              wdk.PositiveIntegerDefault10Max10000
	Offset             wdk.PositiveIntegerOrZero
}

func NewCertificates(db *gorm.DB) *Certificates {
	return &Certificates{db: db}
}

func (c *Certificates) CreateCertificate(certificate *models.Certificate) (uint, error) {
	err := c.db.Create(certificate).Error
	if err != nil {
		return 0, fmt.Errorf("failed to create certificate model: %w", err)
	}
	return certificate.ID, nil
}

func (c *Certificates) DeleteCertificate(userID int, args wdk.RelinquishCertificateArgs) error {
	err := c.db.Delete(&models.Certificate{}, "type = ? AND serial_number = ? AND certifier = ? AND user_id = ?", args.Type, args.SerialNumber, args.Certifier, userID).Error
	if err != nil {
		return fmt.Errorf("failed to delete certificate model: %w", err)
	}

	return nil
}

func (c *Certificates) ListAndCountCertificates(userID int, opts ListCertificatesOptions) ([]*models.Certificate, int64, error) {
	var certificates []*models.Certificate
	var totalRows int64

	query := c.db.Preload("CertificateFields").Where("user_id = ?", userID)

	// Add optional filters based on provided options
	if opts.SerialNumber != nil {
		query = query.Where("serial_number = ?", opts.SerialNumber)
	}
	if opts.Subject != nil {
		query = query.Where("subject = ?", opts.Subject)
	}
	if opts.RevocationOutpoint != nil {
		query = query.Where("revocation_outpoint = ?", opts.RevocationOutpoint)
	}
	if opts.Signature != nil {
		query = query.Where("signature = ?", opts.Signature)
	}
	if len(opts.Certifiers) > 0 {
		query = query.Where("certifier IN ?", opts.Certifiers)
	}
	if len(opts.Types) > 0 {
		query = query.Where("type IN ?", opts.Types)
	}

	// First count results
	query.Model(&models.Certificate{}).Count(&totalRows)

	// Apply pagination
	if opts.Limit > 0 {
		// limit is maxed 10000 so shouldn't overflow
		//nolint:gosec
		query = query.Limit(int(opts.Limit))
	}
	//nolint:gosec
	query = query.Offset(int(opts.Offset))

	err := query.Find(&certificates).Error
	if err != nil {
		return nil, -1, fmt.Errorf("failed to list certificates: %w", err)
	}

	return certificates, totalRows, nil
}
