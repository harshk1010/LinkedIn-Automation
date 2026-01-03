package state

import "time"

// ProfileState represents persisted state of a LinkedIn profile
type ProfileState struct {
	ProfileURL     string
	Visited        bool
	CreatedAt      time.Time
	LastVisitedAt  *time.Time
	LastActionAt   *time.Time
}
