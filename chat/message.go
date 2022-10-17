package chat

import "time"

// Message 메시지
type Message struct {
	Text     string    `json:"text"`
	Sender   string    `json:"sender"`
	SenderID int       `json:"senderID"`
	SentOn   time.Time `json:"sentOn"`
}

// Room 단톡방
type Room struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
