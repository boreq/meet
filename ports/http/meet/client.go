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

		if websocketMessageType != websocket.TextMessage {
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

		if err := c.conn.WriteMessage(websocket.TextMessage, b); err != nil {
			c.log.Warn("could not write a message", "err", err)
			return
		}
	}

}

func (c *Client) unmarshalMessage(reader io.Reader) (meet.IncomingMessage, error) {
	var message IncomingMessage
	if err := json.NewDecoder(reader).Decode(&message); err != nil {
		return nil, errors.Wrap(err, "could not decode the incoming message")
	}

	return c.unmarshalIncomingMessage(message)
}

func (c *Client) unmarshalIncomingMessage(message IncomingMessage) (meet.IncomingMessage, error) {
	mapping, ok := incomingMapping[message.MessageType]
	if !ok {
		return nil, fmt.Errorf("unsupported message type '%+v'", message.MessageType)
	}
	return mapping([]byte(message.Payload))
}

func (c *Client) marshalMessage(msg domain.OutgoingMessage) ([]byte, error) {
	message, err := c.marshalOutgoingMessage(msg)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal outgoing message")
	}

	b, err := json.Marshal(message)
	if err != nil {
		return nil, errors.Wrap(err, "json marshal failed")
	}

	return b, nil
}

func (c *Client) marshalOutgoingMessage(msg domain.OutgoingMessage) (OutgoingMessage, error) {
	message, messageType, err := c.toTransportMessage(msg)
	if err != nil {
		return OutgoingMessage{}, errors.Wrap(err, "could not convert to transport message")
	}

	b, err := json.Marshal(message)
	if err != nil {
		return OutgoingMessage{}, errors.Wrap(err, "could not marshal the message")
	}

	return OutgoingMessage{
		MessageType: messageType,
		Payload:     string(b),
	}, nil
}

func (c *Client) toTransportMessage(message domain.OutgoingMessage) (interface{}, OutgoingMessageType, error) {
	switch msg := message.(type) {
	case domain.HelloMessage:
		return HelloMsg{
			ParticipantUUID: msg.ParticipantUUID.String(),
		}, HelloMessage, nil
	case domain.JoinedMessage:
		return JoinedMsg{
			ParticipantUUID: msg.ParticipantUUID.String(),
		}, JoinedMessage, nil
	case domain.QuitMessage:
		return QuitMsg{
			ParticipantUUID: msg.ParticipantUUID.String(),
		}, QuitMessage, nil
	case domain.NameChangedMessage:
		return NameChangedMsg{
			ParticipantUUID: msg.ParticipantUUID.String(),
			Name:            msg.Name.String(),
		}, NameChangedMessage, nil
	case domain.RemoteSessionDescription:
		return RemoteSessionDescriptionMsg{
			ParticipantUUID:    msg.ParticipantUUID.String(),
			SessionDescription: msg.SessionDescription.String(),
		}, RemoteSessionDescriptionMessage, nil
	case domain.RemoteIceCandidate:
		return RemoteIceCandidateMsg{
			ParticipantUUID: msg.ParticipantUUID.String(),
			IceCandidate:    msg.IceCandidate.String(),
		}, RemoteIceCandidateMessage, nil
	default:
		return nil, "", fmt.Errorf("unsupported message: '%T'", msg)
	}
}

type IncomingMessage struct {
	MessageType IncomingMessageType `json:"messageType"`
	Payload     string              `json:"payload"`
}

type IncomingMessageType string

const (
	SetNameMessage                 IncomingMessageType = "setName"
	LocalSessionDescriptionMessage IncomingMessageType = "localSessionDescription"
	LocalIceCandidateMessage       IncomingMessageType = "localIceCandidate"
)

type OutgoingMessage struct {
	MessageType OutgoingMessageType `json:"messageType"`
	Payload     string              `json:"payload"`
}

type OutgoingMessageType string

const (
	HelloMessage                    OutgoingMessageType = "hello"
	JoinedMessage                   OutgoingMessageType = "joined"
	QuitMessage                     OutgoingMessageType = "quit"
	NameChangedMessage              OutgoingMessageType = "nameChanged"
	RemoteSessionDescriptionMessage OutgoingMessageType = "remoteSessionDescription"
	RemoteIceCandidateMessage       OutgoingMessageType = "remoteIceCandidate"
)
