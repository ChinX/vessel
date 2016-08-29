package models

import "time"

const (
	// StarPoint point type
	StartPoint = "Start"
	// CheckPoint point type
	CheckPoint = "Check"
	// EndPoint point type
	EndPoint = "End"

	// StartPointMark start mark
	StartPointMark = "$StartPointMark$"
	// EndPointMark end mark
	EndPointMark = "$EndPointMark$"
)

// Point data
type Point struct {
	ID         uint64     `json:"id" gorm:"primary_key"`
	PID        uint64     `json:"pid" gorm:"type:int;not null;index"`
	Type       string     `json:"type" binding:"In(Start,Check,End)" gorm:"size:20"`
	Triggers   string     `json:"triggers" gorm:"varchar(255)"` // eg : "a,b,c"
	Conditions string     `json:"conditions" gorm:"size:255"`   // eg : "a,b,c"
	Status     uint       `json:"Status" gorm:"type:tinyint;default:0"`
	CreatedAt  *time.Time `json:"created" `
	UpdatedAt  *time.Time `json:"updated"`
	DeletedAt  *time.Time `json:"deleted"`
}

// PointVersion data
type PointVersion struct {
	ID           uint64        `json:"id" gorm:"primary_key"`
	PvID         uint64        `json:"pvid" gorm:"type:int;not null"`
	PointID      uint64        `json:"PointID" gorm:"type:int;not null;index"`
	State        string        `json:"state" gorm:"column:versionStatus;type:varchar(20);not null;"`
	Detail       string        `json:"detail" gorm:"type:text;"`
	Status       uint          `json:"status" gorm:"type:tinyint;default:0"`
	CreatedAt    *time.Time    `json:"created" `
	UpdatedAt    *time.Time    `json:"updated"`
	DeletedAt    *time.Time    `json:"deleted"`
	Conditions   []string      `json:"-" sql:"-"`
	MateDate     *Point        `json:"-" sql:"-"`
	StageVersion *StageVersion `json:"-" sql:"-"`
}

// TableName point table name in db
func (p *Point) TableName() string {
	return "point"
}

// TableName point version table name in db
func (p *PointVersion) TableName() string {
	return "point_version"
}

// Insert insert new point in db
func (p *Point) Insert() error {
	return db.Create(p).Error
}

// Read read point in db
func (p *Point) Read() error {
	return db.Find(p).Error
}

// ReadPointList Read point list by pid
func (p *Point) ReadPointList() ([]*Point, error) {
	pointList := make([]*Point, 0, 10)
	err := db.Find(&pointList, p).Error
	return pointList, err
}
