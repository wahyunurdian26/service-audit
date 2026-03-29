# Audit Service

The **Audit Service** provides centralized logging and event tracking for all business actions within the OmniPay microservices infrastructure.

## Responsibility
-   Maintain an immutable record of state changes across all services.
-   Store change logs (before/after states) for accounts and transactions.
-   Expose query interfaces for retrieving audit trails.

## Tech Stack
-   **Language**: Go 1.24+
-   **Transport**: AMQP (Subscriber)
-   **Database**: PostgreSQL

## Configuration
-   `SERVICE_PORT`: Internal listening port (default: `8083`).
-   `RABBITMQ_URL`: AMQP connection for event consumption.
-   `DATABASE_URL`: Persistence parameters for audit records.

## Data Persistence
The service stores events in an `audits` table, including metadata like `whodunnit` (user), `event_type`, and a detailed payload of the modified entity.

## Local Development
```bash
go run main.go
```
Requires RabbitMQ to receive events from other microservices.
