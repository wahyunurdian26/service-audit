package service

import (
	"context"

	"microservice/audit-service/model"
	"microservice/audit-service/repository"
)

type AuditService interface {
	RecordEvent(ctx context.Context, log *model.AuditLog) error
	GetListAuditLog(ctx context.Context, page, show int, itemId, itemType string) ([]model.AuditLog, int, error)
}

type auditService struct {
	repo repository.AuditRepository
}

func NewAuditService(repo repository.AuditRepository) AuditService {
	return &auditService{repo: repo}
}
