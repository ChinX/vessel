package models

import (
	"time"
)

const (
	// StarPoint point type
	StarPoint = "Start"
	// CheckPoint point type
	CheckPoint = "Check"
	// EndPoint point type
	EndPoint = "End"

	StartPointMark = "$StartPointMark$"
	EndPointMark   = "$EndPointMark$"
)

// Point data
type Point struct {
	ID         int64      `json:"id" gorm:"primary_key"`
	PID        int64      `json:"pid" gorm:"type:int;not null;index"`
	Type       string     `json:"type" binding:"In(Start,Check,End)" gorm:"size:20"`
	Triggers   string     `json:"triggers" gorm:"varchar(255)"` // eg : "a,b,c"
	Conditions string     `json:"conditions" gorm:"size:255"`   // eg : "a,b,c"
	Status     uint       `json:"Status" gorm:"type:int;default:0"`
	Created    *time.Time `json:"created" `
	Updated    *time.Time `json:"updated"`
	Deleted    *time.Time `json:"deleted"`
}
