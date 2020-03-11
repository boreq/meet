package controller

type ControllerState struct {
	Devices []Device `json:"devices"`
}

type Device struct {
	Id string `json:"id"`
}
