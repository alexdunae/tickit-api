package main

import (
	"fmt"
)

type Config struct {
	Hosts struct {
		AllowedOrigins string
	}

	Servers struct {
		Http       string
		WebSockets string
	}

	Database struct {
		Database string
		Username string
		Password string
		Host     string
	}
}

func (c Config) DSN() string {
	return fmt.Sprintf("%s/%s/%s", c.Database.Database, c.Database.Username, c.Database.Password)
}
