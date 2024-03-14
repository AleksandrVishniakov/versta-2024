package parser

import (
	"encoding/json"
	"github.com/AleksandrVishniakov/versta-2024/orders-service/app/pkg/validator"
	"io"
	"net/http"
)

func Encode[T any](writer io.Writer, obj T) error {
	err := json.NewEncoder(writer).Encode(obj)
	if err != nil {
		return err
	}

	return nil
}

func EncodeResponse[T any](writer http.ResponseWriter, obj T, statusCode int) error {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	return Encode(writer, obj)
}

func Decode[T any](reader io.Reader) (T, error) {
	var obj T

	err := json.NewDecoder(reader).Decode(&obj)
	if err != nil {
		return *new(T), err
	}

	return obj, err
}

func DecodeValid[T validator.Validator](reader io.Reader) (T, error) {
	obj, err := Decode[T](reader)
	if err != nil {
		return *new(T), err
	}

	ok, err := obj.Valid()
	if !ok {
		return *new(T), err
	}

	return obj, err
}
