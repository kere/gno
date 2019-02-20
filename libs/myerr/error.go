package myerr

import (
	"fmt"

	"github.com/kere/gno/libs/log"
)

// M args map
type M map[string]interface{}

// Error class
type Error struct {
	args      M
	message   string
	IsStack   bool
	IsStacked bool
	IsLoged   bool
}

// Error string
func (e *Error) Error() string {
	return e.message
}

// Args map
func (e *Error) Args(args M) *Error {
	e.IsStacked, e.IsLoged = false, false
	e.args = args
	return e
}

// Stack error
func (e *Error) Stack() *Error {
	if !e.IsStack || e.IsStacked {
		return e
	}
	log.App.Stack()
	e.IsStacked = true
	return e
}

// Log error
func (e *Error) Log() *Error {
	if e.IsLoged {
		return e
	}
	log.App.Error(e.message)
	if len(e.args) > 0 {
		log.App.Error("args:", e.args)
	}
	e.IsLoged = true
	return e
}

// New error
func New(m ...interface{}) *Error {
	if len(m) == 1 {
		if err, isok := m[0].(*Error); isok {
			return err
		}
	}
	return &Error{IsStack: true, message: fmt.Sprint(m...)}
}

// NewVersion new
func NewVersion() *Error {
	return &Error{IsStack: false, message: "有人已经修改了数据，请刷新后再试"}
}

// // ExistsErr 已经存在
// type ExistsErr string
//
// func (v ExistsErr) Error() string {
// 	return string(v)
// }
//
// // NewExistsErr new
// func NewExistsErr(s string) ExistsErr {
// 	return ExistsErr(s)
// }
//
// // ForbidErr 禁止
// type ForbidErr string
//
// func (v ForbidErr) Error() string {
// 	return string(v)
// }
//
// // NewForbidErr new
// func NewForbidErr(s string) ForbidErr {
// 	return ForbidErr(s)
// }

// // NotFoundErr 错误
// type NotFoundErr string
//
// func (v NotFoundErr) Error() string {
// 	return string(v)
// }
//
// // NewNotFoundErr new
// func NewNotFoundErr(s string) NotFoundErr {
// 	return NotFoundErr(s)
// }
