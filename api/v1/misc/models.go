package misc

import "time"

type Token struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}
