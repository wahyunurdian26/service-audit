package transport

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/wahyunurdian26/service-audit/endpoint"
	"github.com/wahyunurdian26/service-audit/repository/postgres"
	"github.com/wahyunurdian26/service-audit/service"
	"github.com/wahyunurdian26/service-audit/config"
)

type AmqpServer struct {
	conn      *amqp.Connection
	dbPool    *pgxpool.Pool
	endpoints endpoint.AuditEndpoints
	close     func()
}

func (s *AmqpServer) Endpoints() endpoint.AuditEndpoints {
	return s.endpoints
}

func NewAMQPServer() *AmqpServer {
	cfg := config.LoadConfigs()
	rabbitURL := cfg.MessageBrokerConfig.RabbitMQUrl
	dbURL := cfg.DBConfiguration.DatabaseUrl

	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	_, err = dbPool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS audit_logs (
			id SERIAL PRIMARY KEY,
			transaction_id VARCHAR(50) NOT NULL,
			account_id VARCHAR(50),
			amount NUMERIC,
			merchant_name VARCHAR(100),
			status VARCHAR(20),
			event_type VARCHAR(50),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	repo := postgres.NewAuditRepository(dbPool)
	svc := service.NewAuditService(repo)
	endpoints := endpoint.MakeEndpoints(svc)

	return &AmqpServer{
		conn:      conn,
		dbPool:    dbPool,
		endpoints: endpoints,
		close: func() {
			dbPool.Close()
			conn.Close()
		},
	}
}
