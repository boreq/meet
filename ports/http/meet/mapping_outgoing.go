package meet

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
