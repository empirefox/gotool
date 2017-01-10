# crypt xps

1. get xps
```bash
go get -u github.com/empirefox/gotool/crypt
```

1. generate xps file from `xps-config.json`
```bash
xps -k yourpassword
```

1. generate from `xps-config-dev.json`
```bash
xps -x xps-config-dev.json
```

1. test
```bash
go test
```