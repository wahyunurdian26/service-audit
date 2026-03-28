package repository

import (
	"context"
	"microservice/audit-service/model"
)

type AuditRepository interface {
	Save(ctx context.Context, log *model.AuditLog) error
	GetListAuditLog(ctx context.Context, page, show int, itemId, itemType string) ([]model.AuditLog, int, error)
}
