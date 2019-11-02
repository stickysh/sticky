package inmem

import (
	"fmt"
	"github.com/stickysh/sticky/trigger"
	"sync"
	"time"
)

type scheduleRepo struct {
	mtx sync.RWMutex
	triggers map[trigger.ScheduleID]*trigger.Schedule
}

func NewScheduleRepo() *scheduleRepo {
	return &scheduleRepo{
		triggers: make(map[trigger.ScheduleID]*trigger.Schedule, 0),
	}
}

func (r *scheduleRepo) Store(t *trigger.Schedule) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.triggers[t.ID] = t
	return nil
}

func (r *scheduleRepo) FindNext(dur time.Duration) *trigger.Schedule {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	var item *trigger.Schedule
	minDur := dur
	for _, val := range r.triggers {
		curDur := val.Start.Sub(time.Now())
		if curDur < minDur {
			minDur = curDur
			item = val
		}
	}

	return item
}

func (r *scheduleRepo) Find(id trigger.ScheduleID) (*trigger.Schedule, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	if val, ok := r.triggers[id]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("not found")
}

func (r *scheduleRepo) FindAllByAction(name string) []*trigger.Schedule {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	items := make([]*trigger.Schedule, 0, len(r.triggers))
	for _, val := range r.triggers {
		if val.Action == name {
			items = append(items, val)
		}
	}

	return items
}

func (r *scheduleRepo) FindAllByTime(time time.Time) []*trigger.Schedule {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	items := make([]*trigger.Schedule, 0, len(r.triggers))
	for _, val := range r.triggers {
		if val.Start == time {
			items = append(items, val)
		}
	}

	return items
}

