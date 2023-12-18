package common

import "fmt"

var _ error = (*gostError)(nil)

type gostError struct {
	code   int
	reason int
	msg    string
}

func (e gostError) Error() string {
	return fmt.Sprintf("code:%d,msg:%v", e.code, e.msg)
}

func NewError(code int, reason int, msg string) error {
	return gostError{
		code:   code,
		reason: reason,
		msg:    msg,
	}
}

func GetErrorCode(err error) (int, int) {
	if e, ok := err.(gostError); ok {
		return e.code, e.reason
	} else {
		return ERRORCODE_WRONGTYPE, REASON_UNKNOWN
	}
}

func GetErrorMsg(err error) string {
	if e, ok := err.(gostError); ok {
		return e.msg
	} else {
		return ""
	}
}
