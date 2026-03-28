package transport

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"microservice/audit-service/model"
	"microservice/util/audit"
	"microservice/util/logger"
	"microservice/util/requestid"
)

func (s *AmqpServer) Run() {
	time.Sleep(10 * time.Second) // Wait for infrastructure

	ch, err := s.conn.Channel()
	if err != nil {
		logger.Error("Failed to open a channel: ", err)
		os.Exit(1)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"audit_queue", // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		logger.Error("Failed to declare a queue: ", err)
		os.Exit(1)
	}

	err = ch.ExchangeDeclare(
		"payment_events", // name
		"topic",          // type
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		logger.Error("Failed to declare an exchange: ", err)
		os.Exit(1)
	}

	err = ch.QueueBind(
		q.Name,              // queue name
		"payment.completed", // routing key
		"payment_events",    // exchange
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		logger.Error("Failed to bind a queue: ", err)
		os.Exit(1)
	}

	msgs, err := ch.Consume(
		q.Name, "", true, false, false, false, nil,
	)
	if err != nil {
		logger.Error("Failed to register a consumer: ", err)
		os.Exit(1)
	}

	go func() {
		for d := range msgs {
			reqID := ""
			if d.Headers != nil {
				if id, ok := d.Headers[string(requestid.RequestIdAttr)].(string); ok {
					reqID = id
				}
			}
			ctx := context.WithValue(context.Background(), requestid.RequestId, reqID)

			logger.LogInfo(ctx, "status: received", "Consume_AuditLog")

			// Try to unmarshal into a generic map first to detect format
			var raw map[string]interface{}
			if err := json.Unmarshal(d.Body, &raw); err != nil {
				logger.Error("Failed to unmarshal raw event: ", err)
				continue
			}

			var event model.AuditLog
			if _, ok := raw["activities"]; ok {
				// Detailed Audit Log Format
				var detailed audit.AuditLog
				json.Unmarshal(d.Body, &detailed)
				
				event.TransactionID = detailed.RequestId
				event.AccountID = detailed.ItemId
				event.ItemID = detailed.ItemId
				event.ItemType = detailed.ItemType
				event.Event = detailed.Event
				event.EventType = detailed.Event
				event.WhodunnitID = detailed.WhodunnitId
				event.WhodunnitName = detailed.WhodunnitName
				event.Activities = detailed.Activities
				event.CreatedAt = time.Now()
			} else {
				// Flat Payment Event Format
				json.Unmarshal(d.Body, &event)
				event.EventType = "PaymentCompleted"
				if event.CreatedAt.IsZero() {
					event.CreatedAt = time.Now()
				}
			}

			_, err := s.endpoints.RecordEventEndpoint(ctx, &event)
			if err != nil {
				logger.LogError(ctx, "Audit RecordEvent Error", err)
			} else {
				logger.Info("Audit log successfully saved for ", event.EventType)
			}
		}
	}()

	logger.Info("Audit Service started successfully - Waiting for messages...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-c
	logger.Info("Shutting down Audit Service...")
	if s.close != nil {
		s.close()
	}
}
