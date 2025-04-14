package testabilities

import (
	"fmt"
	"testing"
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/database/models"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/testabilities/testusers"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/txutils"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"gorm.io/gorm"
)

var FirstCreatedAt = time.Date(2006, 02, 01, 15, 4, 5, 7, time.UTC)

type UserUTXOFixture interface {
	OwnedBy(user testusers.User) UserUTXOFixture
	InBasket(basket *wdk.TableOutputBasket) UserUTXOFixture
	P2PKH() UserUTXOFixture
	WithSatoshis(sats int64) UserUTXOFixture
	Stored()
}

type UTXODatabase interface {
	Save(utxo *models.UserUTXO)
}

var defaultBasket = wdk.TableOutputBasket{
	CreatedAt: FirstCreatedAt,
	UpdatedAt: FirstCreatedAt,
	IsDeleted: false,
	BasketID:  1,
	BasketConfiguration: wdk.BasketConfiguration{
		Name:                    wdk.BasketNameForChange,
		NumberOfDesiredUTXOs:    30,
		MinimumDesiredUTXOValue: 1000,
	},
	UserID: 1,
}

type userUtxoFixture struct {
	parent             UTXODatabase
	t                  testing.TB
	index              uint
	userID             int
	txID               string
	vout               uint32
	satoshis           uint64
	estimatedInputSize uint64
	basket             *wdk.TableOutputBasket
}

func newUtxoFixture(t testing.TB, parent UTXODatabase, index uint) *userUtxoFixture {
	basket := defaultBasket
	return &userUtxoFixture{
		t:                  t,
		parent:             parent,
		index:              index,
		basket:             &basket,
		userID:             1,
		txID:               txIDTemplated(index),
		vout:               uint32(index),
		satoshis:           1,
		estimatedInputSize: txutils.P2PKHEstimatedInputSize,
	}
}

func txIDTemplated(index uint) string {
	return fmt.Sprintf("a%010de1b81dd2c9c0c6cd67f9bdf832e9c2bb12a1d57f30cb6ebbe78d9", index)
}

func (f *userUtxoFixture) InBasket(basket *wdk.TableOutputBasket) UserUTXOFixture {
	f.basket = basket
	return f
}

func (f *userUtxoFixture) OwnedBy(user testusers.User) UserUTXOFixture {
	f.userID = user.ID
	f.basket.UserID = user.ID
	return f
}

func (f *userUtxoFixture) P2PKH() UserUTXOFixture {
	f.estimatedInputSize = txutils.P2PKHEstimatedInputSize
	return f
}

func (f *userUtxoFixture) WithSatoshis(satoshis int64) UserUTXOFixture {
	if satoshis < 0 {
		f.t.Fatalf("satoshis must be a positive number, got %d", satoshis)
	}
	f.satoshis = uint64(satoshis)
	return f
}

func (f *userUtxoFixture) Stored() {
	if f.satoshis == 0 {
		return
	}

	utxo := &models.UserUTXO{
		UserID:             f.userID,
		TxID:               f.txID,
		Vout:               f.vout,
		Satoshis:           f.satoshis,
		EstimatedInputSize: f.estimatedInputSize,
		CreatedAt:          FirstCreatedAt.Add(time.Duration(f.index) * time.Second),
		BasketID:           f.basket.BasketID,
		Basket: &models.OutputBasket{
			CreatedAt:               FirstCreatedAt,
			UpdatedAt:               FirstCreatedAt,
			DeletedAt:               gorm.DeletedAt{},
			BasketID:                f.basket.BasketID,
			Name:                    f.basket.Name,
			UserID:                  f.basket.UserID,
			NumberOfDesiredUTXOs:    f.basket.NumberOfDesiredUTXOs,
			MinimumDesiredUTXOValue: f.basket.MinimumDesiredUTXOValue,
		},
	}

	f.parent.Save(utxo)
}
