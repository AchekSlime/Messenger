package main

import (
	"messenger/ws"
)

func main() {
	server := ws.NewServer()

	structCh := make(chan struct{})
	server.StartServer(structCh)
	<-structCh
}
