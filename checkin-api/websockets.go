package main

import (
	"code.google.com/p/go.net/websocket"
	"log"
	"net/http"
	"tickit/common"
)

func startSocketServer() {
	log.Printf("starting websocket listener on %s", config.Servers.WebSockets)
	http.Handle("/record", websocket.Handler(socketScanHandler))
	err := http.ListenAndServe(config.Servers.WebSockets, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

// Records scans, validating only the format of the ticket number and nothing else.
func socketScanHandler(ws *websocket.Conn) {
	defer ws.Close()

	log.Printf("Client connected: %+v", ws)

	for {
		var scans []tickit.Scan

		response := tickit.RecordAPIResponse{Success: false}

		err := websocket.JSON.Receive(ws, &scans)
		if err != nil {
			response.Message = "Error decoding JSON"
			websocket.JSON.Send(ws, response)
			return
		}

		uuids, err := SaveScans(scans)
		if err != nil {
			response.Message = err.Error()
			websocket.JSON.Send(ws, response)
			return
		}

		response.UnsavedUUIDs = uuids
		response.Success = true
		websocket.JSON.Send(ws, response)
	}
}
