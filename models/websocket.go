package models

type WebSocketCommand struct {
	Cmd  string `json:"cmd"`
	ID   string `json:"id"`
	Room int    `json:"room"`
	Data int    `json:"data"`
}

type WebSocketReply struct {
	Cmd   string `json:"cmd"`
	Reply string `json:"reply"`
	Error string `json:"error,omitempty"`
}

type WebSocketEvent struct {
	Event    string          `json:"event"`
	Room     int             `json:"room"`
	Secret   int             `json:"secret"`
	Rankings []*PlayerResult `json:"rankings,omitempty"`
}
