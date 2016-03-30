package main

import (
	ui "github.com/gizak/termui"
	//"fmt"
	"strings"
)

var banner = ` _____           _                       
/  __ \         | |                      
| /  \/ ___   __| | ___  ___  __ _  ___  
| |    / _ \ / _' |/ _ \/ _ \/ _' |/ _ \ 
| \__/\ (_) | (_| |  __/  __/ (_| | (_) |
 \____/\___/ \__,_|\___|\___|\__, |\___/ 
                              __/ |      
                             |___/       `

type UIMessage struct {
	message string
	code    int
}

type uiMessage struct {
	message string
}

var messages = []uiMessage{}

func getMessages() []string {
	ret := []string{}
	slice := len(messages)-18
	if(slice < 0) {
		slice = 0
	}
	for _, item := range messages[slice:] {
		ret = append(ret, item.message)
	}

	return ret
}

var initialized = false

const (
	RenderUI    = 0
	PushMessage = 1
	QuitUI      = 2
)

type Interface struct {
	tip 		*ui.Par
	welcome *ui.Par
	console *ui.Par
}

func NewInterface() Interface {
	ui.ColorMap = map[string]ui.Attribute{
		"fg":           ui.ColorRGB(0, 0, 0),
		"bg":           ui.ColorRGB(0, 0, 0),
		"border.fg":    ui.ColorRGB(0, 0, 0),
		"label.fg":     ui.ColorRGB(0, 0, 0),
		"par.fg":       ui.ColorRGB(0, 0, 0),
		"par.label.bg": ui.ColorRGB(0, 0, 0),
	}

	return Interface {
		ui.NewPar(""),
		ui.NewPar(banner),
		ui.NewPar(""),
	}
}

func (i *Interface) Render() {
	if !initialized {
		initialized = true
		ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(12, 0, i.welcome)),
		ui.NewRow(
			ui.NewCol(12, 0, i.tip)),
		ui.NewRow(
			ui.NewCol(12, 0, i.console)))
		i.tip.Height = 3
		i.tip.Width = 50
		i.tip.BorderLabel = "Tip"

		i.welcome.Height = 10

		i.console.BorderLabel = "Messages"
		i.console.Height = 20
	}

	//ui.Body = ui.NewGrid()
  //ui.Body.Width = ui.TermWidth()
	
	i.tip.Text = ":PRESS CTRL+Q TO QUIT"
	
	i.console.Text = strings.Join(getMessages(), "\n")

	ui.Body.Align()
	ui.Render(ui.Body)
}

var defaultInterface = NewInterface()

func UIStart(global chan GlobalMessage) chan UIMessage {
	messenger := make(chan UIMessage)

	go func() {
		for {
			msg := <-messenger
			switch msg.code {
			case QuitUI:
				ui.StopLoop()
				ui.Close()
				return
			case RenderUI:
				defaultInterface.Render()
			case PushMessage:
				messages = append(messages, uiMessage{msg.message})
				defaultInterface.Render()
			}
		}
	}()

	go func() {
		err := ui.Init()
		if err != nil {
			panic(err)
		}

		ui.Handle("/sys/kbd/C-q", func(ui.Event) {
			messenger <- UIMessage{"", QuitUI}
			global <- GlobalMessage{Main, QuitEverything}
		})

		ui.Handle("/sys/wnd/resize", func(ui.Event) {
			messenger <- UIMessage{"", RenderUI}			
		})

		messenger <- UIMessage{"", RenderUI}

		go ui.Loop()

	}()

	return messenger
}
