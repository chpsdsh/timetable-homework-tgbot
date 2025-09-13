package urlselector

type Room struct {
	room    string
	fullURL string
}

func NewRoom(room string, fullURL string) Room {
	return Room{room: room, fullURL: fullURL}
}

func (r Room) GetShortName() string {
	return r.room
}

func (r Room) GetFullName() string {
	return r.fullURL
}
