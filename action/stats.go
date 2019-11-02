package action

import (
	"time"
	ctrlID "github.com/stickysh/sticky/pkg/id"
)

type InvokStatus string

const (
	Success InvokStatus = "success"
	Running InvokStatus = "running"
	Canceled InvokStatus = "canceled"
	Failure InvokStatus = "failure"
)

type Stat struct {
	ID     string
	Action string
	Start  time.Time
	End    time.Time
	Status InvokStatus
}

func NewStat(name string, start time.Time, status InvokStatus) *Stat{
	return &Stat{
		ID: ctrlID.GenULID(),
		Action: name,
		Start:  start,
		Status: status,
	}
}

func (s *Stat) Success() {
	s.Status = Success
	s.End = time.Now()
}

type StatsRepo interface {
	Store(stat *Stat) error
	FindAll() []*Stat
	FindByAction(actionID string) []*Stat
}
