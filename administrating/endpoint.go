package administrating

import (
	"context"
	"github.com/stickysh/sticky/action"
	"github.com/stickysh/sticky/pkg/endpoint"
	"github.com/stickysh/sticky/trigger"
	"time"
)

type scheduleRequest struct {
	Action    string
	Start     time.Time
	Recurring bool
	Interval  int
	Params    map[string]interface{}
}

type scheduleResponse struct {
	ID        string    `json:"id"`
	Action    string    `json:"action,omitempty"`
	Start     time.Time `json:"start,omitempty"`
	Recurring bool      `json:"recurring"`
	Interval  int       `json:"interval,omitempty"`
	Enabled   bool      `json:"enabled"`
}

func makeCreateScheduleEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		schedReq := req.(scheduleRequest)
		result, err := s.ScheduleRecurringAction(schedReq.Action, schedReq.Interval, action.Params(schedReq.Params))
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}

type scheduleToggleReq struct {
	ID        string
	Enabled   bool
}

func makeToggleScheduleEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		schedReq := req.(scheduleToggleReq)
		err := s.ToggleSchedule(trigger.ScheduleID(schedReq.ID), schedReq.Enabled)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}

type actionCreateReq struct {
	Name        string
}

func makeActionCreateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		createReq := req.(actionCreateReq)
		err := s.CreateAction(createReq.Name)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}

func makeListActionsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		name := req.(string)

		result, err := s.ListSchedule(name)
		if err != nil {
			return nil, err
		}

		list := make([]scheduleResponse, 0, len(result))
		for _, val := range result {
			list = append(list, scheduleResponse{
				string(val.ID),
				val.Action,
				val.Start,
				val.Recurring,
				val.Interval,
				val.Enabled,
			})
		}

		return listSchedule{Schedules: list}, nil
	}
}

type webhookCreateRequest struct {
	Action string
}

func makeCreateWebhookEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		createReq := req.(webhookCreateRequest)
		_, err := s.AddWebhook(createReq.Action)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

}

type listSchedule struct {
	Schedules []scheduleResponse `json:"schedules,omitempty"`
}

func makeListScheduleEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		name := req.(string)

		result, err := s.ListSchedule(name)
		if err != nil {
			return nil, err
		}

		list := make([]scheduleResponse, 0, len(result))
		for _, val := range result {
			list = append(list, scheduleResponse{
				string(val.ID),
				val.Action,
				val.Start,
				val.Recurring,
				val.Interval,
				val.Enabled,
			})
		}

		return listSchedule{Schedules: list}, nil
	}
}

type Stat struct {
	ID         string
	ActionName string
	Start      time.Time
	End        time.Time
	Status     bool
}

type statResponse struct {
	ID       string    `json:"id"`
	Action   string    `json:"action,omitempty"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end,omitempty"`
	Duration float32   `json:"duration"`
	Status   string      `json:"status,omitempty"`
}

type listStats struct {
	Stats []statResponse `json:"schedules,omitempty"`
}

func makeListStatsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		name := req.(string)

		result, err := s.ListStats(name)
		if err != nil {
			return nil, err
		}

		list := make([]statResponse, 0, len(result))
		for _, val := range result {
			list = append(list, statResponse{
				string(val.ID),
				val.Action,
				val.Start,
				val.End,
				float32(float64(val.End.Sub(val.Start)) / float64(time.Second)),
				string(val.Status),
			})
		}

		return listStats{Stats: list}, nil
	}
}
