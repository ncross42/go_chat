package db

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/learningspoons-go/groupchat/chat"
)

// ChatroomSubscription 채팅방 구독(연결)
type ChatroomSubscription struct {
	sub subscription
	C   <-chan chat.Message
}

func sendExistingMessages(chatroomID int, c chan<- chat.Message,
	limit int) error {
	rows, err := db.Query(`WITH msgs AS (
		SELECT m.sender_id, s.username, m.text, m.sent_on
			FROM messages m
			JOIN users s ON m.sender_id = s.id
			WHERE m.chatroom_id = $1
			ORDER BY m.sent_on DESC
			LIMIT $2)
		SELECT * FROM msgs ORDER BY sent_on ASC`, chatroomID, limit)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var msg chat.Message
		if err := rows.Scan(&msg.SenderID, &msg.Sender, &msg.Text,
			&msg.SentOn); err != nil {
			return err
		}
		c <- msg
	}
	return nil
}

// NewChatroomSubscription 채팅방에 접속(구독)하기
func NewChatroomSubscription(chatroomID int) (ChatroomSubscription, error) {
	c := make(chan chat.Message, 128)
	if err := sendExistingMessages(chatroomID, c, 100); err != nil {
		return ChatroomSubscription{}, err
	}

	chatroomSubscription := ChatroomSubscription{
		sub: subscribe(fmt.Sprintf("new_message_%d", chatroomID)),
		C:   c,
	}

	go func() {
		defer close(c)
		for msg := range chatroomSubscription.sub.c {
			var parsedMessage chat.Message
			if err := json.Unmarshal([]byte(msg), &parsedMessage); err != nil {
				log.Println("couldn't parse message", err)
				continue
			}
			c <- parsedMessage
		}
	}()

	return chatroomSubscription, nil
}

// Close 닫기
func (s ChatroomSubscription) Close() {
	s.sub.close()
}

// SendMessage 문자 보내기
func SendMessage(senderID, chatroomID int, text string) error {
	_, err := db.Exec(`INSERT INTO messages (sender_id, chatroom_id, text)
	 		VALUES ($1, $2, $3)`, senderID, chatroomID, text)
	return err
}
