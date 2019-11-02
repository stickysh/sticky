package inmem

import (
	"github.com/stickysh/sticky/action"
	"sync"
)

type actionRepo struct {
	mtx sync.RWMutex
	actions map[string]*action.Action
}

func NewActionRepo() *actionRepo {
	return &actionRepo{
		actions: make(map[string]*action.Action, 0),
	}
}

func (r *actionRepo) Store(a *action.Action) error {

	return nil
}

func (r *actionRepo) FindAll() []*action.Action {
	return nil
}
