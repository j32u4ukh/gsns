package define

const (
	CIPHER string = "GSNS"
)

// Server id
const (
	DbaServer int32 = iota
	AccountServer
	PostMessageServer
	GsnsServer
)

// TODO: Connection config
const (
	MainPort        = 1023
	DbaPort         = 1022
	AccountPort     = 1021
	PostMessagePort = 1020
)

func ServerName(id int32) string {
	switch id {
	case DbaServer:
		return "DbaServer"
	case AccountServer:
		return "DbaServer"
	case PostMessageServer:
		return "DbaServer"
	case GsnsServer:
		return "GsnsServer"
	default:
		return "Unknown"
	}
}
