package repository

import (
	"context"
	"github.com/wahyunurdian26/service-audit/model"
)

type AuditRepository interface {
	Save(ctx context.Context, log *model.AuditLog) error
	GetListAuditLog(ctx context.Context, page, show int, itemId, itemType string) ([]model.AuditLog, int, error)
}
