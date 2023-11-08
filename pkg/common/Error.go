package common

import "fmt"

var _ error = (*gostError)(nil)

type gostError struct {
	code int
	msg  string
}

func (e gostError) Error() string {
	return fmt.Sprintf("code:%d,msg:%v", e.code, e.msg)
}

func NewError(code int, msg string) error {
	return gostError{
		code: code,
		msg:  msg,
	}
}

func GetErrorCode(err error) int {
	if e, ok := err.(gostError); ok {
		return e.code
	} else {
		return ERRORCODE_WRONGTYPE
	}

}

func GetErrorMsg(err error) string {
	if e, ok := err.(gostError); ok {
		return e.msg
	} else {
		return ""
	}
}
