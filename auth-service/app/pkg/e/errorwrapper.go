package e

import "fmt"

func WrapErr(err error, description string) error {
	return fmt.Errorf("%s: %w", description, err)
}

func WrapIfNotNil(err error, description string) error {
	if err != nil {
		return WrapErr(err, description)
	}

	return nil
}
