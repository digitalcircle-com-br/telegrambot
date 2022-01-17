package lib

import "time"

type Sub struct {
	ID        uint
	Chan      string
	ChatID    string
	Username  string
	FirstName string
	LastName  string
	JoinnedAt time.Time
}

type Session struct {
	ID     uint
	ChatID string
}

type PubMsg struct {
	Ch  string `json:"ch"`
	Msg string `json:"msg"`
}
