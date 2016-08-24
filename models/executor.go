package models

const (
	// StateReady  operating state ready
	StateReady = "Ready"
	// StateRunning  operating state running
	StateRunning = "Running"
	// StateDeleted  operating state deleted
	StateDeleted = "Deleted"

	// ResultSuccess  operating result success
	ResultSuccess = "OK"
	// ResultFailed  operating result failed
	ResultFailed = "Error"
	// ResultTimeout  operating result timeout
	ResultTimeout = "Timeout"
)

// ExecutedResult executor operating result
type ExecutedResult struct {
	Name   string
	Status string
	Result interface{}
}
