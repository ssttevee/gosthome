package entity

import "time"

type Optional[T any] struct {
	Has   bool
	Value T
}

type LightRGB struct {
	Red   float32
	Green float32
	Blue  float32
}

type LightCommand struct {
	State            Optional[bool]
	TransitionLength Optional[time.Duration]
	FlashLength      Optional[time.Duration]
	ColorMode        Optional[ColorMode]
	Brightness       Optional[float32]
	ColorBrightness  Optional[float32]
	Rgb              Optional[LightRGB]
	White            Optional[float32]
	ColorTemperature Optional[float32]
	ColdWhite        Optional[float32]
	WarmWhite        Optional[float32]
	Effect           Optional[uint32]
	// bool publish_{true};
	// bool save_{true};
}

func (lc LightCommand) SetState(state bool) LightCommand {
	lc.State.Has = true
	lc.State.Value = state
	return lc
}
func (lc LightCommand) SetBrightness(brightness float32) LightCommand {
	lc.Brightness.Has = true
	lc.Brightness.Value = brightness
	return lc
}
func (lc LightCommand) SetColorMode(colormode ColorMode) LightCommand {
	lc.ColorMode.Has = true
	lc.ColorMode.Value = colormode
	return lc
}
func (lc LightCommand) SetColorBrightness(colorbrightness float32) LightCommand {
	lc.ColorBrightness.Has = true
	lc.ColorBrightness.Value = colorbrightness
	return lc
}
func (lc LightCommand) SetRgb(red float32, green float32, blue float32) LightCommand {
	lc.Rgb.Has = true
	lc.Rgb.Value = LightRGB{
		Red: red, Green: green, Blue: blue,
	}
	return lc
}
func (lc LightCommand) SetWhite(white float32) LightCommand {
	lc.White.Has = true
	lc.White.Value = white
	return lc
}
func (lc LightCommand) SetColorTemperature(colortemperature float32) LightCommand {
	lc.ColorTemperature.Has = true
	lc.ColorTemperature.Value = colortemperature
	return lc
}
func (lc LightCommand) SetColdWhite(coldwhite float32) LightCommand {
	lc.ColdWhite.Has = true
	lc.ColdWhite.Value = coldwhite
	return lc
}
func (lc LightCommand) SetWarmWhite(warmwhite float32) LightCommand {
	lc.WarmWhite.Has = true
	lc.WarmWhite.Value = warmwhite
	return lc
}
func (lc LightCommand) SetTransitionLength(transitionlength time.Duration) LightCommand {
	lc.TransitionLength.Has = true
	lc.TransitionLength.Value = transitionlength
	return lc
}
func (lc LightCommand) SetFlashLength(flashlength time.Duration) LightCommand {
	lc.FlashLength.Has = true
	lc.FlashLength.Value = flashlength
	return lc
}
func (lc LightCommand) SetEffect(effect uint32) LightCommand {
	lc.Effect.Has = true
	lc.Effect.Value = effect
	return lc
}
