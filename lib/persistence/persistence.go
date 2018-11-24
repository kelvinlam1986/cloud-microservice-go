package persistence

type DatabaseHandler interface {
	AddEvent(event Event) ([]byte, error)
	FindEvent([]byte) (Event, error)
	FindEventByName(name string) (Event, error)
	FindAllAvailableEvents() ([]Event, error)
}
