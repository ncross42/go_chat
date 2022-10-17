package ws

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/learningspoons-go/groupchat/chat"
	"github.com/learningspoons-go/groupchat/db"
)

const (
	writeTimeout   = 10 * time.Second
	readTimeout    = 60 * time.Second
	pingPeriod     = 10 * time.Second
	maxMessageSize = 512
)

type conn struct {
	wsConn     *websocket.Conn
	wg         sync.WaitGroup
	sub        db.ChatroomSubscription
	chatroomID int
	senderID   int
}

func newConn(wsConn *websocket.Conn, chatroomID, senderID int) *conn {
	return &conn{
		wsConn:     wsConn,
		chatroomID: chatroomID,
		senderID:   senderID,
	}
}

func (c *conn) run() error {
	sub, err := db.NewChatroomSubscription(c.chatroomID)
	if err != nil {
		return err
	}
	c.sub = sub

	c.wg.Add(2)
	go c.readPump()
	go c.writePump()

	c.wg.Wait()
	return nil
}

func (c *conn) readPump() {
	defer c.wg.Done()
	defer c.sub.Close()

	c.wsConn.SetReadLimit(maxMessageSize)
	c.wsConn.SetReadDeadline(time.Now().Add(readTimeout))
	c.wsConn.SetPongHandler(func(string) error {
		c.wsConn.SetReadDeadline(time.Now().Add(readTimeout))
		return nil
	})

	for {
		var msg chat.Message
		if err := c.wsConn.ReadJSON(&msg); err != nil {
			log.Println("err reading:", err)
			return
		}

		db.SendMessage(c.senderID, c.chatroomID, msg.Text)
	}
}

func (c *conn) writePump() {
	defer c.wg.Done()

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case s, more := <-c.sub.C:
			if !more {
				return
			}
			c.wsConn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err := c.wsConn.WriteJSON(s); err != nil {
				log.Println("err writing:", err)
				return
			}
		case <-ticker.C:
			c.wsConn.WriteControl(
				websocket.PingMessage, nil, time.Now().Add(writeTimeout))
		}
	}
}
