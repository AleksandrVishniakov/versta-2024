package e

import "fmt"

func WrapErr(err error, description string) error {
	return fmt.Errorf("%s: %w", description, err)
}

func WrapErrWithErr(err error, description error) error {
	return fmt.Errorf("%w: %w", description, err)
}

func WrapIfNotNil(err error, description string) error {
	if err != nil {
		return WrapErr(err, description)
	}

	return nil
}

func WrapWithErrIfNotNil(err error, description error) error {
	if err != nil {
		return WrapErrWithErr(err, description)
	}

	return nil
}
