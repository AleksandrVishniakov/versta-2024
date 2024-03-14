package validator

type Validator interface {
	Valid() (bool, error)
}
