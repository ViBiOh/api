package echo

import (
	"net/http"

	"github.com/ViBiOh/httputils/pkg/errors"
	"github.com/ViBiOh/httputils/pkg/logger"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Handler for Echo websocket request request. Should be use with net/http
func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if ws != nil {
			defer func() {
				if err := ws.Close(); err != nil {
					logger.Error("%#v", errors.WithStack(err))
				}
			}()
		}

		if err != nil {
			logger.Error("%v", err)
			return
		}

		for {
			messageType, p, err := ws.ReadMessage()
			if messageType == websocket.CloseMessage {
				return
			}

			if err != nil {
				logger.Error("%v", err)
				return
			}

			if err = ws.WriteMessage(messageType, p); err != nil {
				logger.Error("%v", err)
				return
			}
		}
	})
}
