package interfaces

type Validator interface {
	Valid() (problems map[string]string)
}
