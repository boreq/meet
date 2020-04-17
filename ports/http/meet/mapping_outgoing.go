package meet

type NameChangedMsg struct {
	ParticipantUUID string `json:"participantUUID"`
	Name            string `json:"name"`
}
