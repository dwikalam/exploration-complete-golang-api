package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dwikalam/ecommerce-service/internal/app/types/dto"
	"github.com/dwikalam/ecommerce-service/internal/app/types/interfaces"
)

func Encode[T any](
	w http.ResponseWriter,
	statusCode int,
	message string,
	data T,
) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := dto.Response[T]{
		Message: message,
		Data:    data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return fmt.Errorf("error encode json: %w", err)
	}

	return nil
}

func Decode[T any](r *http.Request) (T, error) {
	var v T

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("error decode json: %w", err)
	}

	return v, nil
}

func DecodeValid[T interfaces.Validator](r *http.Request) (T, map[string]string, error) {
	var v T

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, nil, fmt.Errorf("error decode valid json: %w", err)
	}

	if problems := v.Valid(); len(problems) > 0 {
		return v, problems, fmt.Errorf("error %T: %d problems", v, len(problems))
	}

	return v, nil, nil
}
