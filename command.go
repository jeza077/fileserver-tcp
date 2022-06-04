package main

type commandID int

const (
	CMD_USER commandID = iota
	CMD_JOIN
	CMD_CHANNELS
	CMD_MSG
	CMD_FILE
	CMD_QUIT
)

type command struct {
	id     commandID
	client *client
	args   []string
}
