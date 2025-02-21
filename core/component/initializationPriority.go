package component

//go:generate go-enum

// ENUM(
// bus,// For communication buses like i2c/spi
// io,// For components that represent GPIO pins like PCF8573
// hardware,// For components that deal with hardware and are very important like GPIO switch
// data,// For components that import data from directly connected sensors like DHT.
// processor,// For components that use data from sensors like displays
// bluetooth,
// after_bluetooth,
// wifi,
// ethernet,
// before_connection,// For components that should be initialized after WiFi and before API is connected.
// after_wifi,// For components that should be initialized after WiFi is connected.
// after_connection,// For components that should be initialized after a data connection (API/MQTT) is connected.
// late,// For components that should be initialized at the very end of the setup process.
// )
type InitializationPriority int

type WithInitializationPriorityBus struct{}

func (*WithInitializationPriorityBus) InitializationPriority() InitializationPriority {
	return InitializationPriorityBus
}

type WithInitializationPriorityIo struct{}

func (*WithInitializationPriorityIo) InitializationPriority() InitializationPriority {
	return InitializationPriorityIo
}

type WithInitializationPriorityHardware struct{}

func (*WithInitializationPriorityHardware) InitializationPriority() InitializationPriority {
	return InitializationPriorityHardware
}

type WithInitializationPriorityData struct{}

func (*WithInitializationPriorityData) InitializationPriority() InitializationPriority {
	return InitializationPriorityData
}

type WithInitializationPriorityProcessor struct{}

func (*WithInitializationPriorityProcessor) InitializationPriority() InitializationPriority {
	return InitializationPriorityProcessor
}

type WithInitializationPriorityBluetooth struct{}

func (*WithInitializationPriorityBluetooth) InitializationPriority() InitializationPriority {
	return InitializationPriorityBluetooth
}

type WithInitializationPriorityAfterBluetooth struct{}

func (*WithInitializationPriorityAfterBluetooth) InitializationPriority() InitializationPriority {
	return InitializationPriorityAfterBluetooth
}

type WithInitializationPriorityWifi struct{}

func (*WithInitializationPriorityWifi) InitializationPriority() InitializationPriority {
	return InitializationPriorityWifi
}

type WithInitializationPriorityEthernet struct{}

func (*WithInitializationPriorityEthernet) InitializationPriority() InitializationPriority {
	return InitializationPriorityEthernet
}

type WithInitializationPriorityBeforeConnection struct{}

func (*WithInitializationPriorityBeforeConnection) InitializationPriority() InitializationPriority {
	return InitializationPriorityBeforeConnection
}

type WithInitializationPriorityAfterWifi struct{}

func (*WithInitializationPriorityAfterWifi) InitializationPriority() InitializationPriority {
	return InitializationPriorityAfterWifi
}

type WithInitializationPriorityAfterConnection struct{}

func (*WithInitializationPriorityAfterConnection) InitializationPriority() InitializationPriority {
	return InitializationPriorityAfterConnection
}

type WithInitializationPriorityLate struct{}

func (*WithInitializationPriorityLate) InitializationPriority() InitializationPriority {
	return InitializationPriorityLate
}
