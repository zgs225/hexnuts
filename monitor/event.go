package monitor

import (
	"fmt"
)

const (
	Events_ADD int8 = iota
	Events_Update
	Events_DEL
)

type Event struct {
	T int8
	K string
	V string
}

func (e *Event) String() string {
	return fmt.Sprintf("%d %s %s\n", e.T, e.K, e.V)
}
