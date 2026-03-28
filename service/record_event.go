package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/wahyunurdian26/service-audit/model"
	uLog "github.com/wahyunurdian26/util/logger"
)

func (s *auditService) RecordEvent(ctx context.Context, log *model.AuditLog) error {
	uLog.LogRequest(ctx, "RecordEvent", log)

	err := s.repo.Save(ctx, log)
	if err != nil {
		uLog.LogError(ctx, "s.repo.Save", err)
		return status.Error(codes.Internal, err.Error())
	}
	return nil
}
