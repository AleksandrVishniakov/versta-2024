package main

import (
	"encoding/json"
	"io"
)

func Encode[T any](writer io.Writer, obj T) error {
	err := json.NewEncoder(writer).Encode(obj)
	if err != nil {
		return err
	}

	return nil
}

func Decode[T any](reader io.Reader) (T, error) {
	var obj T

	err := json.NewDecoder(reader).Decode(&obj)
	if err != nil {
		return *new(T), err
	}

	return obj, err
}
