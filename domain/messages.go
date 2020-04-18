package domain

type OutgoingMessage interface{}

type NameChangedMessage struct {
	ParticipantUUID ParticipantUUID
	Name            ParticipantName
}
