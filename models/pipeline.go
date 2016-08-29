package models

import (
	"fmt"
	"time"
)

// PipelineTemplate template for request data
type PipelineTemplate struct {
	Kind       string    `json:"kind" binding:"In(CCloud)"`
	APIVersion string    `json:"apiVersion" binding:"In(v1)"`
	MetaData   *Pipeline `json:"metadata" binding:"Required"`
}

// Pipeline data
type Pipeline struct {
	ID        uint64     `json:"id" gorm:"primary_key"`
	Namespace string     `json:"namespace" binding:"Required" gorm:"type:varchar(20);not null;unique_index:idxs_namespace_name"`
	Name      string     `json:"name" binding:"Required" gorm:"type:varchar(20);not null;unique_index:idxs_namespace_name"`
	Timeout   uint64     `json:"timeout" gorm:"type:int;"`
	Status    uint       `json:"Status" gorm:"type:tinyint;default:0"`
	CreatedAt *time.Time `json:"created" `
	UpdatedAt *time.Time `json:"updated"`
	DeletedAt *time.Time `json:"deleted"`
	Stages    []*Stage   `json:"stages" binding:"Required" sql:"-"`
	Points    []*Point   `json:"points" binding:"Required" sql:"-"`
}

// PipelineVersion data
type PipelineVersion struct {
	ID        uint64     `json:"id" gorm:"primary_key"`
	PID       uint64     `json:"Pid" gorm:"type:int;not null;index"`
	State     string     `json:"state" gorm:"column:state;type:varchar(20)"`
	Detail    string     `json:"detail" gorm:"type:text;"`
	Status    uint       `json:"status" gorm:"type:tinyint;default:0"`
	CreatedAt *time.Time `json:"created" `
	UpdatedAt *time.Time `json:"updated"`
	DeletedAt *time.Time `json:"deleted"`
	MetaData  *Pipeline  `json:"-" sql:"-"`
}

// PipelineResult data
type PipelineResult struct {
	PID       uint64 `json:"pid"`
	PvID      uint64 `json:"pvid"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Detail    string `json:"detail"`
}

// TableName pipeline table name in db
func (p *Pipeline) TableName() string {
	return "pipeline"
}

// Insert insert new pipeline in db
func (p *Pipeline) Insert() error {
	db.Begin()
	if err := db.Create(p).Error; err != nil || p.ID == 0 {
		db.Rollback()
		return err
	}

	for _, stage := range p.Stages {
		stage.PID = p.ID
		stage.Namespace = p.Namespace
		if err := stage.Insert(); err != nil || stage.ID == 0 {
			db.Rollback()
			return err
		}
	}

	for _, point := range p.Points {
		point.PID = p.ID
		if err := point.Insert(); err != nil || point.ID == 0 {
			db.Rollback()
			return err
		}
	}
	db.Commit()
	return nil
}

// Read read pipeline in db
func (p *Pipeline) Read() error {
	err := db.Find(p).Error
	if err != nil {
		return err
	}

	newStage := &Stage{PID: p.ID}
	p.Stages, err = newStage.ReadStageList()
	if err != nil {
		return err
	}

	newPoint := &Point{PID: p.ID}
	p.Points, err = newPoint.ReadPointList()
	if err != nil {
		return err
	}
	return nil
}

// CheckExist check pipeline exist by name and namespace
func (p *Pipeline) CheckExist() error {
	check := &Pipeline{Namespace: p.Namespace, Name: p.Name}
	if check.Read() == nil && check.ID != 0 {
		return fmt.Errorf("Pipeline = %v in namespane = %v is already exist", p.Name, p.Namespace)
	}
	for _, stage := range p.Stages {
		stage.Namespace = p.Namespace
		if err := stage.CheckExist(); err != nil {
			return err
		}
	}
	return nil
}

// TableName pipeline version table name in db
func (p *PipelineVersion) TableName() string {
	return "pipeline_version"
}

// Insert insert new stage in db
func (p *PipelineVersion) Insert() error {
	return db.Create(p).Error
}

// Read read stage in db
func (p *PipelineVersion) Read() error {
	return db.Find(p).Error
}
