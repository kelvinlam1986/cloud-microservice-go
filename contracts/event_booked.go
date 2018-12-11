package contracts

type EventBookedEvent struct {
	EventId string `json:"eventId"`
	UserId string `json:"userId"`
}

func (e *EventBookedEvent) EventName() string {
	return "eventBooked"
}
