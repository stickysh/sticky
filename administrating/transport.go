package administrating

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	ctrlhttp "github.com/stickysh/sticky/pkg/http"
)

func MakeHandler(srv Service) http.Handler {
	r := httprouter.New()

	createScheduleHandler := ctrlhttp.NewServer(
		makeCreateScheduleEndpoint(srv),
		decodeScheduleRequest,
		encodeResponse,
	)

	toggleScheduleHandler := ctrlhttp.NewServer(
		makeToggleScheduleEndpoint(srv),
		decodeToggleRequest,
		encodeResponse,
	)

	listScheduleHandler := ctrlhttp.NewServer(
		makeListScheduleEndpoint(srv),
		decodeActionName,
		encodeResponse,
	)

	listStatsHandler := ctrlhttp.NewServer(
		makeListStatsEndpoint(srv),
		decodeActionName,
		encodeResponse,
	)

	actionCreateHandler := ctrlhttp.NewServer(
		makeActionCreateEndpoint(srv),
		decodeActionCreateRequest,
		encodeResponse,
	)

	actionsListHandler := ctrlhttp.NewServer(
		makeListActionsEndpoint(srv),
		decodeActionCreateRequest,
		encodeResponse,
	)

	actionAddSecretHandler := ctrlhttp.NewServer(
		makeActionCreateEndpoint(srv),
		decodeActionName,
		encodeResponse,
	)

	createWebhookHandler := ctrlhttp.NewServer(
		makeCreateWebhookEndpoint(srv),
		decodeWebhookCreateRequest,
		encodeResponse,
	)

	// TODO: Add Auth
	// Actions
	r.Handler(http.MethodGet, "/admin/v1/actions", actionsListHandler)
	r.Handler(http.MethodPost, "/admin/v1/actions", actionCreateHandler)
	r.Handler(http.MethodPost, "/admin/v1/actions/:name/secrets", actionAddSecretHandler)
	r.Handler(http.MethodPost, "/admin/v1/actions/:name/schedules", listScheduleHandler)



	// List details regarding an action
	r.Handler(http.MethodGet, "/admin/v1/actions/:name", nil)

	// Webhook
	r.Handler(http.MethodPost, "/admin/v1/webhook", createWebhookHandler)

	// Scheduling
	r.Handler(http.MethodPost, "/admin/v1/schedule", createScheduleHandler)
	r.Handler(http.MethodPatch, "/admin/v1/schedule/:id", toggleScheduleHandler)


	// Stats
	r.Handler(http.MethodGet, "/admin/v1/stats/:name", listStatsHandler)

	return r
}

func decodeScheduleRequest(_ context.Context, r *http.Request) (interface{}, error) {

	var body struct {
		Action    string `json:"action"`
		Start     time.Time `json:"start"`
		Recurring bool `json:"recurring"`
		Interval  int	`json:"interval"`
		Params    map[string]interface{} `json:"params"`

	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return scheduleRequest{
		Action: body.Action,
		Start: body.Start,
		Recurring: body.Recurring,
		Interval: body.Interval,
		Params: body.Params,
	}, nil
}

func decodeToggleRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	params := httprouter.ParamsFromContext(ctx)
	id := params.ByName("id")
	var body struct {
		Enabled bool `json:"recurring"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return scheduleToggleReq{
		ID: id,
		Enabled: body.Enabled,
	}, nil
}

func decodeActionCreateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return actionCreateReq{
		Name: body.Name,
	}, nil
}

func decodeWebhookCreateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		Action    string `json:"action"`
		SignHeaderName string `json:"signHeaderName"`

	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return webhookCreateRequest{
		Action: body.Action,
	}, nil
}

func decodeActionName(ctx context.Context, r *http.Request) (interface{}, error) {
	params := httprouter.ParamsFromContext(ctx)
	name := params.ByName("name")

	return name, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	//if e, ok := response.(errorer); ok && e.error() != nil {
	//	encodeError(ctx, e.error(), w)
	//	return nil
	//}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}




