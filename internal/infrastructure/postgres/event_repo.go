package postgres

import (
	"context"
	"encoding/json"
	"strconv"

	"behavix-ai/internal/domain/event"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// EventRepository implements event.Repository using PostgreSQL.
type EventRepository struct {
	pool *pgxpool.Pool
}

// NewEventRepository returns an event repository backed by PostgreSQL.
func NewEventRepository(pool *pgxpool.Pool) *EventRepository {
	return &EventRepository{pool: pool}
}

// InsertBatch inserts a batch of event records.
func (r *EventRepository) InsertBatch(ctx context.Context, records []event.Record) error {
	if len(records) == 0 {
		return nil
	}
	for _, rec := range records {
		props, _ := json.Marshal(rec.Properties)
		if props == nil {
			props = []byte("{}")
		}
		ctxVal, _ := json.Marshal(rec.Context)
		if ctxVal == nil {
			ctxVal = []byte("{}")
		}
		_, err := r.pool.Exec(ctx,
			`INSERT INTO events (id, tenant_id, user_id, session_id, event_name, event_timestamp, received_at, properties, context)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			rec.ID, rec.TenantID, rec.UserID, rec.SessionID, rec.EventName, rec.EventTimestamp, rec.ReceivedAt, props, ctxVal,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// List returns events for a tenant with optional filters.
func (r *EventRepository) List(ctx context.Context, tenantID uuid.UUID, limit int, eventName, userID string) ([]event.Record, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	query := `SELECT id, tenant_id, user_id, COALESCE(session_id, ''), event_name, event_timestamp, received_at, properties, context
	          FROM events WHERE tenant_id = $1`
	args := []interface{}{tenantID}
	n := 2
	if eventName != "" {
		query += ` AND event_name = $` + strconv.Itoa(n)
		args = append(args, eventName)
		n++
	}
	if userID != "" {
		query += ` AND user_id = $` + strconv.Itoa(n)
		args = append(args, userID)
		n++
	}
	query += ` ORDER BY received_at DESC LIMIT $` + strconv.Itoa(n)
	args = append(args, limit)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []event.Record
	for rows.Next() {
		var rec event.Record
		var props, ctxVal []byte
		err := rows.Scan(&rec.ID, &rec.TenantID, &rec.UserID, &rec.SessionID, &rec.EventName, &rec.EventTimestamp, &rec.ReceivedAt, &props, &ctxVal)
		if err != nil {
			return nil, err
		}
		_ = json.Unmarshal(props, &rec.Properties)
		_ = json.Unmarshal(ctxVal, &rec.Context)
		if rec.Properties == nil {
			rec.Properties = make(map[string]interface{})
		}
		if rec.Context == nil {
			rec.Context = make(map[string]interface{})
		}
		out = append(out, rec)
	}
	return out, rows.Err()
}

// Compile-time check that EventRepository implements event.Repository.
var _ event.Repository = (*EventRepository)(nil)
