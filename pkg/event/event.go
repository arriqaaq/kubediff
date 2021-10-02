package event

type Event struct {
	Type string
	Kind string
	Obj  interface{}
	Diff interface{}
}

func NewEvent(event string, kind string, obj interface{}, diff interface{}) Event {
	return Event{event, kind, obj, diff}
}
