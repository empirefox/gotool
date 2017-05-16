# crypt xps

1. get xps
```bash
go get -u github.com/empirefox/gotool/crypt/xps
```

1. generate xps file from `xps-config.json`
```bash
xps -k yourpassword
```

1. generate from `xps-config-dev.json`
```bash
xps -x xps-config-dev.json
```

1. Extract to prod
```bash
xps -d ./prod [-x xps-config.json] [-k password]
```

1. test
```bash
go test
```

1. api
```go
func LoadConfig(config interface{}, opts *ConfigOptions) (err error)
```

Parse ConfigOptions from env:
```go
err := crypt.LoadConfig(config, nil)
```