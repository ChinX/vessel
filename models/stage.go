package models

import (
	"time"

	"k8s.io/kubernetes/pkg/api/v1"
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
	ID            int64       `json:"-"`
	Name          string      `json:"name"  binding:"Required" gorm:"varchar(30)"`
	Type          string      `json:"type"  binding:"In(container,vm,pc)" gorm:"varchar(20)"`
	Dependencies  string      `json:"dependencies,omitempty" gorm:"varchar(255)"` // eg : "a,b,c"
	Artifacts     []*Artifact `json:"artifacts" binding:"Required" gorm:"-"`
	Volumes       []*Volume   `json:"volumes,omitempty" gorm:"-"`
	ArtifactsJSON string      `json:"-" gorm:"column:artifacts;type:text;not null"` // json type
	VolumesJSON   string      `json:"-" gorm:"column:volumes;type:text;not null"`   // json type
	Status        uint        `json:"status" gorm:"type:int;default:0"`
	Created       *time.Time  `json:"created" `
	Updated       *time.Time  `json:"updated"`
	Deleted       *time.Time  `json:"deleted"`
}

// StageVersion data
type StageVersion struct {
	ID      int64      `json:"id" gorm:"primary_key"`
	PvID    int64      `json:"pvid" gorm:"type:int;not null"`
	SID     int64      `json:"sid" gorm:"type:int;not null"`
	Detail  string     `json:"detail" gorm:"type:text;"`
	Status  uint       `json:"status" gorm:"type:int;default:0"`
	Created *time.Time `json:"created" `
	Updated *time.Time `json:"updated"`
	Deleted *time.Time `json:"deleted"`
	Stage   Stage      `json:"-"`
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
