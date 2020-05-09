package domain

type OutgoingMessage interface{}

type HelloMessage struct {
	ParticipantUUID ParticipantUUID
}

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

type RemoteSessionDescription struct {
	ParticipantUUID    ParticipantUUID
	SessionDescription SessionDescription
}

type RemoteIceCandidate struct {
	ParticipantUUID ParticipantUUID
	IceCandidate    IceCandidate
}
