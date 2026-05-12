// Package repository provides PostgreSQL data-access operations for the
// company domain using pgx/v5 connection pooling.
package repository

import (
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository provides all PostgreSQL CRUD operations for the eiker-company-db service.
type Repository struct {
	// pool is the pgx connection pool shared by all operations.
	pool *pgxpool.Pool
}

// New constructs a Repository backed by the given connection pool.
func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// jsonMarshalNullable returns nil when m is empty, otherwise the JSON-encoded bytes.
// This avoids storing an explicit "null" literal in JSONB columns when no data is provided.
func jsonMarshalNullable(m map[string]interface{}) ([]byte, error) {
	if len(m) == 0 {
		return nil, nil
	}
	return json.Marshal(m)
}
