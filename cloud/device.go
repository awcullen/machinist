package cloud

// Device is a communication channel to a local network endpoint.
type Device interface {
	Name() string
	Location() string
	DeviceID() string
	Kind() string
	Stop() error
}
