package invoking

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/stickysh/sticky/pkg/endpoint"
	ctrlhttp "github.com/stickysh/sticky/pkg/http"
)


type invokeRequest struct {
	Action string
	Params map[string]interface{}
}

type webhookRequest struct {
	ID string
	Action string
	Params map[string]interface{}
	Headers map[string]string
}

func MakeHandler(srv Service) http.Handler {
	r := httprouter.New()

	invokingHandler := ctrlhttp.NewServer(
		makeInvokeActionEndpoint(srv),
		decodeInvokeRequest,
		encodeInvokeResponse,
	)

	// TODO: Add security token
	r.Handler(http.MethodPost, "/invoking/v1/actions/:name", invokingHandler)


	return r
}


func decodeInvokeRequest(ctx context.Context, r *http.Request) (interface{}, error){

	var actionParams map[string]interface{}

	params := httprouter.ParamsFromContext(ctx)
	name := params.ByName("name")

	if err := json.NewDecoder(r.Body).Decode(&actionParams); err != nil {
		return nil, err
	}

	return &invokeRequest{
		Action: name,
		Params: actionParams,
	}, nil
}

func encodeInvokeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	return json.NewEncoder(w).Encode(resp)
}

func makeInvokeActionEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {


		payload := req.(invokeRequest)
		result, err := s.Run(payload.Action, payload.Params)
		if err != nil {

		}
		return result, nil
	}
}

