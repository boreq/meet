package domain

type OutgoingMessage interface{}

type JoinedMessage struct {
	ParticipantUUID ParticipantUUID
}

type QuitMessage struct {
	ParticipantUUID ParticipantUUID
}

type NameChangedMessage struct {
	ParticipantUUID ParticipantUUID
	Name            ParticipantName
}

type RemoteSessionDescriptionReceived struct {
	ParticipantUUID    ParticipantUUID
	SessionDescription SessionDescription
}
