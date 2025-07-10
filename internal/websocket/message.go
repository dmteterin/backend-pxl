package websocket

import "time"

const (
	MessageTypeTextUpdate  = "text_update"
	MessageTypeUserJoined  = "user_joined"
	MessageTypeUserLeft    = "user_left"
	MessageTypeHeartbeat   = "heartbeat"
	MessageTypeUserList    = "user_list"
	MessageTypePixelUpdate = "pixel_update"
	MessageTypeCanvasState = "canvas_state"
	MessageTypeCanvasClear = "canvas_clear"
)

type Message struct {
	Type      string    `json:"type"`
	Content   string    `json:"content,omitempty"`
	ClientID  string    `json:"client_id,omitempty"`
	UserList  []string  `json:"user_list,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	// Pixel art fields
	X      int               `json:"x"`
	Y      int               `json:"y"`
	Color  string            `json:"color,omitempty"`
	Pixels map[string]string `json:"pixels,omitempty"`
}

type TextUpdateMessage struct {
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

func NewTextUpdateMessage(content string) *Message {
	return &Message{
		Type:      MessageTypeTextUpdate,
		Content:   content,
		Timestamp: time.Now(),
	}
}

func NewUserJoinedMessage(clientID string) *Message {
	return &Message{
		Type:      MessageTypeUserJoined,
		ClientID:  clientID,
		Timestamp: time.Now(),
	}
}

func NewUserLeftMessage(clientID string) *Message {
	return &Message{
		Type:      MessageTypeUserLeft,
		ClientID:  clientID,
		Timestamp: time.Now(),
	}
}

func NewUserListMessage(userList []string) *Message {
	return &Message{
		Type:      MessageTypeUserList,
		UserList:  userList,
		Timestamp: time.Now(),
	}
}

func NewPixelUpdateMessage(x, y int, color, clientID string) *Message {
	return &Message{
		Type:      MessageTypePixelUpdate,
		X:         x,
		Y:         y,
		Color:     color,
		ClientID:  clientID,
		Timestamp: time.Now(),
	}
}

func NewCanvasStateMessage(pixels map[string]string) *Message {
	return &Message{
		Type:      MessageTypeCanvasState,
		Pixels:    pixels,
		Timestamp: time.Now(),
	}
}

func NewCanvasClearMessage() *Message {
	return &Message{
		Type:      MessageTypeCanvasClear,
		Timestamp: time.Now(),
	}
}
