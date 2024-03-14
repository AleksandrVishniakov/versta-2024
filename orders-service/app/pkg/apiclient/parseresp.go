package apiclient

import (
	"github.com/AleksandrVishniakov/versta-2024/orders-service/app/pkg/parser"
	"net/http"

	"github.com/AleksandrVishniakov/versta-2024/orders-service/app/pkg/e"
)

func ScanResponse[T any](resp *http.Response, dest *T) error {
	if resp.StatusCode < 400 {
		if dest == nil {
			return nil
		}

		obj, err := parser.Decode[T](resp.Body)
		if err != nil {
			return err
		}

		*dest = obj

		return nil
	}

	apiError, err := parser.Decode[e.ResponseError](resp.Body)
	if err != nil {
		return err
	}

	return &apiError
}
