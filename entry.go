package confucius

import "sync"

const (
	Undefined Status = iota
	Inactive
	Ok
	Serving
	Stopping
	Stopped
)

type Status uint8

// containerEntry is a wrapper for services, which is used to name a service, handle its status and mutes it
type containerEntry struct {
	sync.Mutex
	service Service
	name    string
	status  Status
}

func (e *containerEntry) getStatus() Status {
	e.Lock()
	defer e.Unlock()

	return e.status
}

func (e *containerEntry) setStatus(status Status) {
	e.Lock()
	e.status = status
	e.Unlock()
}

func (e *containerEntry) hasStatus(status Status) bool {
	return status == e.getStatus()
}
