package ws_connector

import (
	"fmt"
	"testing"
)

type Broker struct {
	Info    map[string][][]byte
	WantErr bool
}

func (b *Broker) Publish(channel string, data []byte) error {
	if b.WantErr {
		return fmt.Errorf("error")
	}
	if _, ok := b.Info[channel]; !ok {
		b.Info[channel] = make([][]byte, 0)
	}
	b.Info[channel] = append(b.Info[channel], data)
	return nil
}

func TestWSConnector_Publish(t *testing.T) {
	b := &Broker{
		Info:    make(map[string][][]byte),
		WantErr: false,
	}
	ws := New(b)
	err := ws.Publish("channel", []byte("test"))
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if len(b.Info["channel"]) != 1 {
		t.Fatalf("expected 1, got %d", len(b.Info["channel"]))
	}
	if string(b.Info["channel"][0]) != "test" {
		t.Fatalf("expected test, got %s", string(b.Info["channel"][0]))
	}
}

func TestWSConnector_Publish_Error(t *testing.T) {
	b := &Broker{
		Info:    make(map[string][][]byte),
		WantErr: true,
	}
	ws := New(b)
	err := ws.Publish("channel", []byte("test"))
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
