package status

import (
	"fmt"

	"google.golang.org/grpc/codes"
)

func Error(code codes.Code, msg string) error {
	return fmt.Errorf("%s: %s", code.String(), msg)
}

func Errorf(code codes.Code, format string, args ...interface{}) error {
	return Error(code, fmt.Sprintf(format, args...))
}
