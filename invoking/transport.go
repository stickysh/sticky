package invoking

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/stickysh/sticky/pkg/endpoint"
	ctrlhttp "github.com/stickysh/sticky/pkg/http"
)


type invokeRequest map[string]interface{}



func MakeHandler(srv Service) http.Handler {
	r := httprouter.New()

	invokingHandler := ctrlhttp.NewServer(
		makeInvokeActionEndpoint(srv),
		decodeInvokeRequest,
		encodeInvokeResponse,
	)

	webhookHandler := ctrlhttp.NewServer(
		makeWebhookActionEndpoint(srv),
		decodeInvokeRequest,
		encodeInvokeResponse,
	)

	// TODO: Add security token
	r.Handler(http.MethodPost, "/invoking/v1/actions/:name", invokingHandler)

	r.Handler(http.MethodPost, "/invoking/v1/webhook/:name/:id", webhookHandler)

	return r
}


func decodeInvokeRequest(_ context.Context, r *http.Request) (interface{}, error){
	var params invokeRequest

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return nil, err
	}


	return params, nil
}

func encodeInvokeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	return json.NewEncoder(w).Encode(resp)
}

func makeInvokeActionEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		params := httprouter.ParamsFromContext(ctx)
		name := params.ByName("name")

		payload := req.(invokeRequest)
		result, err := s.RunAction(name, payload)
		if err != nil {

		}
		return result, nil
	}
}

func makeWebhookActionEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		params := httprouter.ParamsFromContext(ctx)
		name := params.ByName("name")

		payload := req.(invokeRequest)
		result, err := s.RunAction(name, payload)
		if err != nil {

		}
		return result, nil
	}
}
