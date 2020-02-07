package validator

type InputFieldCode string

const (
	CodeMissing       InputFieldCode = "missing"
	CodeMissingField                 = "missing_field"
	CodeInvalid                      = "invalid"
	CodeAlreadyExists                = "already_exists"
)

type InputError struct {
	Message string
	Field   string
	Code    InputFieldCode
}

func (i *InputError) Error() string {
	return i.Message
}
