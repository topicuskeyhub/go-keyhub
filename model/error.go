package model

import "fmt"

// ErrorReport error report returned by KeyHub Api
type ErrorReport struct {
	Code             int      `json:"code"`
	Reason           string   `json:"reason"`
	Exception        string   `json:"exception,omitempty"`
	Message          string   `json:"message"`
	ApplicationError string   `json:"applicationError,omitempty"`
	StackTrace       []string `json:"stacktrace,omitempty"`
}

func (er ErrorReport) Error() string {
	return er.Message
}

// Wrap  rap the errorReport within an error of type KeyhubApiError
func (er ErrorReport) Wrap(format string, any ...any) error {
	return &KeyhubApiError{
		Message: fmt.Sprintf(format, any...),
		Report:  ErrorReport{},
	}
}

// NewKeyhubApiError Create new KeyhubApiError from ErrorReport
func NewKeyhubApiError(errorReport ErrorReport, format string, any ...any) *KeyhubApiError {

	err := KeyhubApiError{
		Message: fmt.Sprintf(format, any...),
		Report:  errorReport,
	}
	return &err
}

type KeyhubApiError struct {
	Message string
	Report  ErrorReport
}

func (e KeyhubApiError) Error() string {
	return fmt.Sprintf("%s Error: %s", e.Message, e.Report.Message)
}
