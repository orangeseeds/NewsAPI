package storage

import (
	"fmt"
)

type errType int

const (
	None              errType = 0
	NeoError          errType = 1
	InvalidConstraint errType = 2
	NoMatchFound      errType = 3
)

type DBError struct {
	Type errType
	Msg  string
}

func (e *DBError) Error() string {
	return fmt.Sprintf("type: %d, %s", e.Type, e.Msg)
}

func MatchErr[T error](v error, m T) bool {
	switch any(v).(type) {
	case T:
		return true
	default:
		return false
	}
}

func MatchDBErr(v error) bool {
	return MatchErr(v, &DBError{})
}

func GetType(v error) errType {
	if MatchDBErr(v) {
		e := v.(*DBError)
		return e.Type
	}
	return None
}

func NewError(t errType, msg string) *DBError {
	return &DBError{
		Type: t,
		Msg:  msg,
	}
}
