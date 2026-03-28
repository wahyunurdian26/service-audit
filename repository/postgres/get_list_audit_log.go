package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"microservice/audit-service/model"
)

var (
	baseGetListAuditLogQuery = `
		SELECT count(*) OVER() AS total_count, id, transaction_id, account_id, amount, merchant_name, status, event_type, created_at,
		       COALESCE(item_id, ''), COALESCE(item_type, ''), COALESCE(event, ''), COALESCE(whodunnit_id, 0), COALESCE(whodunnit_name, ''), activities
		FROM audit_logs
	`
)

func (r *auditRepository) GetListAuditLog(ctx context.Context, page, show int, itemId, itemType string) ([]model.AuditLog, int, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	offset := (page - 1) * show
	
	where := []string{}
	args := []interface{}{}
	argIdx := 1

	if itemId != "" {
		where = append(where, fmt.Sprintf("item_id = $%d", argIdx))
		args = append(args, itemId)
		argIdx++
	}
	if itemType != "" {
		where = append(where, fmt.Sprintf("lower(item_type) = $%d", argIdx))
		args = append(args, strings.ToLower(itemType))
		argIdx++
	}

	query := baseGetListAuditLogQuery
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}
	
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, show, offset)

	rows, err := r.db.Query(ctxTimeout, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []model.AuditLog
	var total int
	for rows.Next() {
		var log model.AuditLog
		var activitiesJSON []byte
		err := rows.Scan(
			&total, &log.ID, &log.TransactionID, &log.AccountID, &log.Amount, &log.MerchantName, &log.Status, &log.EventType, &log.CreatedAt,
			&log.ItemID, &log.ItemType, &log.Event, &log.WhodunnitID, &log.WhodunnitName, &activitiesJSON,
		)
		if err != nil {
			return nil, 0, err
		}
		if len(activitiesJSON) > 0 && string(activitiesJSON) != "null" {
			json.Unmarshal(activitiesJSON, &log.Activities)
		}
		logs = append(logs, log)
	}
	return logs, total, nil
}
