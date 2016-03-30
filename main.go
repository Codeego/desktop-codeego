package main

import (
	//"os"
)

type GlobalMessage struct {
	target  string
	message interface{}
}

const (
	QuitEverything = 666
	Main           = ""
)

func main() {
	global := make(chan GlobalMessage)

	ui := UIStart(global)
	server := ServerStart(global)

	for {
		packet := <-global
		switch packet.target {
		case "ui":
			ui <- packet.message.(UIMessage)
		case "server":
			server <- packet.message.(ServerMessage)
		default:
			code := packet.message.(int)
			switch code {
			case QuitEverything:
				server <- ServerMessage{"", ServerQuit}
				println("See you soon :)")
				//os.Exit(0)
				return
			}
		}
	}
}
