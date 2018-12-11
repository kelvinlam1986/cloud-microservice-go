package persistence

type DatabaseHandler interface {
	AddUser(user User) ([]byte, error)
	AddEvent(event Event) ([]byte, error)
	AddBookingForUser([]byte, Booking) error
	AddLocation(Location) (Location, error)
	FindUser(string, string) (User, error)
	FindBookingForUser([]byte) ([]Booking, error)
	FindEvent([]byte) (Event, error)
	FindEventByName(name string) (Event, error)
	FindAllAvailableEvents() ([]Event, error)
	FindLocation(string) (Location, error)
	FindAllLocation() ([]Location, error)
}
