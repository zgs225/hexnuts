package monitor

const (
	Events_ADD int8 = iota
	Events_DEL
)

type Event struct {
	T int8
	K string
	V string
}
