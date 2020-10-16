package models

import (
	"database/sql"
	"encoding/json"
)

type JobStatus int

const (
	PENDING = iota
	DONE
	FAILED
)

type Job struct {
	ID         string          `json:"id" db:"id"`
	Type       string          `json:"type" db:"type"`
	Status     JobStatus       `json:"status,omitempty" db:"status"`
	Payload    json.RawMessage `json:"payload" db:"payload"`
	EntityType sql.NullString  `json:"entity_type,omitempty" db:"entity_type"`
	EntityID   sql.NullString  `json:"entity_id,omitempty" db:"entity_id"`
}
