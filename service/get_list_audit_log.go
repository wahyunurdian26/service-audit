package service

import (
	"context"
	
	"microservice/audit-service/model"
	uLog "microservice/util/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (svc *auditService) GetListAuditLog(ctx context.Context, page, show int, itemId, itemType string) ([]model.AuditLog, int, error) {
	if page < 1 {
		page = 1
	}
	if show < 1 {
		show = 10
	}

	logs, total, err := svc.repo.GetListAuditLog(ctx, page, show, itemId, itemType)
	if err != nil {
		uLog.LogError(ctx, "svc.repo.GetListAuditLog", err)
		return nil, 0, status.Error(codes.Internal, err.Error())
	}

	return logs, total, nil
}
