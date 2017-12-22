package model

import (
	"fmt"
)

type AuthUser struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	UserName string `yaml:"username"`
	PassWord string `yaml:"password"`
}

func (this AuthUser) String() string {
	return fmt.Sprintf(
		"Smtp Server: %s, Smtp server port: %d, UserName: %s, Password: %s",
		this.Host,
		this.Port,
		this.UserName,
		this.PassWord,
	)
}

func (this AuthUser) Islegal() bool {
	if this.Host == "" || this.Port < 0 || this.Port > 65535 || this.UserName == "" || this.PassWord == "" {
		return false
	}

	return true
}
