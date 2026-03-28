package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"microservice/audit-service/repository"
	"microservice/util/logger"
)

const (
	timeout = 5
)

type auditRepository struct {
	db *pgxpool.Pool
}

func NewAuditRepository(db *pgxpool.Pool) repository.AuditRepository {
	// Auto-migrate schema
	migrationQuery := `
		ALTER TABLE audit_logs 
		ADD COLUMN IF NOT EXISTS item_id VARCHAR(255),
		ADD COLUMN IF NOT EXISTS item_type VARCHAR(255),
		ADD COLUMN IF NOT EXISTS event VARCHAR(255),
		ADD COLUMN IF NOT EXISTS whodunnit_id BIGINT,
		ADD COLUMN IF NOT EXISTS whodunnit_name VARCHAR(255),
		ADD COLUMN IF NOT EXISTS activities JSONB;
	`
	_, err := db.Exec(context.Background(), migrationQuery)
	if err != nil {
		logger.Warn("Warning: failed to auto-migrate audit_logs table: ", err)
	}

	return &auditRepository{db: db}
}
