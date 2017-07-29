package interfaces

// ServiceEvent contains the information related to an event that happened
// to a service.
type ServiceEvent struct {
	Name string
}

// ServiceWatcher defines an interface for a service watcher against which
// it is possible to subscribe for service changes.
type ServiceWatcher interface {
	ServiceRepository

	// Subscribe allows providing a channel through which service events will
	// be distributed.
	Subscribe(events chan<- ServiceEvent)
}
