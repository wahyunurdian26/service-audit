package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/wahyunurdian26/service-audit/model"
)

var (
	saveAuditLogQuery = `
		INSERT INTO audit_logs (
			transaction_id, account_id, amount, merchant_name, status, event_type,
			item_id, item_type, event, whodunnit_id, whodunnit_name, activities
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
)

func (r *auditRepository) Save(ctx context.Context, log *model.AuditLog) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	activitiesJSON, _ := json.Marshal(log.Activities)

	_, err := r.db.Exec(ctxTimeout, saveAuditLogQuery, 
		log.TransactionID, log.AccountID, log.Amount, log.MerchantName, log.Status, log.EventType,
		log.ItemID, log.ItemType, log.Event, log.WhodunnitID, log.WhodunnitName, activitiesJSON,
	)
	return err
}
