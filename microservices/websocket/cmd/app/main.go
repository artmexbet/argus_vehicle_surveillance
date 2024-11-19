package main

import (
	"Argus/microservices/websocket/internal/router"
	"Argus/microservices/websocket/internal/service"
	"Argus/pkg/models"
	natsConnector "Argus/pkg/nats-connector"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"log/slog"
)

func main() {
	natsConn, err := natsConnector.New(&natsConnector.Config{Host: "broker", Port: "4222"})
	if err != nil {
		log.Fatal(err)
	}

	svc := service.New()
	_, err = natsConn.CreateReader("camera-01-ws", func(msg *nats.Msg) {
		slog.Info("Read msg from camera-01-ws")
		var yoloJson models.YOLOJson
		slog.Info(string(msg.Data))
		err := json.Unmarshal(msg.Data, &yoloJson)
		if err != nil {
			log.Println(err)
			return
		}
		go svc.Send(yoloJson)
		slog.Info("Send msg to websocket")
	})

	r := router.New(&router.Config{Port: ":3000"}, svc)
	log.Fatal(r.Start())
}
