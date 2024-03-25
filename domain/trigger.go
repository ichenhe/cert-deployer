package domain

// Trigger can monitor specific event and execute deployments when typical event happened.
type Trigger interface {
	GetName() string
	StartMonitoring() error
	// Close closes the trigger. All inner coroutines should be canceled.
	Close()
}
