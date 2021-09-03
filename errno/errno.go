package errno

import (
	"fmt"
)

type MicroServiceError struct {
	Code    int
	Message string
}

func (k *MicroServiceError) Error() string {
	return fmt.Sprintf("koala error, code:%d message:%v", k.Code, k.Message)
}

var (
	NotHaveInstance = &MicroServiceError{
		Code:    1,
		Message: "not have instance",
	}
	ConnFailed = &MicroServiceError{
		Code:    2,
		Message: "connect failed",
	}
	InvalidNode = &MicroServiceError{
		Code:    3,
		Message: "invalid node",
	}
	AllNodeFailed = &MicroServiceError{
		Code:    4,
		Message: "all node failed",
	}
)

func IsConnectError(err error) bool {

	koalaErr, ok := err.(*MicroServiceError)
	if !ok {
		return false
	}
	var result bool
	if koalaErr == ConnFailed {
		result = true
	}
	return result
}
