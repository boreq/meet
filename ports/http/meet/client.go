package meet

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/boreq/errors"
	"github.com/boreq/meet/application/meet"
	"github.com/boreq/meet/domain"
	"github.com/boreq/meet/internal/logging"
	"github.com/gorilla/websocket"
)

type Client struct {
	receive chan meet.IncomingMessage
	send    chan domain.OutgoingMessage
	conn    *websocket.Conn
	log     logging.Logger
}

func NewClient(conn *websocket.Conn) meet.Client {
	c := &Client{
		conn:    conn,
		receive: make(chan meet.IncomingMessage),
		send:    make(chan domain.OutgoingMessage),
		log:     logging.New("ports/client"),
	}

	go c.runReceive()
	go c.runSend()

	return meet.Client{
		Receive: c.receive,
		Send:    c.send,
	}
}

func (c *Client) runReceive() {
	defer close(c.receive)

	for {
		websocketMessageType, reader, err := c.conn.NextReader()
		if err != nil {
			c.log.Debug("next reader error", "err", err)
			return
		}

		if websocketMessageType != websocket.BinaryMessage {
			c.log.Warn("invalid received message type", "websocketMessageType", websocketMessageType)
			return
		}

		msg, err := c.unmarshalMessage(reader)
		if err != nil {
			c.log.Warn("could not unmarshal the incoming message", "err", err)
			return
		}

		c.receive <- msg
	}
}

func (c *Client) runSend() {
	defer func() {
		if err := c.conn.Close(); err != nil {
			c.log.Warn("could not close the connection", "err", err)
		}
	}()

	for {
		msg, ok := <-c.send
		if !ok {
			return
		}

		b, err := c.marshalMessage(msg)
		if err != nil {
			c.log.Warn("could not marshal a message", "err", err)
			return
		}

		if err := c.conn.WriteMessage(websocket.BinaryMessage, b); err != nil {
			c.log.Warn("could not write a message", "err", err)
			return
		}
	}

}

func (c *Client) unmarshalMessage(reader io.Reader) (meet.IncomingMessage, error) {
	var message Message
	if err := json.NewDecoder(reader).Decode(&message); err != nil {
		return nil, errors.Wrap(err, "could not decode the incoming message")
	}

	return c.unmarshalIncomingMessage(message)
}

func (c *Client) unmarshalIncomingMessage(message Message) (meet.IncomingMessage, error) {
	switch message.MessageType {
	case SetNameMessage:
		var msg meet.SetNameMessage
		if err := json.Unmarshal(message.Payload, &msg); err != nil {
			return nil, errors.Wrap(err, "could not unmarshal set name message")
		}
		return msg, nil
	default:
		return nil, fmt.Errorf("unsupported message type '%+v'", message.MessageType)
	}
}

func (c *Client) marshalMessage(msg domain.OutgoingMessage) ([]byte, error) {
	message, err := c.marshalOutgoingMessage(msg)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal outgoing message")
	}

	return json.Marshal(message)
}

func (c *Client) marshalOutgoingMessage(msg domain.OutgoingMessage) (Message, error) {
	b, err := json.Marshal(msg)
	if err != nil {
		return Message{}, errors.Wrap(err, "could not marshal the message")
	}

	messageType, err := c.getMessageType(msg)
	if err != nil {
		return Message{}, errors.Wrap(err, "could not get a message type")
	}

	return Message{
		MessageType: messageType,
		Payload:     b,
	}, nil
}

func (c *Client) getMessageType(msg domain.OutgoingMessage) (MessageType, error) {
	switch msg.(type) {
	case domain.NameChangedMessage:
		return NameChangedMessage, nil
	default:
		return "", fmt.Errorf("unsupported message: '%T'", msg)
	}
}

type Message struct {
	MessageType MessageType `json:"messageType"`
	Payload     []byte      `json:"payload"`
}

type MessageType string

const (
	SetNameMessage     MessageType = "set_name"
	NameChangedMessage MessageType = "name_changed"
)
