package models

const (
	// StateReady  state ready
	StateReady = "Ready"
	// StateRunning  state running
	StateRunning = "Running"
	// StateDeleted  state deleted
	StateDeleted = "Deleted"

	// ResultSuccess  result success
	ResultSuccess = "OK"
	// ResultFailed  result failed
	ResultFailed = "Error"
	// ResultTimeout  result timeout
	ResultTimeout = "Timeout"
)

// ExecutedResult executor operating result
type ExecutedResult struct {
	Name   string
	Status string
	Detail string
}
