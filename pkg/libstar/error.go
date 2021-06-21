package libstar

import "fmt"

type Err struct {
	Code    int
	Message string
}

func NewErr(message string, v ...interface{}) (e *Err) {
	e = &Err{
		Code:    0xFFff,
		Message: fmt.Sprintf(message, v...),
	}
	return
}

func (e *Err) String() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func (e *Err) Error() string {
	return e.String()
}
