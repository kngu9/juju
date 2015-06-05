// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package process

// Status values specific to workload processes.
const (
	StatusPending Status = iota
	StatusActive
	StatusFailed
	StatusStopped
)

// Status represents the status of a worload process.
type Status string

// String implements fmt.Stringer.
func (s Status) String() string {
	switch status {
	case StatusPending:
		return "pending"
	case StatusActive:
		return "active"
	case StatusFailed:
		return "failed"
	case StatusStopped:
		return "stopped"
	}
	return "Unknown"
}
