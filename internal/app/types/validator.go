package types

type Validator interface {
	Valid() (problems map[string]string)
}
