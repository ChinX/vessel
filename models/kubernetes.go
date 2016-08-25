package models

const (
	// LabelKey K8S label key
	LabelKey = "app"
	// WatchAdded K8S watch added
	WatchAdded = "ADDED"
	// WatchDeleted K8S watch deleted
	WatchDeleted = "DELETED"
)

// K8SRes call result from K8S
type K8SRes struct {
	Result string
	Detail string
}
