package structure

import "time"

type TimeFields struct {
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	StartedAt time.Time `json:"started_at" bson:"started_at"`
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at"`
}
