package model

import "time"

type Watermark struct {
	Job           string    `json:"job,omitempty" bson:"job"`
	LastTimestamp time.Time `json:"lastTimestamp" bson:"watermark_ts"`
	LastId        int64     `json:"lastId,omitempty" bson:"watermark_id"`
	UpdatedAt     time.Time `json:"updatedAt" bson:"updated_at"`
}
