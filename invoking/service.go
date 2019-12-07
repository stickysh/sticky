package invoking

import (
	"fmt"
	"github.com/stickysh/sticky/invoke"
	"time"

	"github.com/stickysh/sticky/action"
	"github.com/stickysh/sticky/trigger"
)

type EventAction int

const (
	EventActionDone = iota
)

type Service interface {
	Run(name string, params map[string]interface{}) (interface{}, error)

	AddSchedule(name string, schedID trigger.ScheduleID) error
	RemoveSchedule(name string, schedID trigger.ScheduleID) error


	TriggerActionWithWebhook(action string, webhookID trigger.WebhookID, params map[string]interface{}) error

}

type service struct {
	actionRepo   action.ActionRepo
	scheduleRepo trigger.ScheduleRepo
	statsRepo    action.StatsRepo
	actionTimer  *invoke.ActionTimer
	provider     *invoke.ActionProvider
}

func NewService(actRepo action.ActionRepo, scheduleRepo trigger.ScheduleRepo, statsRepo action.StatsRepo, actionTimer *invoke.ActionTimer, provider *invoke.ActionProvider) *service {
	return &service{
		actionRepo:   actRepo,
		scheduleRepo: scheduleRepo,
		statsRepo:    statsRepo,
		actionTimer:  actionTimer,
		provider:     provider,
	}
}


func (s *service) Run(name string, params map[string]interface{}) (interface{}, error) {
	stat := action.NewStat(name, time.Now(), action.Running)
	defer func(stat *action.Stat){
		s.statsRepo.Store(stat)
	}(stat)

	if !s.provider.ActionExists(name) {
		return nil, fmt.Errorf("action %v does not exists", name)
	}

	// TODO: Add Error handeling
	ac, _ := s.provider.BuildAction(name)

	env := ac.BuildEnv()
	payload := s.provider.EncodePayload(params)
	result := s.provider.InvokeAction(ac.Name, payload, env)

	return result, nil
}

func (s *service) AddSchedule(name string, schedID trigger.ScheduleID) error {
	sched, _ := s.scheduleRepo.Find(schedID)
	s.actionTimer.AddSchedule(sched, s.Run)

	return nil
}

func (s *service) RemoveSchedule(name string, schedID trigger.ScheduleID) error {
	sched, _ := s.scheduleRepo.Find(schedID)
	if !sched.Enabled {
		s.actionTimer.RemoveSchedule(sched.ID)
	}

	return nil
}


func (s *service) TriggerActionWithWebhook(action string, webhookID trigger.WebhookID, params map[string]interface{}) error {
	return nil
}




