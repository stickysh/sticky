package invoking

import (
	"github.com/stickysh/sticky/invoke"
	"time"

	"github.com/stickysh/sticky/action"
	"github.com/stickysh/sticky/trigger"
)

type Service interface {
	RunAction(name string, params map[string]interface{}) (interface{}, error)
	AddActionSchedule(name string, schedID trigger.ScheduleID) error
	RemoveActionSchedule(name string, schedID trigger.ScheduleID) error


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


func (s *service) RunAction(name string, params map[string]interface{}) (interface{}, error) {
	stat := action.NewStat(name, time.Now(), action.Running)
	defer func(stat *action.Stat){
		s.statsRepo.Store(stat)
	}(stat)

	// TODO: Add Error handeling
	ac, _ := s.provider.BuildAction(name)

	env := ac.BuildEnv()
	payload := s.provider.EncodePayload(params)
	result := s.provider.InvokeAction(name, payload, env)

	return result, nil
}

func (s *service) AddActionSchedule(name string, schedID trigger.ScheduleID) error {
	sched, _ := s.scheduleRepo.Find(schedID)
	s.actionTimer.AddSchedule(sched, s.RunAction)

	return nil
}

func (s *service) RemoveActionSchedule(name string, schedID trigger.ScheduleID) error {
	sched, _ := s.scheduleRepo.Find(schedID)
	if !sched.Enabled {
		s.actionTimer.RemoveSchedule(sched.ID)
	}

	return nil
}


func (s *service) TriggerActionWithWebhook(action string, webhookID trigger.WebhookID, params map[string]interface{}) error {
	return nil
}




