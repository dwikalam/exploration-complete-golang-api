package defaultcodec

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dwikalam/ecommerce-service/internal/app/helperdependency/response"
	"github.com/dwikalam/ecommerce-service/internal/app/type/icommon"
	"github.com/dwikalam/ecommerce-service/internal/app/type/wrappertype"
)

func Encode[T any](
	w http.ResponseWriter,
	statusCode int,
	message string,
	data T,
) error {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")

	response := response.Default[T]{
		Message: message,
		Data:    data,
	}

	return json.NewEncoder(w).Encode(response)
}

func Decode[T any](r *http.Request) (T, error) {
	var v T

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("error decode json: %w", err)
	}

	return v, nil
}

func DecodeValid[T icommon.Validator](r *http.Request) (payload T, problems wrappertype.ProblemsMap, err error) {
	if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return payload, nil, fmt.Errorf("error decode valid json: %w", err)
	}

	if problems = payload.Valid(r.Context()); len(problems) > 0 {
		return payload, problems, fmt.Errorf("error %T: %d problems", payload, len(problems))
	}

	return payload, nil, nil
}
