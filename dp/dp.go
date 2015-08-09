package dp

import "github.com/empirefox/gotool/paas"

var (
	Mode DevOrProd
)

type DevOrProd struct {
	IsDev     bool
	HttpProto string
	HttpPort  string
	WsProto   string
	WsPort    string
}

func init() {
	if Mode.HttpProto == "" {
		SetDevMode(false)
	}
}

func SetDevMode(isDev bool) {
	Mode.IsDev = isDev
	if isDev {
		Mode.HttpProto = "http"
		Mode.HttpPort = paas.PortInTest
		Mode.WsProto = "ws"
		Mode.WsPort, _ = paas.GetWsPorts()
	} else {
		Mode.HttpProto = "https"
		Mode.HttpPort = ""
		Mode.WsProto = "wss"
		_, Mode.WsPort = paas.GetWsPorts()
	}
}
