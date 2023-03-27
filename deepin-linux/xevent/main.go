package main

import (
	x "github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/util/keybind"
)

func main() {
	xConn, err := x.NewConn()
	if err != nil {
		x.Logger.Println("failed to get X conn:", err)
		return
	}

	rootWin := xConn.GetDefaultScreen().Root
	err = keybind.GrabKeyboard(xConn, rootWin)
	if err != nil {
		return
	}
	eventChan := make(chan x.GenericEvent, 10)
	xConn.AddEventChan(eventChan)
	go func() {
		for ev := range eventChan {
			switch ev.GetEventCode() {
			case x.PropertyNotifyEventCode:
				event, _ := x.NewPropertyNotifyEvent(ev)
				x.Logger.Println(event)
			default:
				x.Logger.Println(ev.GetEventCode())
			}
		}
	}()
	select {}
}
