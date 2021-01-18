package log

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

var (
	stackSourceFileName     = "file"
	stackSourceLineName     = "line"
	stackSourceFunctionName = "func"
	stackSourcePkgName      = "pkg"
)

// ErrStack error stack implmenet fmt.Stringer interface
type ErrStack []map[string]string

func (stack ErrStack) String() string {
	var s strings.Builder

	for i, m := range stack {
		if s.Len() > 0 && i != len(stack) {
			s.WriteString("-> ")
		}

		for _, key := range []string{"func", "file", "line"} {
			s.WriteString(stack.stackInfo(key, m[key]))
		}
	}
	return s.String()
}

func (stack ErrStack) stackInfo(k, v string) string {
	return fmt.Sprintf("[%s]:%v ", k, v)
}

// GetStack marshalling stack information
func GetStack(err error) ErrStack {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}
	sterr, ok := err.(stackTracer)
	if !ok {
		err = errors.WithStack(err)
		sterr = err.(stackTracer)
	}
	st := sterr.StackTrace()
	maxFrames := len(st)

	s := &state{}
	out := make([]map[string]string, 0, maxFrames)
	for _, frame := range st[:maxFrames] {
		out = append(out, map[string]string{
			stackSourceFileName:     frameField(frame, s, 's'),
			stackSourceLineName:     frameField(frame, s, 'd'),
			stackSourceFunctionName: frameField(frame, s, 'n'),
		})
	}
	return out
}

func frameField(f errors.Frame, s *state, c rune) string {
	f.Format(s, c)
	return string(s.b)
}

type state struct {
	b []byte
}

// Write implement fmt.Formatter interface.
func (s *state) Write(b []byte) (n int, err error) {
	s.b = b
	return len(b), nil
}

// Width implement fmt.Formatter interface.
func (s *state) Width() (wid int, ok bool) {
	return 0, false
}

// Precision implement fmt.Formatter interface.
func (s *state) Precision() (prec int, ok bool) {
	return 0, false
}

// Flag implement fmt.Formatter interface.
func (s *state) Flag(c int) bool {
	// if c == '+' {
	// 	return true
	// }
	return false
}
