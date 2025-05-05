package models

import (
	"encoding/json"
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"gorm.io/datatypes"
)

type ProvenTxReq struct {
	CreatedAt time.Time
	UpdatedAt time.Time

	TxID string `gorm:"type:varchar(64);primaryKey"`

	Status   wdk.ProvenTxReqStatus `gorm:"default:unknown"`
	Attempts uint
	Notified bool

	RawTx     []byte
	InputBeef []byte

	History datatypes.JSON
}

func (p *ProvenTxReq) AddNote(when time.Time, what string, attrs map[string]any) {
	var history HistoryModel
	if p.History != nil {
		// in case of unmarshalling error, we will just create a new history
		_ = json.Unmarshal(p.History, &history)
	}

	note := HistoryNote{
		When:  when,
		What:  what,
		Attrs: attrs,
	}
	history.Notes = append(history.Notes, note)

	data, _ := json.Marshal(history)
	p.History = data
}
