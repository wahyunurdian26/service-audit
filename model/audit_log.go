package model

import (
	"time"
	"github.com/wahyunurdian26/util/audit"
)

type AuditLog struct {
	ID            int64            `json:"id"`
	TransactionID string           `json:"transaction_id"`
	AccountID     string           `json:"account_id"`
	Amount        float64          `json:"amount"`
	MerchantName  string           `json:"merchant_name"`
	Status        string           `json:"status"`
	EventType     string           `json:"event_type"`
	CreatedAt     time.Time        `json:"created_at"`
	
	// Detailed Audit Fields
	ItemID        string           `json:"item_id"`
	ItemType      string           `json:"item_type"`
	Event         string           `json:"event"`
	WhodunnitID   int64            `json:"whodunnit_id"`
	WhodunnitName string           `json:"whodunnit_name"`
	Activities    []audit.Activity `json:"activities"`
}
