package models

import (
	"fmt"

	"github.com/containerops/vessel/utils/timer"
	"k8s.io/kubernetes/pkg/api/v1"
	"time"
)

const (
	// StageContainer stage type
	StageContainer = "container"
	// StageVM stage type
	StageVM = "vm"
	// StagePC stage type
	StagePC = "pc"
)

// Stage stage data
type Stage struct {
	ID                  uint64
	PID                 uint64
	Namespace           string           `json:"namespace" gorm:"type:varchar(20);not null;unique_index:idxs_namespace_name"`
	Name                string           `json:"name" binding:"Required" gorm:"type:varchar(20);not null;unique_index:idxs_namespace_name"`
	PipelineName        string           `json:"-" gorm:"-"`
	Replicas            uint64           `json:"replicas" binding:"Required"`
	Dependencies        string           `json:"dependence"`
	StatusCheckURL      string           `json:"statusCheckLink"`
	StatusCheckInterval uint64           `json:"statusCheckInterval"`
	StatusCheckCount    uint64           `json:"statusCheckCount"`
	Image               string           `json:"image" binding:"Required"`
	Port                uint64           `json:"port" binding:"Required"`
	EnvName             string           `json:"envName"`
	EnvValue            string           `json:"envValue"`
	Status              uint             `json:"-"`
	Hourglass           *timer.Hourglass `json:"-" gorm:"-"`
}

// StageVersion data
type StageVersion struct {
	ID           uint64     `json:"id" gorm:"primary_key"`
	PvID         uint64     `json:"pvid" gorm:"type:int;not null"`
	SID          uint64     `json:"sid" gorm:"type:int;not null"`
	State        string     `json:"state" gorm:"column:state;type:varchar(20)"`
	Detail       string     `json:"detail" gorm:"type:text;"`
	Status       uint       `json:"status" gorm:"type:tinyint;default:0"`
	Created      *time.Time `json:"created" `
	Updated      *time.Time `json:"updated"`
	Deleted      *time.Time `json:"deleted"`
	MateDate     *Stage     `json:"-" gorm:"-"`
	Dependencies []string   `json:"-" gorm:"-"`
}

// Artifact data
type Artifact struct {
	Name      string     `json:"name"`
	Path      string     `json:"path"`
	Lifecycle *Lifecycle `json:"lifecycle,omitempty"`
	Container *Container `json:"container,omitempty"`
}

// Lifecycle data
type Lifecycle struct {
	Before  []string `json:"before,omitempty"`
	Runtime []string `json:"runtime,omitempty"`
	After   []string `json:"after,omitempty"`
}

// Volume data
type Volume struct {
	Name                 string                               `json:"name"`
	HostPath             *v1.HostPathVolumeSource             `json:"hostPath,omitempty"`
	EmptyDir             *v1.EmptyDirVolumeSource             `json:"emptyDir,omitempty"`
	AWSElasticBlockStore *v1.AWSElasticBlockStoreVolumeSource `json:"awsElasticBlockStore,omitempty"`
	CephFS               *v1.CephFSVolumeSource               `json:"cephfs,omitempty"`
}

// Container data
type Container struct {
	WorkingDir     string             `json:"workingDir,omitempty"`
	Ports          []v1.ContainerPort `json:"ports,omitempty"`
	Env            []v1.EnvVar        `json:"env,omitempty"`
	VolumeMounts   []v1.VolumeMount   `json:"volumeMounts,omitempty"`
	LivenessProbe  *v1.Probe          `json:"livenessProbe,omitempty"`
	ReadinessProbe *v1.Probe          `json:"readinessProbe,omitempty"`
	PullPolicy     v1.PullPolicy      `json:"PullPolicy,omitempty"`
	Stdin          bool               `json:"stdin,omitempty"`
	StdinOnce      bool               `json:"stdinOnce,omitempty"`
	TTY            bool               `json:"tty,omitempty"`
}

// TableName stage table name in db
func (s *Stage) TableName() string {
	return "stage"
}

// Insert insert new stage in db
func (s *Stage) Insert() error {
	if err := s.objToJSON(); err != nil {
		return err
	}
	return db.Create(s).Error
}

// Read read stage in db
func (s *Stage) Read() error {
	if err := db.Find(s).Error; err != nil {
		return err
	}
	return s.jsonToObj()
}

// ReadStageList Read stage list by pid
func (s *Stage) ReadStageList() ([]*Stage, error) {
	stageList := make([]*Stage, 0, 10)
	if err := db.Find(&stageList, s).Error; err != nil {
		return nil, err
	}
	for _, stage := range stageList {
		if err := stage.jsonToObj(); err != nil {
			return nil, err
		}
	}
	return stageList, nil
}

// CheckExist check stage exist by name and namespace
func (s *Stage) CheckExist() error {
	check := &Stage{Namespace: s.Namespace, Name: s.Name}
	if check.Read() == nil && check.ID != 0 {
		return fmt.Errorf("Stage = %v in namespane = %v is already exist", s.Name, s.Namespace)
	}
	return nil
}

func (s *Stage) objToJSON() error {
	return nil
}

func (s *Stage) jsonToObj() error {
	return nil
}

// TableName stage version table name in db
func (s *StageVersion) TableName() string {
	return "stage_version"
}

// Insert insert new stage in db
func (s *StageVersion) Insert() error {
	return db.Create(s).Error
}

// Read read stage in db
func (s *StageVersion) Read() error {
	return db.Find(s).Error
}
