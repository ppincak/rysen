package core

// Base worker interface
type Worker interface {

	// Start the worker, should be run as a goroutine as it is blocking
	Start()

	// Stop the worker
	Stop()
}
