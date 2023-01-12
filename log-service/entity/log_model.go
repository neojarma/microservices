package entity

import "time"

type LogEntry struct {
	Name      string    `json:"name" bson:"name"`
	Data      string    `json:"data" bson:"data"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}
