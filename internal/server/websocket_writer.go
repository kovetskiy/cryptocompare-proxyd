package server

import "github.com/gorilla/websocket"

type websocketWriter struct {
	connection *websocket.Conn
}

func (writer websocketWriter) Write(data []byte) (int, error) {
	underlying, err := writer.connection.NextWriter(websocket.TextMessage)
	if err != nil {
		return 0, err
	}

	defer underlying.Close()

	return underlying.Write(data)
}
