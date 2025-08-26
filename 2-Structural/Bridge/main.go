package main

type Device interface {
	State() bool
	TurnOn()
	TurnOff()
	GetChannel() int
	SetChannel(int)
}

type TV struct {
	state   bool
	channel int
}

func NewTV() *TV {
	return &TV{
		state:   false,
		channel: 1,
	}
}

func (t *TV) State() bool { return t.state }

func (t *TV) TurnOn() { t.state = true }

func (t *TV) TurnOff() { t.state = false }

func (t *TV) GetChannel() int { return t.channel }

func (t *TV) SetChannel(c int) { t.channel = c }

type Radio struct {
	state   bool
	channel int
}

func NewRadio() *Radio {
	return &Radio{
		state:   false,
		channel: 1,
	}
}

func (r *Radio) State() bool { return r.state }

func (r *Radio) TurnOn() { r.state = true }

func (r *Radio) TurnOff() { r.state = false }

func (r *Radio) GetChannel() int { return r.channel }

func (r *Radio) SetChannel(c int) { r.channel = c }

type RemoteControl struct {
	device Device
}

func NewRemoteControl(d Device) *RemoteControl {
	return &RemoteControl{device: d}
}

func (rc *RemoteControl) Power() {
	if rc.device.State() {
		rc.device.TurnOff()
	} else {
		rc.device.TurnOn()
	}
}

func (rc *RemoteControl) ChannelUp() {
	c := rc.device.GetChannel()
	rc.device.SetChannel(c + 1)
}

func (rc *RemoteControl) ChannelDown() {
	c := rc.device.GetChannel()
	if c > 1 {
		rc.device.SetChannel(c - 1)
	}
}

type AdvancedRemoteControl struct {
	*RemoteControl
}

func NewAdvancedRemoteControl(d Device) *AdvancedRemoteControl {
	return &AdvancedRemoteControl{
		RemoteControl: NewRemoteControl(d),
	}
}

func (arc *AdvancedRemoteControl) Mute() {
	arc.RemoteControl.device.SetChannel(0)
}

func main() {
	tv := NewTV()
	radio := NewRadio()

	//

	remote := NewRemoteControl(tv)
	remote.Power()
	remote.ChannelUp()
	remote.ChannelDown()

	//

	arc := NewAdvancedRemoteControl(radio)
	arc.Power()
	arc.ChannelUp()
	arc.Mute()
}
