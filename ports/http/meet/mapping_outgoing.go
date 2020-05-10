package meet

type HelloMsg struct {
	ParticipantUUID string `json:"participantUUID"`
}

type JoinedMsg struct {
	ParticipantUUID string `json:"participantUUID"`
}

type QuitMsg struct {
	ParticipantUUID string `json:"participantUUID"`
}

type NameChangedMsg struct {
	ParticipantUUID string `json:"participantUUID"`
	Name            string `json:"name"`
}

type VisualisationStateMsg struct {
	ParticipantUUID string `json:"participantUUID"`
	State           string `json:"state"`
}

type RemoteSessionDescriptionMsg struct {
	ParticipantUUID    string `json:"participantUUID"`
	SessionDescription string `json:"sessionDescription"`
}

type RemoteIceCandidateMsg struct {
	ParticipantUUID string `json:"participantUUID"`
	IceCandidate    string `json:"iceCandidate"`
}
