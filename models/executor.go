package models

const (
	// StateReady  operating state ready
	StateReady uint = 0//"Ready"
	// StateRunning  operating state running
	StateRunning uint = 1//"Running"
	// StateDeleted  operating state deleted
	StateDeleted uint = 2//"Deleted"

	// ResultSuccess  operating result success
	ResultSuccess = "OK"
	// ResultFailed  operating result failed
	ResultFailed = "Error"
	// ResultTimeout  operating result timeout
	ResultTimeout = "Timeout"
)

type Executor interface {
	Start(readyMap map[string]bool, finishChan chan *ExecutedResult) bool
	Stop(readyMap map[string]bool, finishChan chan *ExecutedResult) bool
	GetFrom() []string
}

// ExecutedResult executor operating result
type ExecutedResult struct {
	Name   string
	Status string
	Detail string
}
