package administrating

import (
	"fmt"
	"github.com/stickysh/sticky/packing"
	"time"

	"github.com/stickysh/sticky/action"
	"github.com/stickysh/sticky/trigger"
)

type EventHandler interface {
	ScheduleWasAdded(event trigger.SchedulingEvent)
	ScheduleWasDisabled(event trigger.SchedulingEvent)
}

type HandlerFunc func (name string, schedID trigger.ScheduleID) error

type Service interface {
	ListActions() ([]*action.Action, error)
	CreateAction(name string) error
	ActionCodeModified(name string) error

	ListSchedule(name string) ([]*trigger.Schedule, error)
	ScheduleAction(name string, start time.Time) (trigger.ScheduleID, error)
	ScheduleRecurringAction(name string, interval int, params action.Params) (trigger.ScheduleID, error)
	ToggleSchedule(id trigger.ScheduleID, enabled bool) error


	AddWebhook(name string) (*trigger.Webhook, error)

	ListStats(name string) ([]*action.Stat, error)
}

type service struct {
	actionRepo   action.ActionRepo
	schedRepo    trigger.ScheduleRepo
	webhookRepo  trigger.WebhookRepo
	statsRepo    action.StatsRepo
	actionPacker *packing.ActionPack
	schedHandler EventHandler
}

func NewService(actRepo action.ActionRepo, schedRepo trigger.ScheduleRepo, statsRepo action.StatsRepo, actionPacker *packing.ActionPack, schedHandler EventHandler) *service {
	return &service{
		actionRepo: actRepo,
		schedRepo:  schedRepo,
		statsRepo: statsRepo,
		actionPacker: actionPacker,
		schedHandler: schedHandler,
	}
}

func (s *service) CreateAction(name string) error {
	return s.actionPacker.Create(name)
}

func (s *service) ActionCodeModified(name string) error {
	return s.actionPacker.Deploy(name)
}

func (s *service) ListActions() ([]*action.Action, error) {
	panic("implement me")
	return nil, nil
}

func (s *service) AddWebhook(action string) (*trigger.Webhook, error) {
	wh := trigger.NewWebhook(action)
	s.webhookRepo.Store(wh)
	return nil, nil
}


func (s *service) ListSchedule(name string) ([]*trigger.Schedule, error) {

	return s.schedRepo.FindAllByAction(name), nil
}

func (s *service) ListStats(name string) ([]*action.Stat, error) {

	return s.statsRepo.FindByAction(name), nil
}

func (s *service) ScheduleAction(name string, start time.Time) (trigger.ScheduleID, error) {
	panic("implement me")
}

func (s *service) ScheduleRecurringAction(name string, interval int, params action.Params) (trigger.ScheduleID, error) {
	sched := trigger.NewRecurringSchedule(name, time.Now(), interval, params)
	err := s.schedRepo.Store(sched)
	if err != nil {
		return "", fmt.Errorf("could not save")
	}

	s.schedHandler.ScheduleWasAdded(trigger.SchedulingEvent{
		Action: name,
		ScheduleID: sched.ID,
	})

	return sched.ID, nil
}

func (s *service) ToggleSchedule(id trigger.ScheduleID, enabled bool) error {
	sched, err := s.schedRepo.Find(id)
	if err != nil {
		// TODO add error for missing item
		return err
	}

	sched.Enabled = !sched.Enabled

	err = s.schedRepo.Store(sched)
	if err != nil {
		return err
	}

	if sched.Enabled {

	}

	s.schedHandler.ScheduleWasDisabled(trigger.SchedulingEvent{
		Action:     sched.Action,
		ScheduleID: sched.ID,
	})

	return nil
}

type schedulingEventHandler struct {
	ScheduleAddHandlerFunc func (name string, schedID trigger.ScheduleID) error
	ScheduleDisableHandlerFunc func (name string, schedID trigger.ScheduleID) error
}

func (h *schedulingEventHandler) ScheduleWasAdded(event trigger.SchedulingEvent) {
	h.ScheduleAddHandlerFunc(event.Action, event.ScheduleID)
}

func (h *schedulingEventHandler) ScheduleWasDisabled(event trigger.SchedulingEvent) {
	h.ScheduleDisableHandlerFunc(event.Action, event.ScheduleID)
}


func NewEventHandler(addHandlerFunc, disableHandlerFunc HandlerFunc) EventHandler {
	return &schedulingEventHandler{
		ScheduleAddHandlerFunc: addHandlerFunc,
		ScheduleDisableHandlerFunc: disableHandlerFunc,
	}
}

