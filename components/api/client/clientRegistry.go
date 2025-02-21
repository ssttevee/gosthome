package client

func (c *Client) BinarySensorByKey(key uint32) (ret *BinarySensorComponent, ok bool) {
	cm, ok := c.reg.BinarySensorByKey(key)
	if !ok {
		return
	}
	ret, ok = cm.(*BinarySensorComponent)
	return
}
func (c *Client) CoverByKey(key uint32) (ret *CoverComponent, ok bool) {
	cm, ok := c.reg.CoverByKey(key)
	if !ok {
		return
	}
	ret, ok = cm.(*CoverComponent)
	return
}
func (c *Client) FanByKey(key uint32) (ret *FanComponent, ok bool) {
	cm, ok := c.reg.FanByKey(key)
	if !ok {
		return
	}
	ret, ok = cm.(*FanComponent)
	return
}
func (c *Client) LightByKey(key uint32) (ret *LightComponent, ok bool) {
	cm, ok := c.reg.LightByKey(key)
	if !ok {
		return
	}
	ret, ok = cm.(*LightComponent)
	return
}
func (c *Client) SensorByKey(key uint32) (ret *SensorComponent, ok bool) {
	cm, ok := c.reg.SensorByKey(key)
	if !ok {
		return
	}
	ret, ok = cm.(*SensorComponent)
	return
}
func (c *Client) SwitchByKey(key uint32) (ret *SwitchComponent, ok bool) {
	cm, ok := c.reg.SwitchByKey(key)
	if !ok {
		return
	}
	ret, ok = cm.(*SwitchComponent)
	return
}
func (c *Client) ButtonByKey(key uint32) (ret *ButtonComponent, ok bool) {
	cm, ok := c.reg.ButtonByKey(key)
	if !ok {
		return
	}
	ret, ok = cm.(*ButtonComponent)
	return
}
func (c *Client) TextSensorByKey(key uint32) (ret *TextSensorComponent, ok bool) {
	cm, ok := c.reg.TextSensorByKey(key)
	if !ok {
		return
	}
	ret, ok = cm.(*TextSensorComponent)
	return
}
func (c *Client) ServiceByKey(key uint32) (ret *ServiceComponent, ok bool) {
	cm, ok := c.reg.ServiceByKey(key)
	if !ok {
		return
	}
	ret, ok = cm.(*ServiceComponent)
	return
}
func (c *Client) CameraByKey(key uint32) (ret *CameraComponent, ok bool) {
	cm, ok := c.reg.CameraByKey(key)
	if !ok {
		return
	}
	ret, ok = cm.(*CameraComponent)
	return
}
func (c *Client) ClimateByKey(key uint32) (ret *ClimateComponent, ok bool) {
	cm, ok := c.reg.ClimateByKey(key)
	if !ok {
		return
	}
	ret, ok = cm.(*ClimateComponent)
	return
}
func (c *Client) NumberByKey(key uint32) (ret *NumberComponent, ok bool) {
	cm, ok := c.reg.NumberByKey(key)
	if !ok {
		return
	}
	ret, ok = cm.(*NumberComponent)
	return
}
func (c *Client) DateByKey(key uint32) (ret *DateComponent, ok bool) {
	cm, ok := c.reg.DateByKey(key)
	if !ok {
		return
	}
	ret, ok = cm.(*DateComponent)
	return
}
func (c *Client) TimeByKey(key uint32) (ret *TimeComponent, ok bool) {
	cm, ok := c.reg.TimeByKey(key)
	if !ok {
		return
	}
	ret, ok = cm.(*TimeComponent)
	return
}
func (c *Client) DatetimeByKey(key uint32) (ret *DatetimeComponent, ok bool) {
	cm, ok := c.reg.DatetimeByKey(key)
	if !ok {
		return
	}
	ret, ok = cm.(*DatetimeComponent)
	return
}
func (c *Client) TextByKey(key uint32) (ret *TextComponent, ok bool) {
	cm, ok := c.reg.TextByKey(key)
	if !ok {
		return
	}
	ret, ok = cm.(*TextComponent)
	return
}
func (c *Client) SelectByKey(key uint32) (ret *SelectComponent, ok bool) {
	cm, ok := c.reg.SelectByKey(key)
	if !ok {
		return
	}
	ret, ok = cm.(*SelectComponent)
	return
}
func (c *Client) SirenByKey(key uint32) (ret *SirenComponent, ok bool) {
	cm, ok := c.reg.SirenByKey(key)
	if !ok {
		return
	}
	ret, ok = cm.(*SirenComponent)
	return
}
func (c *Client) LockByKey(key uint32) (ret *LockComponent, ok bool) {
	cm, ok := c.reg.LockByKey(key)
	if !ok {
		return
	}
	ret, ok = cm.(*LockComponent)
	return
}
func (c *Client) ValveByKey(key uint32) (ret *ValveComponent, ok bool) {
	cm, ok := c.reg.ValveByKey(key)
	if !ok {
		return
	}
	ret, ok = cm.(*ValveComponent)
	return
}
func (c *Client) MediaPlayerByKey(key uint32) (ret *MediaPlayerComponent, ok bool) {
	cm, ok := c.reg.MediaPlayerByKey(key)
	if !ok {
		return
	}
	ret, ok = cm.(*MediaPlayerComponent)
	return
}
func (c *Client) AlarmControlPanelByKey(key uint32) (ret *AlarmControlPanelComponent, ok bool) {
	cm, ok := c.reg.AlarmControlPanelByKey(key)
	if !ok {
		return
	}
	ret, ok = cm.(*AlarmControlPanelComponent)
	return
}
func (c *Client) EventByKey(key uint32) (ret *EventComponent, ok bool) {
	cm, ok := c.reg.EventByKey(key)
	if !ok {
		return
	}
	ret, ok = cm.(*EventComponent)
	return
}
func (c *Client) UpdateByKey(key uint32) (ret *UpdateComponent, ok bool) {
	cm, ok := c.reg.UpdateByKey(key)
	if !ok {
		return
	}
	ret, ok = cm.(*UpdateComponent)
	return
}

