package model

import (
	"fmt"
)

type Mail struct {
	Tos     string `json:"tos" form:"tos"`
	Subject string `json:"subject" form:"subject"`
	Content string `json:"content" form:"content"`
}

func (this Mail) String() string {
	return fmt.Sprintf("tos: %s, subject: %s, content: %s",
		this.Tos,
		this.Subject,
		this.Content,
	)
}
