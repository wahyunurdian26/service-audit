package service

import (
	"context"
	"testing"

	"github.com/wahyunurdian26/service-audit/mock"
	"github.com/wahyunurdian26/service-audit/model"
	"go.uber.org/mock/gomock"
)

func TestRecordEvent(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockAuditRepository(ctrl)
	svc := NewAuditService(mockRepo)

	event := &model.AuditLog{
		AccountID:    "123456",
		Amount:       50000,
		MerchantName: "Merchant Shopee",
		Status:       "SUCCESS",
		EventType:    "PaymentCompleted",
		WhodunnitName: "Wahyu Nurdian",
	}

	// Expect Save to be called
	mockRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)

	err := svc.RecordEvent(ctx, event)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
