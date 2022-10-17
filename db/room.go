package db

import "github.com/learningspoons-go/groupchat/chat"

// CreateRoom 방 만들기
func CreateRoom(name string) (id int, err error) {
	err = db.QueryRow(`INSERT INTO chatrooms (name)
		VALUES ($1)
		RETURNING id`, name).Scan(&id)
	return
}

// RoomExists 방 존재 여부
func RoomExists(id int) (exists bool, err error) {
	err = db.QueryRow(`SELECT EXISTS(
			SELECT 1 FROM chatrooms WHERE id = $1)`, id).
		Scan(&exists)
	return
}

// GetRooms 모든 채팅방 리턴
func GetRooms() ([]chat.Room, error) {
	rooms := []chat.Room{}
	rows, err := db.Query(`SELECT id, name FROM chatrooms`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var room chat.Room
		if err := rows.Scan(&room.ID, &room.Name); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}
