package inmem

import (
	"github.com/stickysh/sticky/action"
	"sync"
)


type statsRepo struct {
	mtx sync.RWMutex
	stats map[string]*action.Stat
}

func NewStatsRepo() *statsRepo {
	return &statsRepo{
		stats: make(map[string]*action.Stat, 0),
	}
}

func (r *statsRepo) Store(stat *action.Stat) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.stats[stat.ID] = stat

	return nil
}

func (r *statsRepo) FindAll() []*action.Stat {

	stats := make([]*action.Stat, 0, len(r.stats))

	for _, v := range r.stats {
		stats = append(stats, v)
	}

	return stats
}

func (r *statsRepo) FindByAction(actionID string) []*action.Stat {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	list := make([]*action.Stat, 0)

	for _, val := range r.stats {
		if val.Action == actionID {
			list = append(list, val)
		}
	}

	return list
}


