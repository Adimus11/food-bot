package objects

const MessageEventType = "message"

type MessageEvent struct {
	Message string `json:"message"`
}

const CardEventType = "card"

type CardEvent struct {
	DishID      string `json:"dish_id"`
	Title       string `json:"title_id"`
	Description string `json:"description"`
	Image       string `json:"img"`
	Link        string `json:"link"`
}

const RatingRequestedEventType = "rating_requested"
const RatingSetEventType = "rating_set"

type RatingEvent struct {
	DishID string `json:"dish_id"`
	Rating *int   `json:"rating,omitempty"`
}

const ChatIdleEventType = "chat_idle"

type ChatIdleEvent struct {
	Reason string `json:"reason"`
}

var EventTypes = []string{
	MessageEventType,
	CardEventType,
	RatingRequestedEventType,
	RatingSetEventType,
	ChatIdleEventType,
	SelectEventType,
}

type StatusResponse struct {
	Status string `json:"status"`
}

const SelectEventType = "select"

type DishSelection struct {
	Message          string    `json:"message,omitempty"`
	Options          []*Option `json:"options,omitempty"`
	SelectedOptionID string    `json:"selected_option_id"`
}

type Option struct {
	OptionID   int    `json:"option_id"`
	OptionText string `json:"option_text"`
	DishID     string `json:"dish_id"`
}
