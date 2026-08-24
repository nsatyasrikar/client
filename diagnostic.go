package client

import "fmt"

type Severity uint8

const (
	SeverityWarning Severity = iota + 1
	SeverityError
)

type Diagnostic struct {
	Code, Message, Path string
	Severity            Severity
}
type Diagnostics []Diagnostic

func (d Diagnostics) Error() string {
	if len(d) == 0 {
		return ""
	}
	return fmt.Sprintf("%s: %s", d[0].Code, d[0].Message)
}
