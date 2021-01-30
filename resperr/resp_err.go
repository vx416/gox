package resperr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func NewRespErr(status int, msg ...string) error {
	errResp := &ErrResponse{
		HTTPStatus: status,
		Details:    make(map[string]interface{}),
	}

	if len(msg) == 1 {
		errResp.Message = msg[0]
	} else {
		errResp.Message = http.StatusText(status)
	}

	return errors.WithStack(errResp)
}

// ErrResponse error response
type ErrResponse struct {
	HTTPStatus int
	Message    string
	Details    map[string]interface{}
	internal   error
}

func (e ErrResponse) Error() string {
	var b strings.Builder
	_, _ = b.WriteRune('[')
	_, _ = b.WriteString(strconv.Itoa(e.HTTPStatus))
	_, _ = b.WriteRune(']')
	_, _ = b.WriteRune(' ')
	_, _ = b.WriteString(e.Message)
	if e.internal != nil {
		_, _ = b.WriteRune('(')
		_, _ = b.WriteString(e.internal.Error())
		_, _ = b.WriteRune(')')
	}
	return b.String()
}

func (e ErrResponse) MarshalJSON() ([]byte, error) {
	data := make(map[string]map[string]interface{})
	data["error"] = map[string]interface{}{
		"code":    e.HTTPStatus,
		"message": e.Message,
		"details": e.Details,
	}

	return json.Marshal(data)
}

func (e ErrResponse) Unwrap() error {
	return e.internal
}

func (e *ErrResponse) AddDetail(key string, val interface{}) {
	e.Details[key] = val
}

func ToErrResponse(err error) *ErrResponse {
	cause := errors.Cause(err)
	errResp, ok := cause.(*ErrResponse)
	if !ok {
		return &ErrResponse{
			HTTPStatus: http.StatusInternalServerError,
			Message:    http.StatusText(http.StatusInternalServerError),
			internal:   err,
			Details:    make(map[string]interface{}),
		}
	}

	return errResp
}

func WithDetails(err error, details map[string]interface{}) error {
	errResp := ToErrResponse(err)
	for k, v := range details {
		errResp.Details[k] = v
	}
	return errResp
}

func WithNewMsgf(err error, format string, args ...interface{}) error {
	errResp := ToErrResponse(err)
	errResp.Message = fmt.Sprintf(format, args...)
	return errResp
}
