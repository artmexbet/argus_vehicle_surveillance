package router

import (
	"github.com/gofiber/websocket/v2"
	"log/slog"
)

func (r *Router) GetInfo(c *websocket.Conn) {
	defer c.Close()
	for msg := range r.svc.GetChannel() {
		slog.Info("Read msg from channel")
		err := c.WriteJSON(msg)
		slog.Info("Send JSON msg")
		if err != nil {
			slog.Error(err.Error())
			return
		}
	}
}
