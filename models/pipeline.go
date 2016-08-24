package models

import "time"

// PipelineTemplate template for request data
type PipelineTemplate struct {
	Kind       string    `json:"kind" binding:"In(CCloud)"`
	APIVersion string    `json:"apiVersion" binding:"In(v1)"`
	MetaData   *Pipeline `json:"metadata" binding:"Required"`
}

// Pipeline pipeline data
type Pipeline struct {
	ID        int64      `json:"id" gorm:"primary_key"`
	Namespace string     `json:"namespace" binding:"Required" gorm:"type:int;not null;idxs_namespace_name"`
	Name      string     `json:"name" binding:"Required" gorm:"type:int;not null;idxs_namespace_name"`
	Timeout   uint64     `json:"timeout" gorm:"type:int;"`
	Stages    []*Stage   `json:"stages" grom:"-"`
	Points    []*Point   `json:"points" grom:"-"`
	Status    uint       `json:"Status" gorm:"type:int;default:0"`
	Created   *time.Time `json:"created" `
	Updated   *time.Time `json:"updated"`
	Deleted   *time.Time `json:"deleted"`
}

// PipelineVersion data
type PipelineVersion struct {
	ID       int64      `json:"id" gorm:"primary_key"`
	PID      int64      `json:"Pid" gorm:"type:int;not null;index"`
	Status   uint       `json:"status" gorm:"type:int;default:0"`
	Detail   string     `json:"detail" gorm:"type:text;"`
	Created  *time.Time `json:"created" `
	Updated  *time.Time `json:"updated"`
	Deleted  *time.Time `json:"deleted"`
	Pipeline Pipeline   `json:"-"`
}

// PipelineResult pipeline result
type PipelineResult struct {
	PID       uint   `json:"pid"`
	PvID      uint   `json:"pvid"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}