func (c *Client) BinarySensors() (ret []*BinarySensorComponent) {
	cmps := c.reg.BinarySensors()
	ret = make([]*BinarySensorComponent, len(cmps))
	for i, ec := range cmps {
		ret[i] = ec.(*BinarySensorComponent)
	}
	return
}
func (c *Client) Covers() (ret []*CoverComponent) {
	cmps := c.reg.Covers()
	ret = make([]*CoverComponent, len(cmps))
	for i, ec := range cmps {
		ret[i] = ec.(*CoverComponent)
	}
	return
}
func (c *Client) Fans() (ret []*FanComponent) {
	cmps := c.reg.Fans()
	ret = make([]*FanComponent, len(cmps))
	for i, ec := range cmps {
		ret[i] = ec.(*FanComponent)
	}
	return
}
func (c *Client) Lights() (ret []*LightComponent) {
	cmps := c.reg.Lights()
	ret = make([]*LightComponent, len(cmps))
	for i, ec := range cmps {
		ret[i] = ec.(*LightComponent)
	}
	return
}
func (c *Client) Sensors() (ret []*SensorComponent) {
	cmps := c.reg.Sensors()
	ret = make([]*SensorComponent, len(cmps))
	for i, ec := range cmps {
		ret[i] = ec.(*SensorComponent)
	}
	return
}
func (c *Client) Switches() (ret []*SwitchComponent) {
	cmps := c.reg.Switches()
	ret = make([]*SwitchComponent, len(cmps))
	for i, ec := range cmps {
		ret[i] = ec.(*SwitchComponent)
	}
	return
}
func (c *Client) Buttons() (ret []*ButtonComponent) {
	cmps := c.reg.Buttons()
	ret = make([]*ButtonComponent, len(cmps))
	for i, ec := range cmps {
		ret[i] = ec.(*ButtonComponent)
	}
	return
}
func (c *Client) TextSensors() (ret []*TextSensorComponent) {
	cmps := c.reg.TextSensors()
	ret = make([]*TextSensorComponent, len(cmps))
	for i, ec := range cmps {
		ret[i] = ec.(*TextSensorComponent)
	}
	return
}
func (c *Client) Services() (ret []*ServiceComponent) {
	cmps := c.reg.Services()
	ret = make([]*ServiceComponent, len(cmps))
	for i, ec := range cmps {
		ret[i] = ec.(*ServiceComponent)
	}
	return
}
func (c *Client) Cameras() (ret []*CameraComponent) {
	cmps := c.reg.Cameras()
	ret = make([]*CameraComponent, len(cmps))
	for i, ec := range cmps {
		ret[i] = ec.(*CameraComponent)
	}
	return
}
func (c *Client) Climates() (ret []*ClimateComponent) {
	cmps := c.reg.Climates()
	ret = make([]*ClimateComponent, len(cmps))
	for i, ec := range cmps {
		ret[i] = ec.(*ClimateComponent)
	}
	return
}
func (c *Client) Numbers() (ret []*NumberComponent) {
	cmps := c.reg.Numbers()
	ret = make([]*NumberComponent, len(cmps))
	for i, ec := range cmps {
		ret[i] = ec.(*NumberComponent)
	}
	return
}
func (c *Client) Dates() (ret []*DateComponent) {
	cmps := c.reg.Dates()
	ret = make([]*DateComponent, len(cmps))
	for i, ec := range cmps {
		ret[i] = ec.(*DateComponent)
	}
	return
}
func (c *Client) Times() (ret []*TimeComponent) {
	cmps := c.reg.Times()
	ret = make([]*TimeComponent, len(cmps))
	for i, ec := range cmps {
		ret[i] = ec.(*TimeComponent)
	}
	return
}
func (c *Client) Datetimes() (ret []*DatetimeComponent) {
	cmps := c.reg.Datetimes()
	ret = make([]*DatetimeComponent, len(cmps))
	for i, ec := range cmps {
		ret[i] = ec.(*DatetimeComponent)
	}
	return
}
func (c *Client) Texts() (ret []*TextComponent) {
	cmps := c.reg.Texts()
	ret = make([]*TextComponent, len(cmps))
	for i, ec := range cmps {
		ret[i] = ec.(*TextComponent)
	}
	return
}
func (c *Client) Selects() (ret []*SelectComponent) {
	cmps := c.reg.Selects()
	ret = make([]*SelectComponent, len(cmps))
	for i, ec := range cmps {
		ret[i] = ec.(*SelectComponent)
	}
	return
}
func (c *Client) Sirens() (ret []*SirenComponent) {
	cmps := c.reg.Sirens()
	ret = make([]*SirenComponent, len(cmps))
	for i, ec := range cmps {
		ret[i] = ec.(*SirenComponent)
	}
	return
}
func (c *Client) Locks() (ret []*LockComponent) {
	cmps := c.reg.Locks()
	ret = make([]*LockComponent, len(cmps))
	for i, ec := range cmps {
		ret[i] = ec.(*LockComponent)
	}
	return
}
func (c *Client) Valves() (ret []*ValveComponent) {
	cmps := c.reg.Valves()
	ret = make([]*ValveComponent, len(cmps))
	for i, ec := range cmps {
		ret[i] = ec.(*ValveComponent)
	}
	return
}
func (c *Client) MediaPlayers() (ret []*MediaPlayerComponent) {
	cmps := c.reg.MediaPlayers()
	ret = make([]*MediaPlayerComponent, len(cmps))
	for i, ec := range cmps {
		ret[i] = ec.(*MediaPlayerComponent)
	}
	return
}
func (c *Client) AlarmControlPanels() (ret []*AlarmControlPanelComponent) {
	cmps := c.reg.AlarmControlPanels()
	ret = make([]*AlarmControlPanelComponent, len(cmps))
	for i, ec := range cmps {
		ret[i] = ec.(*AlarmControlPanelComponent)
	}
	return
}
func (c *Client) Events() (ret []*EventComponent) {
	cmps := c.reg.Events()
	ret = make([]*EventComponent, len(cmps))
	for i, ec := range cmps {
		ret[i] = ec.(*EventComponent)
	}
	return
}
func (c *Client) Updates() (ret []*UpdateComponent) {
	cmps := c.reg.Updates()
	ret = make([]*UpdateComponent, len(cmps))
	for i, ec := range cmps {
		ret[i] = ec.(*UpdateComponent)
	}
	return
}
