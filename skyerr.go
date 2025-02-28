package skyerr

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

type SkyErr interface {
	Code() int32
	Error() string
	Msg() string
	String() string
	AppendStack()
}

type throwStack struct {
	file     string
	line     int
	funcName string
}

type throwErr struct {
	code   int32
	err    string //go 错误
	msg    string //错误友好提示信息
	stacks []*throwStack
}

var (
	defaultErr = newThrow
)

const (
	errCode int32 = -1
)

func SetDefErr(defErr func(int32, error, ...any) SkyErr) {
	defaultErr = defErr
}

func newThrow(code int32, err error, opts ...any) SkyErr {
	var throw *throwErr
	if code == 0 {
		return throw
	}
	throw = &throwErr{code: code}
	if err != nil {
		throw.err = err.Error()
	}
	if len(opts) > 0 {
		throw.msg = fmt.Sprintf(opts[0].(string), opts[1:]...)
	}
	return throw
}

// Code return value of code,if throw is nil,return 0;
func (throw *throwErr) Code() int32 {
	if throw == nil {
		return 0
	}
	return throw.code
}

func (throw *throwErr) Error() string {
	if throw == nil {
		return ""
	}
	return throw.err
}

func (throw *throwErr) Msg() string {
	if throw == nil {
		return ""
	}
	return throw.msg
}

func (throw *throwErr) String() string {
	if throw == nil || throw.code == 0 {
		return "ok(0)"
	}
	stack := throw.stringStack()
	if stack == "" {
		return fmt.Sprintf("code:%d, msg:%s", throw.code, throw.msg)
	} else {
		return fmt.Sprintf("code:%d, msg:%s, err:%s", throw.code, throw.msg, stack)
	}
}

func (throw *throwErr) AppendStack() {
	if throw == nil {
		return
	}
	pc, file, line, ok := runtime.Caller(2)
	if ok {
		splits := strings.Split(file, "/")
		throw.stacks = append(throw.stacks, &throwStack{
			file:     splits[len(splits)-1],
			line:     line,
			funcName: runtime.FuncForPC(pc).Name(),
		})
	}
}

func (throw *throwErr) stringStack() string {
	if throw == nil {
		return ""
	}
	if len(throw.stacks) <= 0 {
		return throw.err
	}
	buf := &strings.Builder{}
	buf.WriteString(throw.err)
	buf.WriteString("\n")
	for _, stack := range throw.stacks {
		buf.WriteString("        at ")
		buf.WriteString(stack.funcName)
		buf.WriteString("(")
		buf.WriteString(stack.file)
		buf.WriteString(":")
		buf.WriteString(strconv.Itoa(stack.line))
		buf.WriteString(")\n")
	}
	return buf.String()
}

// SkyError only set error
func SkyError(err any) SkyErr {
	if err == nil {
		return defaultErr(0, nil)
	}
	var in SkyErr
	switch e := err.(type) {
	case SkyErr:
		in = e
	case error:
		in = defaultErr(errCode, e)
	default:
		in = defaultErr(errCode, nil, fmt.Sprint(e))
	}
	if in.Error() != "" {
		in.AppendStack()
	}
	return in
}

// SkyErrorM only set information
func SkyErrorM(code int32, opts ...any) SkyErr {
	return defaultErr(code, nil, opts...)
}

// SkyErrorF  with error
func SkyErrorF(code int32, err error, opts ...any) SkyErr {
	in := defaultErr(code, err, opts...)
	if in.Error() != "" {
		in.AppendStack()
	}
	return in
}
