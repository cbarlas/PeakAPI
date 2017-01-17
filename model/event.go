// Package model has the structure and representation of the Event data
package model

// Representation of the Event data
type Event struct {
	ApiKey    string    `json:"apiKey"`
	UserID    uint64    `json:"userID"`
	Timestamp int64		`json:"timestamp"`
}

type Events []*Event
