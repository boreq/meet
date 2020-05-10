package domain

type VisualisationState struct {
	state string
}

func NewVisualisationState(state string) (VisualisationState, error) {
	return VisualisationState{
		state: state,
	}, nil
}

func (s VisualisationState) String() string {
	return s.state
}
