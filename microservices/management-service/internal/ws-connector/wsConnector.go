package ws_connector

type IBroker interface {
	Publish(channel string, data []byte) error
}

type WSConnector struct {
	broker IBroker
}

func New(broker IBroker) *WSConnector {
	return &WSConnector{broker: broker}
}

func (ws *WSConnector) Publish(channel string, data []byte) error {
	return ws.broker.Publish(channel, data)
}
