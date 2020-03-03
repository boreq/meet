package domain

type Rack struct {
	uuid       RackUUID
	lights     []Device
	airPumps   []Device
	waterPumps []Device
}

func NewRack(uuid RackUUID) (*Rack, error) {
	return &Rack{}, nil
}
