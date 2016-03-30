package main

import (
  "./tcp_server"
  "regexp"
)


type ServerMessage struct {
	message string
	code    int
}

const (
	ServerQuit = 1
)

func ServerStart(global chan GlobalMessage) chan ServerMessage {
	messenger := make(chan ServerMessage)

	server := tcp_server.New("localhost:14322")

	go func() {
		defer server.Listen()
	}()

	go func() {
		for {
			msg := <-messenger
			switch msg.code {
			case ServerQuit:
				server.Close()
				return
			}
		}
	}()

  server.OnNewMessage(func(c *tcp_server.Client, message string) {
    re := regexp.MustCompile("\n")
    message = re.ReplaceAllString(message, "")

    global <- GlobalMessage{"ui", UIMessage{message, PushMessage}}
  })

	return messenger
}
