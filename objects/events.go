package objects

// Chat events

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

const RatingEventType = "rating_event"

type RatingEvent struct {
	DishID string `json:"dish_id"`
	Rating *int   `json:"rating,omitempty"`
}

const ChatIdleEventType = "chat_idle"

type ChatIdleEvent struct {
	Reason string `json:"reason"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

const SelectEventType = "select"

type DishSelection struct {
	Options          []*Option `json:"options,omitempty"`
	SelectedOptionID *int      `json:"selected_option_id,omitempty"`
}

func (ds *DishSelection) GetDishIDForOption() string {
	if ds.SelectedOptionID == nil {
		return ""
	}

	for _, option := range ds.Options {
		if option.OptionID == *ds.SelectedOptionID {
			return option.DishID
		}
	}

	return ""
}

type Option struct {
	OptionID   int    `json:"option_id"`
	OptionText string `json:"option_text"`
	DishID     string `json:"dish_id"`
}

var EventTypes = []string{
	MessageEventType,
	CardEventType,
	RatingEventType,
	ChatIdleEventType,
	SelectEventType,
}
