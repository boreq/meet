package meet

type JoinMeetingRequest struct {
	Sdp string `json:"sdp"`
}

type JoinMeetingResponse struct {
	Sdp string `json:"sdp"`
}
