package dp

var (
	Mode DevOrProd
)

type DevOrProd struct {
	IsDev      bool
	HttpProto  string
	HttpPrefix string
	WsProto    string
	WsPrefix   string
}

func init() {
	if Mode.HttpPrefix == "" {
		SetDevMode(false)
	}
}

func SetDevMode(isDev bool) {
	Mode.IsDev = isDev
	if isDev {
		Mode.HttpProto = "http"
		Mode.HttpPrefix = "http://"
		Mode.WsProto = "ws"
		Mode.WsPrefix = "ws://"
	} else {
		Mode.HttpProto = "https"
		Mode.HttpPrefix = "https://"
		Mode.WsProto = "wss"
		Mode.WsPrefix = "wss://"
	}
}
