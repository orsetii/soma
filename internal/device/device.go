package device

type ResponseType struct {
	Requested bool
	Responded bool
}

// Information about a networking device
type DeviceInfo struct {
	IP           string
	Hostname     string
	MAC          string
	IcmpResponse ResponseType
	ArpResponse  ResponseType
	TcpResponse  ResponseType
}
